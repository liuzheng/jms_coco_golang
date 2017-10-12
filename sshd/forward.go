package sshd

import (
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
	//"golang.org/x/crypto/ssh/agent"
	"coco/api"
)

func proxy(reqs1, reqs2 <-chan *ssh.Request, channel1, channel2 ssh.Channel) {
	var closer sync.Once
	closeFunc := func() {
		channel1.Close()
		channel2.Close()
	}

	defer closer.Do(closeFunc)

	closerChan := make(chan bool, 1)

	go func() {
		io.Copy(channel1, channel2)
		closerChan <- true
	}()

	go func() {
		io.Copy(channel2, channel1)
		closerChan <- true
	}()

	for {
		select {
		case req := <-reqs1:
			if req == nil {
				return
			}
			b, err := channel2.SendRequest(req.Type, req.WantReply, req.Payload)
			if err != nil {
				return
			}
			req.Reply(b, nil)

		case req := <-reqs2:
			if req == nil {
				return
			}
			b, err := channel1.SendRequest(req.Type, req.WantReply, req.Payload)
			if err != nil {
				return
			}
			req.Reply(b, nil)
		case <-closerChan:
			return
		}
	}
}

type channelOpenDirectMsg struct {
	RAddr string
	RPort uint32
	LAddr string
	LPort uint32
}

// ChannelForward establishes a secure channel forward (ssh -W) to the server
// requested by the user, assuming it is a permitted host.
func (s *Server) ChannelForward(session *Session, newChannel ssh.NewChannel) {
	var msg channelOpenDirectMsg
	ssh.Unmarshal(newChannel.ExtraData(), &msg)
	address := fmt.Sprintf("%s:%d", msg.RAddr, msg.RPort)

	permitted := false
	//for _, remote := range session.Machines {
	//	if remote == address {
	//		permitted = true
	//		break
	//	}
	//}

	if !permitted {
		newChannel.Reject(ssh.Prohibited, "remote host access denied for user")
		return
	}

	// Log the selection
	//if s.Selected != nil {
	//	if err := s.Selected(session, address); err != nil {
	//		newChannel.Reject(ssh.Prohibited, "access denied")
	//		return
	//	}
	//}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		newChannel.Reject(ssh.ConnectionFailed, fmt.Sprintf("error: %v", err))
		return
	}

	channel, reqs, err := newChannel.Accept()

	go ssh.DiscardRequests(reqs)
	var closer sync.Once
	closeFunc := func() {
		channel.Close()
		conn.Close()
	}

	go func() {
		io.Copy(channel, conn)
		closer.Do(closeFunc)
	}()

	go func() {
		io.Copy(conn, channel)
		closer.Do(closeFunc)
	}()
}

type rw struct {
	io.Reader
	io.Writer
}

func HostKeyCallback(hostID string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

// SessionForward performs a regular forward, providing the user with an
// interactive remote host selection if necessary. This forwarding type
// requires agent forwarding in order to work.
func (s *Server) SessionForward(session *Session, newChannel ssh.NewChannel, chans <-chan ssh.NewChannel) {

	// Okay, we're handling this as a regular session
	sesschan, sessReqs, err := newChannel.Accept()
	if err != nil {
		return
	}

	stderr := sesschan.Stderr()

	var remote api.Machine
	switch len(session.Machines) {
	case 0:
		fmt.Fprintf(stderr, "User has no permitted remote hosts.\r\n")
		sesschan.Close()
		return
	case 1:
		remote = session.Machines[0]
	default:
		comm := rw{Reader: sesschan, Writer: stderr}
		if s.Interactive == nil {
			remote, err = DefaultInteractive(comm, session)
		} else {
			remote, err = s.Interactive(comm, session)
		}
		if err != nil {
			sesschan.Close()
			return
		}
	}

	// Log the selection
	//if s.Selected != nil {
	//	if err = s.Selected(session, remote); err != nil {
	//		fmt.Fprintf(stderr, "Remote host selection denied.\r\n")
	//		sesschan.Close()
	//		return
	//	}
	//}

	fmt.Fprintf(stderr, "Connecting to %s@%s:%d\r\n", remote.Users[0].Username, remote.Ip, remote.Port)

	// Set up the agent

	_, agentReqs, err := session.Conn.OpenChannel("auth-agent@openssh.com", nil)
	//if err != nil {
	//	fmt.Fprintf(stderr, "sshmux requires either agent forwarding or secure channel forwarding.\r\n")
	//	fmt.Fprintf(stderr, "Either enable agent forwarding (-A), or use a ssh -W proxy command.\r\n")
	//	fmt.Fprintf(stderr, "For more info, see the Jumpserver's wiki.\r\n")
	//	sesschan.Close()
	//	return
	//}
	//defer agentChan.Close()
	go ssh.DiscardRequests(agentReqs)

	// Set up the client

	//ag := agent.NewClient(agentChan)

	singer, _ := remote.Signer()
	clientConfig := &ssh.ClientConfig{
		User: remote.Users[0].Username,
		Auth: []ssh.AuthMethod{
			//ssh.PublicKeysCallback(ag.Signers),
			ssh.PublicKeys(singer),
		},
		HostKeyCallback: HostKeyCallback,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", remote.Ip, remote.Port), clientConfig)
	if err != nil {
		fmt.Fprintf(stderr, "[gabriel]Connect failed: %v\r\n", err)
		sesschan.Close()
		return
	}

	// Handle all incoming channel requests
	go func() {
		for newChannel = range chans {
			if newChannel == nil {
				return
			}

			channel2, reqs2, err := client.OpenChannel(newChannel.ChannelType(), newChannel.ExtraData())
			if err != nil {
				x, ok := err.(*ssh.OpenChannelError)
				if ok {
					newChannel.Reject(x.Reason, x.Message)
				} else {
					newChannel.Reject(ssh.Prohibited, "remote server denied channel request")
				}
				continue
			}

			channel, reqs, err := newChannel.Accept()
			if err != nil {
				channel2.Close()
				continue
			}
			go proxy(reqs, reqs2, channel, channel2)
		}
	}()

	// Forward the session channel
	channel2, reqs2, err := client.OpenChannel("session", []byte{})
	if err != nil {
		fmt.Fprintf(stderr, "Remote session setup failed: %v\r\n", err)
		sesschan.Close()
		return
	}

	// Proxy the channel and its requests
	maskedReqs := make(chan *ssh.Request, 1)
	go func() {
		for req := range sessReqs {
			if req.Type == "auth-agent-req@openssh.com" {
				continue
			}
			maskedReqs <- req
		}
	}()
	proxy(maskedReqs, reqs2, sesschan, channel2)

}
