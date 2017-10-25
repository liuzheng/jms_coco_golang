package sshd

import (
	//"errors"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"coco/util/log"
	"bytes"
	"coco/api"
	"coco/util"
	"fmt"
	"coco/util/errors"
)

// Server is the coco server instance.
type Server struct {
	// Auther checks if a connection is permitted, and returns a user if
	// recognized.. Returning nil error indicates that the login was allowed,
	// regardless of whether the user was recognized or not. To disallow a
	// connection, return an error.
	Auther func(ssh.ConnMetadata, ssh.PublicKey) (*User, errors.Error)

	// Setup takes a Session, the most important task being filling out the
	// permitted remote hosts. Returning an error here will send the error to
	// the user and terminate the connection. This is not as clean as denying
	// the user in Auther, but can be used in case the denial was too dynamic.
	Setup func(*Session) errors.Error

	// Interactive is called to ask the user to select a host on the list of
	// potential remote hosts. This is only called in the case wehre more than
	// one option is available. If an error is returned, it is presented to the
	// user and the connection is terminated. The io.ReadWriter is to be used
	// for user interaction.
	Interactive func(io.ReadWriter, *Session) (api.Machine, errors.Error)

	// Selected is called when a remote host has been decided upon. The main
	// purpose of this callback is logging, but returning an error will
	// terminate the connection, allowing it to be used as a last-minute
	// bailout.
	//Selected  func(*Session, string) error
	sshConfig *ssh.ServerConfig
}

// HandleConn takes a net.Conn and runs it through coco.
func (s *Server) HandleConn(c net.Conn) {
	sshConn, chans, reqs, err := ssh.NewServerConn(c, s.sshConfig)
	if err == io.EOF {
		log.Info("HandleConn", "User leave")
		return
	} else if log.HandleErr("HandleConn", err, "Network issue") {
		return
	}
	defer sshConn.Close()

	if sshConn.Permissions == nil || sshConn.Permissions.Extensions == nil {
		return
	}
	log.Debug("HandleConn", "%v", sshConn.Permissions)

	ext := sshConn.Permissions.Extensions
	pk := &publicKey{
		publicKey:     []byte(ext["pubKey"]),
		publicKeyType: ext["pubKeyType"],
	}

	user, err := s.Auther(sshConn, pk)
	if log.HandleErr("sshd Serve", err) {
		return
	}
	session := &Session{
		Conn:      sshConn,
		User:      user,
		PublicKey: pk,
	}

	s.Setup(session)

	go ssh.DiscardRequests(reqs)
	newChannel := <-chans
	if newChannel == nil {
		return
	}

	switch newChannel.ChannelType() {
	case "session":
		sesschan, _, err := newChannel.Accept()
		if err != nil {
			log.Panic("session", "%v", err)
		}
		conn := rw{Reader: sesschan, Writer: sesschan.Stderr()}
		menu, err := NewMenu(conn, session)
		if log.HandleErr("Session", err, "") {
			return
		}

		menu.Welcome()
		menu.Manager()
	case "direct-tcpip":
		s.ChannelForward(session, newChannel)
	default:
		newChannel.Reject(ssh.UnknownChannelType, "connection flow not supported by coco")
	}
	return
}

// Serve is an Accept loop that sends the accepted connections through
// HandleConn.
func (s *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if log.HandleErr("sshd Serve", err) {
			return err
		}

		go s.HandleConn(conn)
	}
}

func (s *Server) auth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	perm := &ssh.Permissions{
		Extensions: map[string]string{
			"pubKey":     string(key.Marshal()),
			"pubKeyType": key.Type(),
		},
	}

	_, err := s.Auther(conn, key)
	if err == nil {
		return perm, nil
	}

	return nil, err
}

// New returns a Server initialized with the provided signer and callbacks.
func New() *Server {
	hostPrivateKey, err := ioutil.ReadFile(*util.Hostkey)
	if log.HandleErr("sshd Run", err) {
		panic(err)
	}

	signer, err := ssh.ParsePrivateKey(hostPrivateKey)
	if log.HandleErr("sshd Run", err) {
		panic(err)
	}
	// sshforward setup
	auth := func(c ssh.ConnMetadata, key ssh.PublicKey) (user *User, autherr errors.Error) {
		var candidate ssh.PublicKey
		log.Debug("sshd auth", "start auth : %v", c.User())
		as := api.New()

		t := key.Type()
		k := key.Marshal()
		PublicKey, autherr := as.GetUserPubKey(c.User())

		if log.HandleErr("sshd auth", autherr, autherr.Error()) {
			return nil, autherr
		}
		Token, err := as.GetLoginToken(c.User(), PublicKey.Ticket)
		if log.HandleErr("Login", err, "login issue") {
			return nil, err
		}
		user = &User{
			Name:  c.User(),
			Token: Token,
			Api:   as,
		}
		authFile := []byte(PublicKey.Key)
		log.Debug("sshd auth", "%v", PublicKey)

		candidate, _, _, _, _ = ssh.ParseAuthorizedKey(authFile)
		if t == candidate.Type() && bytes.Compare(k, candidate.Marshal()) == 0 {
			return user, nil
		}

		log.Warn("sshd auth", "%s: access denied (username: %s)", c.RemoteAddr(), c.User())
		return nil, errors.New("access denied", 403)
	}

	setup := func(session *Session) (setuperr errors.Error) {
		var username string
		if session.User != nil {
			username = session.User.Name
		} else {
			username = "unknown user"
			return errors.New("unknow User", 404)
		}
		log.Info("sshd setup", "%s: %s authorized (username: %s)", session.Conn.RemoteAddr(), username, session.Conn.User())

		session.Machines, setuperr = session.User.Api.GetList("", 0)
		if log.HandleErr("sshd setup", setuperr, "") {
			return setuperr
		}
		return nil
	}
	server := &Server{
		Auther: auth,
		Setup:  setup,
	}

	server.sshConfig = &ssh.ServerConfig{
		PublicKeyCallback: server.auth,
		//ServerVersion: util.Version,
	}
	server.sshConfig.AddHostKey(signer)

	return server
}

func Run() {

	server := New()

	// Set up listener
	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", *util.Ip, *util.SshPort))
	if err != nil {
		panic(err)
	}

	server.Serve(l)
}
