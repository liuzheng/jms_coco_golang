package websocket

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
	"coco/api"
)

type HttpAPI struct {
	API *api.Server
}
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
	Id       int    `json:"id"`
	Children []Host `json:"children"`
}

type Host struct {
	Name    string   `json:"name"`
	Uuid    string   `json:"uuid"`
	Type    string   `json:"type"`
	Users   []string `json:"users"`
}

func (h *HttpAPI) Nav(w http.ResponseWriter, r *http.Request) {
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

func (h *HttpAPI) Checklogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		io.WriteString(w, `{"logined":true}`)
	} else if r.Method == "POST" {
		fmt.Println(r)
	}
}
func (h *HttpAPI) HostGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		MachineGroup, _ := h.API.GetGroupList()
		HG := []HostGroup{}
		for _, MG := range MachineGroup {
			Machines, _ := h.API.GetList("", MG.Gid)
			H := []Host{}
			for _, M := range Machines {
				users := []string{}
				for _, u := range M.Users {
					users = append(users, u.Username)
				}
				H = append(H, Host{
					Name:    M.Name,
					Uuid:    M.Uuid,
					Type:    M.Type,
					Users:   users,
				})
			}
			HG = append(HG, HostGroup{
				Name:     MG.Name,
				Id:       MG.Gid,
				Children: H,
			})
		}
		rb, _ := json.Marshal(HG)
		io.WriteString(w, string(rb))
	}
}
