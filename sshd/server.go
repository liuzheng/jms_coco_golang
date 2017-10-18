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
	"strings"
)

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
	API       *api.Server
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
	if log.HandleErr("sshd Serve", err) {
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
		sshConn.Close()
		return
	}

	switch newChannel.ChannelType() {
	case "session":
		sesschan, _, err := newChannel.Accept()
		if err != nil {
			log.Panic("session", "%v", err)
		}
		conn := rw{Reader: sesschan, Writer: sesschan.Stderr()}
		menu, err := NewMenu(conn, session, s.API)
		if err != nil {
			log.Error("Session", "%v", err)
			session.Close()
			return
		}

		menu.Welcome()
	loop:
		for {
			fmt.Fprint(conn, "\r\nOpt>")
			var buf []byte
			b := make([]byte, 1)
			for {
				n, err := conn.Read(b)
				if err != nil {
					log.Error("Server", "%v", err)
					fmt.Fprint(conn, "^D")
					fmt.Fprint(conn, "\r\nGoodbye\r\n")
					break loop
				}
				if n >= 0 {
					switch b[0] {
					case '\r':
						if len(buf) == 1 {
							fmt.Fprint(conn, "\r\n")
							switch strings.ToUpper(string(buf[0])) {
							case "P":
								// 输入 P/p 显示您有权限的主机.
								menu.GetMachineList()
							case "G":
								// 输入 G/g 显示您有权限的主机组.
								menu.GetHostGroup()
								//case "E":
								////  输入 E/e 批量执行命令.(未完成)
								//case "U":
								////  输入 U/u 批量上传文件.(未完成)
								//case "D":
								////  输入 D/d 批量下载文件.(未完成)
							case "H":
								//  输入 H/h 帮助.
								menu.GetHelp()
							case "Q":
								//  输入 Q/q 退出.
								fmt.Fprint(conn, "Goodbye\r\n")
								break loop
							default:
								fmt.Fprint(conn, "TO BE CONTINUED")
							}
						} else if len(buf) > 1 {
							fmt.Fprint(conn, "\r\n")
							switch strings.ToUpper(string(buf[0])) {
							case "/":
								// 输入 / + IP, 主机名 or 备注 搜索. 如: /ip
								menu.Search(string(buf[1:]))
							case "G":
								//  输入 G/g + 组ID 显示该组下主机. 如: g1
								menu.GetHostGroupList(string(buf[1:]))
							default:
								// 输入 ID 直接登录 或 输入部分 IP,主机名,备注 进行搜索登录(如果唯一).

							}
						}
						continue loop
						//res, err := strconv.ParseInt(string(buf), 10, 64)
						//if log.HandleErr("DefaultInteractive", err) {
						//	fmt.Fprint(conn, "input not a valid integer. Please try again")
						//	continue loop
						//}
						//if int(res) >= 2 || res < 0 {
						//	fmt.Fprint(conn, "No such server. Please try again")
						//	continue loop
						//}
						//return remotes[int(res)], nil
					case 0x03:
						fmt.Fprint(conn, "^C")
						continue loop
					case 0x04:
						fmt.Fprint(conn, "^D\r\nGoodbye\r\n")
						break loop
						//return api.Machine{}, errors.New("user terminated session")
					case 0x7f:
						if l := len(buf); l > 0 {
							buf = buf[:l-1]
							fmt.Fprint(conn, "\b \b")
						}
						continue
					default:
						fmt.Fprintf(conn, "%s", b)
					}
					buf = append(buf, b[0])
				}
			}
		}
		session.Close()
		//for {
		//	s.SessionForward(session, newChannel, chans)
		//}
		//sesschan.Close()

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
	as := api.New()

	hostPrivateKey, err := ioutil.ReadFile(*util.Hostkey)
	if log.HandleErr("sshd Run", err) {
		panic(err)
	}

	signer, err := ssh.ParsePrivateKey(hostPrivateKey)
	if log.HandleErr("sshd Run", err) {
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
		//log.HandleErr("sshd New", aerr) //TODO： 这里的方法有点问题
		authFile := []byte(PublicKey.Key)

		candidate, _, _, _, _ = ssh.ParseAuthorizedKey(authFile)
		if t == candidate.Type() && bytes.Compare(k, candidate.Marshal()) == 0 {
			return &user, nil
		}

		log.Warn("sshd auth", "%s: access denied (username: %s)", c.RemoteAddr(), c.User())
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
		session.Machines, _ = as.GetList("", 0)
		//log.HandleErr("sshd Run", aerr) //TODO： 这里的方法有点问题
		return nil
	}
	server := &Server{
		Auther: auth,
		Setup:  setup,
		API:    as,
	}

	server.sshConfig = &ssh.ServerConfig{
		PublicKeyCallback: server.auth,
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
