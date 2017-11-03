package client

import (
	"fmt"
	"os"
	"syscall"
	"os/signal"
	"github.com/zhouhui8915/go-socket.io-client"
	"bufio"
)

func main() {
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["infos"] = "user"
	//opts.Query["pwd"] = "pass"
	uri := "http://127.0.0.1:9250/socket.io/"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		fmt.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func() {
		fmt.Printf("on error\n")
	})
	client.On("connection", func() {
		fmt.Printf("on connect\n")
	})
	client.On("message", func(msg string) {
		fmt.Printf("on message:%v\n", msg)
	})
	client.On("disconnection", func() {
		fmt.Printf("on disconnect\n")
	})
	client.Emit("infos", `{"ip":"10.59.169.171","port":3389,"screen":{"width":769,"height":726},"domain":"","username":"Administrator","password":"Administrator","locale":"en-US"}`)

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		client.Emit("message", command)
		fmt.Printf("send message:%v\n", command)
	}
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
MSGFOR:
	for {
		select {
		case <-signals:
			break MSGFOR
		}
	}

}
