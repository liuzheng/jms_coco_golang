package client

import (
	//"os"
	//"syscall"
	//"os/signal"
	//"bufio"
	"coco/util/log"
	"github.com/googollee/go-socket.io"
	"encoding/json"
)

type screen struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type hostinfo struct {
	Ip       string `json:"ip"`
	Port     int `json:"port"`
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
	Locale   string `json:"locale"`
	Screen   screen `json:"screen"`
}

func RdpWebSocket() (server *socketio.Server) {
	log.Info("RdpWebSocket", "start")
	server, err := socketio.NewServer(nil)
	if log.HandleErr("RdpWebSocket", err, "初始化错误") {
		return
	}
	//uri := "ws://127.0.0.1:9250/socket.io/"
	//if log.HandleErr("RdpWebSocket", err, "socket error") {
	//	return
	//}

	server.On("connection", func(so socketio.Socket) {
		log.Info("RdpWebSocket", "on connection")
		so.On("disconnect", func() {
			log.Debug("RdpWebSocket", "on disconnect")
			so.Emit("disconnect")
		})
		so.On("infos", func(msg hostinfo) {
			log.Debug("WS", "%v", msg)
			infos, _ := json.Marshal(msg)
			log.Debug("WS", "%v", infos)

			//client.Emit("infos", string(infos))
		})
		//so.On("mouse", func(x, y, button int, isPressed bool) {
		//	client.Emit("mouse", x, y, button, isPressed)
		//})
		//so.On("wheel", func(x, y, step int, isNegative, isHorizontal bool) {
		//	client.Emit("wheel", x, y, step, isNegative, isHorizontal)
		//})
		//so.On("scancode", func(code int, isPressed bool) {
		//	client.Emit("scancode", code, isPressed)
		//})
		//so.On("unicode", func(code int, isPressed bool) {
		//	client.Emit("unicode", code, isPressed)
		//})
		//so.On("connect", func() {
		//	client.Emit("connect")
		//})
		//so.On("close", func() {
		//	client.Emit("close")
		//})
		//so.On("bitmap", func(bitmap interface{}) {
		//	client.Emit("bitmap", bitmap)
		//})
		//so.On("error", func(error interface{}) {
		//	client.Emit("error", error)
		//})

	})
	server.On("error", func(so socketio.Socket, err error) {
		log.HandleErr("RdpWebSocket", err)
	})
	return server
}
//func main() {
//	opts := &socketio_client.Options{
//		Transport: "websocket",
//		Query:     make(map[string]string),
//	}
//	opts.Query["infos"] = "user"
//	//opts.Query["pwd"] = "pass"
//	uri := "http://127.0.0.1:9250/socket.io/"
//
//	client, err := socketio_client.NewClient(uri, opts)
//	if err != nil {
//		fmt.Printf("NewClient error:%v\n", err)
//		return
//	}
//
//	client.On("error", func() {
//		fmt.Printf("on error\n")
//	})
//	client.On("connection", func() {
//		fmt.Printf("on connect\n")
//	})
//	client.On("message", func(msg string) {
//		fmt.Printf("on message:%v\n", msg)
//	})
//	client.On("disconnection", func() {
//		fmt.Printf("on disconnect\n")
//	})
//	client.Emit("infos", `{"ip":"10.59.169.171","port":3389,"screen":{"width":769,"height":726},"domain":"","username":"Administrator","password":"Administrator","locale":"en-US"}`)
//
//	reader := bufio.NewReader(os.Stdin)
//	for {
//		data, _, _ := reader.ReadLine()
//		command := string(data)
//		client.Emit("message", command)
//		fmt.Printf("send message:%v\n", command)
//	}
//	// trap SIGINT to trigger a shutdown.
//	signals := make(chan os.Signal, 1)
//	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
//	MSGFOR:
//	for {
//		select {
//		case <-signals:
//			break MSGFOR
//		}
//	}
//
//}
