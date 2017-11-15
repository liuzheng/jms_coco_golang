package websocket

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
)

type navlist struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Children []navChild `json:"children"`
}
type navChild struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Click   string `json:"click"`
	Href    string `json:"href"`
	Disable bool   `json:"disable"`
}

type HostGroup struct {
	Name     string `json:"name"`
	Id       string `json:"id"`
	Children []Host `json:"children"`
}

type Host struct {
	Name    string `json:"name"`
	Uuid    string `json:"uuid"`
	Type    string `json:"type"`
	Token   string `json:"token"`
	Machine string `json:"machine"`
}

func Nav(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		res := []navlist{
			{Id: "File",
				Name: "Server",
				Children: []navChild{
					{
						Id:      "NewConnection",
						Href:    "Aaaa",
						Name:    "New connection",
						Disable: true,
					},
					{Id: "Connect",
						Href: "Aaaa",
						Name: "Connect",
						Disable: true,
					},
					{
						Id:    "Disconnect",
						Click: "Disconnect",
						Name:  "Disconnect",
					},
					{
						Id:    "DisconnectAll",
						Click: "DisconnectAll",
						Name:  "Disconnect all",
					},
					{
						Id:      "Duplicate",
						Href:    "Aaaa",
						Name:    "Duplicate",
						Disable: true,
					},
					{
						Id:      "Upload",
						Href:    "Aaaa",
						Name:    "Upload",
						Disable: true,
					},
					{
						Id:      "Download",
						Href:    "Aaaa",
						Name:    "Download",
						Disable: true,
					},
					{
						Id:      " Search",
						Href:    "Aaaa",
						Name:    "Search",
						Disable: true,
					},
					{
						Id:    "Reload",
						Click: "ReloadLeftbar",
						Name:  "Reload",
					},
				},},
			{
				Id:   "View",
				Name: "View",
				Children: []navChild{
					{
						Id:    "HindLeftManager",
						Click: "HideLeft",
						Name:  "Hind left manager",
					},
					{
						Id:      "SplitVertical",
						Href:    "Aaaa",
						Name:    "Split vertical",
						Disable: true,
					},
					{
						Id:      "CommandBar",
						Href:    "Aaaa",
						Name:    "Command bar",
						Disable: true,
					},
					{
						Id:      "ShareSession",
						Href:    "Aaaa",
						Name:    "Share session (read/write)",
						Disable: true,
					},
					{
						Id:      "Language",
						Href:    "Aaaa",
						Name:    "Language",
						Disable: true,
					}},
			}, {
				Id:   "Help",
				Name: "Help",
				Children: []navChild{
					{
						Id:    "EnterLicense",
						Click: "EnterLicense",
						Name:  "Enter License",
					},
					{
						Id:    "Website",
						Click: "Website",
						Name:  "Website",
					},
					{
						Id:    "BBS",
						Click: "BBS",
						Name:  "BBS",
					}},
			},
		}
		rb, _ := json.Marshal(res)
		io.WriteString(w, string(rb))
	}
}

func Checklogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		io.WriteString(w, `{"logined":true}`)
	} else if r.Method == "POST" {
		fmt.Println(r)
	}
}
func HostGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		res := []HostGroup{
			{
				Name: "ops",
				Id:   "ccc",
				Children: []Host{
					{
						Name:  "ops-linux",
						Uuid:  "xxxx",
						Type:  "ssh",
						Token: "sshxxx",
					},
					{
						Name:    "ops-win",
						Uuid:    "win-aasdf",
						Type:    "rdp",
						Token:   "rdpxxx",
						Machine: "sss",
					},
				},
			},
		}
		rb, _ := json.Marshal(res)
		io.WriteString(w, string(rb))
	}
}
