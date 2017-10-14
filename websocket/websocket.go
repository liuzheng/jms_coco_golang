package websocket

import (
	"golang.org/x/crypto/ssh"
	"os"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/googollee/go-socket.io"
	"io"
	"coco/api"
	"coco/client"
	"coco/util"
	"coco/util/log"
	"net/http"
	"fmt"
)

type TTY struct {
	Key       api.UserPubKey
	UserToken api.UserToken
	Machines  []api.Machine
}

func (t *TTY) GetTermSize() (termWidth, termHeight int, err error) {

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Error("GetTermSize", "创建文件描述符: %v", err)
		return
	}

	termWidth, termHeight, err = terminal.GetSize(fd)
	if err != nil {
		log.Error("GetTermSize", "获取窗口宽高: %v", err)
		return
	}
	defer terminal.Restore(fd, oldState)
	return

}

//func (t *TTY) GetTermSession(machineID string) (session *client.Session, err error) {
//
//	client, err := NewClient(machineID)
//	if err != nil && !client.Success {
//		log.Error("GetTermSession", "Failed to newClient: %v", err)
//		return
//	}
//	defer client.Close()
//	session.session, err = client.SSH.NewSession()
//	if err != nil {
//		log.Error("GetTermSession", "Failed to create session: %v", err)
//	}
//	defer session.session.Close()
//	return
//}
func (t *TTY) GetMachine(machineID string) (machine api.Machine, err error) {
	// TODO: get the machine info
	return
}

//func (t *TTY) GetHostUsername(kid string, hid string) string {
//
//	return
//}

func New() (server *socketio.Server) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("ServerInit", "%v", err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Info("ServerInit", "on connection")

		as := api.New()
		var session *client.Session
		var soin io.WriteCloser
		var soout io.Reader
		t := TTY{}

		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		buf := make([]byte, 10240)

		so.On("data", func(msg string) {
			log.Debug("ServerInit", "reselve %v", msg)
			soin.Write([]byte(msg))
		})
		so.On("login", func(username string) {
			t.Key, err = as.GetUserPubKey(username)
			log.HandleErr("sshd New", err)
			//if err != nil {
			//	log.Error("ServerInit", "%v", err)
			//}
			//t.UserToken, err = userauth.GetLoginToken()
			//if err != nil {
			//	log.Error("ServerInit", "%v", err)
			//}
			t.Machines, err = as.GetList()
			if err != nil {
				log.Error("ServerInit", "%v", err)
			}
			// TODO: so.Emit the machine list
		})
		so.On("machine", func(machineID string) {
			log.Debug("ServerInit", "try to login into %v", machineID)

			remote, err := t.GetMachine(machineID)
			if err != nil {
				log.Error("ServerInit", "%v", err)
				so.Emit("data", err.Error())
				so.Emit("disconnect")
				return
			}
			credit, err := as.GetLoginCredit(remote.Sid, remote.Users[0].Uid)
			if err != nil {
				log.Error("ServerInit", "GetLoginCredit : %v", err)
			}
			connect, err := client.New(remote, credit)
			session, err = connect.NewSession()
			soin, err = session.Session.StdinPipe()
			log.HandleErr("ServerInit", err)
			soout, err = session.Session.StdoutPipe()
			log.HandleErr("ServerInit", err)

			if err := session.Session.RequestPty("xterm", 24, 80, modes); err != nil {
				log.Error("ServerInit", "1request for pseudo terminal failed: %v", err)
				return
			}
			err = session.Session.Shell()
			if err != nil {
				log.Error("ServerInit", "执行Shell出错: %v", err)
				return
			}
			//err = session.Wait()
			//if err != nil {
			//	log.Error("执行Wait出错: ", err)
			//	return
			//}

			go func() {
				for {
					n, err := soout.Read(buf)
					if err == io.EOF {
						so.Emit("disconnect")
						log.Info("ServerInit", "websocket is disconnected")
						break
					}
					if n > 0 {
						so.Emit("data", string(buf[:n]))
					}
				}
			}()
			//go func() {
			//
			//	err = session.Wait()
			//	if err != nil {
			//		log.Error("执行Wait出错: ", err)
			//		return
			//	}
			//}()

		})
		so.On("resize", func(size []int) {
			log.Debug("ServerInit", "resize to: %v", size)
			if err := session.Resize(size[1], size[0]); err != nil {
				log.Error("ServerInit", "request for pseudo terminal failed: %v", err)
				return
			}
			//so.Emit("data", buffString())
		})
		so.On("disconnect", func() {
			log.Debug("ServerInit", "on disconnect")
			session.Close()

			so.Emit("disconnect")
		})
		//so.On("key", func(key string) {
		//	//fmt.Println(key)
		//	var timestring string
		//	var check string
		//	var check1 string
		//	var timestamp int64
		//	list := strings.Split(key, "\n")
		//	if len(list) > 4 {
		//		for _, i := range list {
		//			if len(i) == 32 {
		//				key1 := fmt.Sprintf("%x", md5.Sum([]byte(i)))
		//				timestring += string(key1[0]) + string(key1[31])
		//				check += string(i[31])
		//				check1 = string(key1[28:32])
		//				if len(timestring) == 8 {
		//					timestamp, _ = strconv.ParseInt(timestring, 16, 32)
		//				}
		//			} else {
		//				check = "zzzzzz"
		//				check1 = "aaaa"
		//				break
		//			}
		//		}
		//		if check[0:4] == check1 {
		//			tm := time.Unix(timestamp, 0)
		//			so.Emit("popup", "License until "+tm.Format("2006-01-02"))
		//		} else {
		//			so.Emit("popup", "Error key")
		//		}
		//	} else {
		//		so.Emit("popup", "Error key")
		//	}
		//})
		//so.On("api", func(msg string) {
		//	log.Debug("ServerInit", "api: %v", msg)
		//	nav := `[{"id":"File","name":"Server","children":[{"id":"NewConnection","href":"Aaaa","name":"New connection","disable":true},{"id":"Connect","href":"Aaaa","name":"Connect","disable":true},{"id":"Disconnect","click":"Disconnect","name":"Disconnect"},{"id":"DisconnectAll","click":"DisconnectAll","name":"Disconnect all"},{"id":"Duplicate","href":"Aaaa","name":"Duplicate","disable":true},{"id":"Upload","href":"Aaaa","name":"Upload","disable":true},{"id":"Download","href":"Aaaa","name":"Download","disable":true},{"id":"Search","href":"Aaaa","name":"Search","disable":true},{"id":"Reload","click":"ReloadLeftbar","name":"Reload"}]},{"id":"View","name":"View","children":[{"id":"HindLeftManager","click":"HideLeft","name":"Hind left manager"},{"id":"SplitVertical","href":"Aaaa","name":"Split vertical","disable":true},{"id":"CommandBar","href":"Aaaa","name":"Command bar","disable":true},{"id":"ShareSession","href":"Aaaa","name":"Share session (read/write)","disable":true},{"id":"Language","href":"Aaaa","name":"Language","disable":true}]},{"id":"Edit","name":"Edit","children":[{"id":"Host","name":"Host","href":"HostEdit"}]},{"id":"Help","name":"Help","children":[{"id":"EnterLicense","click":"EnterLicense","name":"Enter License"},{"id":"Website","click":"Website","name":"Website"},{"id":"BBS","click":"BBS","name":"BBS"}]}]`
		//	//leftbar := `[{"title":"xxxx","key":"1","folder":false,"machine":"localhost"},{"title":"Folder 2","key":"2","folder":true,"children":[{"title":"Node 2.1","key":"3","machine":"aa"},{"title":"Node 2.2","key":"4","machine":"bb"}]},{"title":"Folder 3","key":"2","folder":true,"children":[{"title":"Node 2.1","key":"3"},{"title":"Node 2.2","key":"4"}]}]`
		//	if msg == "nav" {
		//		so.Emit("nav", string(nav))
		//		//} else if msg == "leftbar" {
		//		//	so.Emit("leftbar", string(leftbar))
		//	} else if msg == "leftbar" {
		//		so.Emit("leftbar", "changed")
		//	} else if msg == "all" {
		//		so.Emit("nav", string(nav))
		//		//so.Emit("leftbar", string(leftbar))
		//	}
		//
		//})

	})
	server.On("error", func(so socketio.Socket, err error) {
		log.HandleErr("ServerInit", err)
	})
	return server
}

func Run() {
	server := New()
	http.Handle("/socket.io/", server)
	log.Fatal("WS Run", "%v", http.ListenAndServe(fmt.Sprintf("%s:%d", *util.Ip, *util.WsPort), nil))
}
