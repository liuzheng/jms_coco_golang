package sshd

import (
	"errors"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"coco/util/log"
	"bytes"
	"coco/api"
	"coco/util"
	"fmt"
)

var (
	Hostkey = util.Hostkey
)

// User describes an authenticable user.
type User struct {
	// The public key of the user.
	PublicKey ssh.PublicKey
	AuthKeys  string
	// The name the user will be referred to as. *NOT* the username used when
	// starting the session.
	Name string
}

// Session describes the current user session.
type Session struct {
	// Conn is the ssh.ServerConn associated with the connection.
	Conn *ssh.ServerConn

	// User is the current user, or nil if unknown.
	User *User

	// Remotes is the allowed set of remote hosts.
	//Remotes []string

	// PublicKey is the public key used in this session.
	PublicKey ssh.PublicKey

	Machines []api.Machine
}

// Server is the sshmux server instance.
type Server struct {
	// Auther checks if a connection is permitted, and returns a user if
	// recognized.. Returning nil error indicates that the login was allowed,
	// regardless of whether the user was recognized or not. To disallow a
	// connection, return an error.
	Auther func(ssh.ConnMetadata, ssh.PublicKey) (*User, error)

	// Setup takes a Session, the most important task being filling out the
	// permitted remote hosts. Returning an error here will send the error to
	// the user and terminate the connection. This is not as clean as denying
	// the user in Auther, but can be used in case the denial was too dynamic.
	Setup func(*Session) error

	// Interactive is called to ask the user to select a host on the list of
	// potential remote hosts. This is only called in the case wehre more than
	// one option is available. If an error is returned, it is presented to the
	// user and the connection is terminated. The io.ReadWriter is to be used
	// for user interaction.
	Interactive func(io.ReadWriter, *Session) (api.Machine, error)

	// Selected is called when a remote host has been decided upon. The main
	// purpose of this callback is logging, but returning an error will
	// terminate the connection, allowing it to be used as a last-minute
	// bailout.
	//Selected  func(*Session, string) error
	sshConfig *ssh.ServerConfig
}

type publicKey struct {
	publicKey     []byte
	publicKeyType string
}

func (p *publicKey) Marshal() []byte {
	b := make([]byte, len(p.publicKey))
	copy(b, p.publicKey)
	return b
}

func (p *publicKey) Type() string {
	return p.publicKeyType
}

func (p *publicKey) Verify([]byte, *ssh.Signature) error {
	return errors.New("verify not implemented")
}

// HandleConn takes a net.Conn and runs it through sshmux.
func (s *Server) HandleConn(c net.Conn) {
	sshConn, chans, reqs, err := ssh.NewServerConn(c, s.sshConfig)
	if err != nil {
		c.Close()
		return
	}

	if sshConn.Permissions == nil || sshConn.Permissions.Extensions == nil {
		sshConn.Close()
		return
	}

	ext := sshConn.Permissions.Extensions
	pk := &publicKey{
		publicKey:     []byte(ext["pubKey"]),
		publicKeyType: ext["pubKeyType"],
	}

	user, err := s.Auther(sshConn, pk)

	session := &Session{
		Conn:      sshConn,
		User:      user,
		PublicKey: pk,
	}

	s.Setup(session)

	go ssh.DiscardRequests(reqs)
	newChannel := <-chans
	if newChannel == nil {
		sshConn.Close()
		return
	}

	switch newChannel.ChannelType() {
	case "session":
		s.SessionForward(session, newChannel, chans)
	case "direct-tcpip":
		s.ChannelForward(session, newChannel)
	default:
		newChannel.Reject(ssh.UnknownChannelType, "connection flow not supported by sshmux")
	}
}

// Serve is an Accept loop that sends the accepted connections through
// HandleConn.
func (s *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go s.HandleConn(conn)
	}
}

func (s *Server) auth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	k := key.Marshal()
	t := key.Type()
	perm := &ssh.Permissions{
		Extensions: map[string]string{
			"pubKey":     string(k),
			"pubKeyType": t,
		},
	}

	_, err := s.Auther(conn, key)
	if err == nil {
		return perm, nil
	}

	return nil, err
}

// New returns a Server initialized with the provided signer and callbacks.
func New(as *api.Server) *Server {

	hostPrivateKey, err := ioutil.ReadFile(*Hostkey)
	if err != nil {
		panic(err)
	}

	signer, err := ssh.ParsePrivateKey(hostPrivateKey)
	if err != nil {
		log.Error("sshd Run", "%v", err)
		panic(err)
	}
	// sshforward setup
	auth := func(c ssh.ConnMetadata, key ssh.PublicKey) (*User, error) {
		var candidate ssh.PublicKey
		t := key.Type()
		k := key.Marshal()
		log.Debug("auth", "%v", c.User())
		user := User{Name: c.User()}
		PublicKey, _ := as.GetUserPubKey(c.User())
		authFile := []byte(PublicKey.Key)

		candidate, _, _, _, _ = ssh.ParseAuthorizedKey(authFile)
		if t == candidate.Type() && bytes.Compare(k, candidate.Marshal()) == 0 {
			return &user, nil
		}

		//if hasDefaults {
		//	return nil, nil
		//}

		//log.Warn("sshd auth", "%s: access denied (username: %s)", c.RemoteAddr(), c.User())
		return nil, errors.New("access denied")
	}

	setup := func(session *Session) error {
		var username string
		if session.User != nil {
			username = session.User.Name
		} else {
			username = "unknown user"
		}
		log.Info("sshd setup", "%s: %s authorized (username: %s)", session.Conn.RemoteAddr(), username, session.Conn.User())
		session.Machines, _ = as.GetList()
		//if Hosts == nil {
		//	// TODO: need write more code
		//	log.Error("sshd", "got nothing need write more code")
		//}
		////outer:
		//for _, h := range Hosts {
		//	session.Remotes = append(session.Remotes, h.Ip)
		//}
		//session.Machines =
		return nil
	}
	server := &Server{
		Auther: auth,
		Setup:  setup,
	}

	server.sshConfig = &ssh.ServerConfig{
		PublicKeyCallback: server.auth,
	}
	server.sshConfig.AddHostKey(signer)

	return server
}

func Run() {
	as := api.New()

	sshserver := New(as)
	//sshserver.Selected = func(session *Session, remote string) error {
	//	var username string
	//	if session.User != nil {
	//		username = session.User.Name
	//	} else {
	//		username = "unknown user"
	//	}
	//	log.Debug("sshserver", "%s: %s connecting to %s", session.Conn.RemoteAddr(), username, remote)
	//	return nil
	//}
	// Set up listener
	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", as.Ip, as.SshPort))
	if err != nil {
		panic(err)
	}

	sshserver.Serve(l)
}
