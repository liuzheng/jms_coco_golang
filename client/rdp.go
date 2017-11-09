package client

import (
	//"os"
	//"syscall"
	//"os/signal"
	//"bufio"
	"net/http"
	"io"
	//"coco/util/log"
	"encoding/json"
)

type hostinfo struct {
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *hostinfo) toString() (r string) {
	rb, _ := json.Marshal(h)
	return string(rb)
}

func Rdp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.Form.Get("token")
		res := hostinfo{
			Ip:       "127.0.0.1",
			Port:     3389,
			Username: "liuzheng",
			Password: "liuzheng",
			Domain:   "",
		}
		io.WriteString(w, res.toString())
	}
}
