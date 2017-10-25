package sshd

import (
	"fmt"
	"strings"
	"coco/api"
	"coco/util/errors"
	"coco/util/log"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"strconv"
)

type Menu struct {
	Conn    rw
	Session *Session
	api     *api.Server
	term    *terminal.Terminal
	gid     int
}

var (
	help = []string{
		"输入 \033[32mID\033[0m 直接登录 或 输入\033[32m部分 IP,主机名,备注\033[0m 进行搜索登录(如果唯一).",
		"输入 \033[32m/\033[0m + \033[32mIP, 主机名 or 备注 \033[0m搜索. 如: /ip",
		"输入 \033[32mP/p\033[0m 显示您有权限的主机.",
		"输入 \033[32mG/g\033[0m 显示您有权限的主机组.",
		"输入 \033[32mG/g\033[0m\033[0m + \033[32m组ID\033[0m 显示该组下主机. 如: g1",
		"输入 \033[32mE/e\033[0m 批量执行命令.(未完成)", //TODO: 暂时不管这个功能
		"输入 \033[32mU/u\033[0m 批量上传文件.(未完成)", //TODO: 暂时不管这个功能
		"输入 \033[32mD/d\033[0m 批量下载文件.(未完成)", //TODO: 暂时不管这个功能
		"输入 \033[32mH/h\033[0m 帮助.",
		"输入 \033[32mQ/q\033[0m 退出.",
	}
)
// 初始化menu
func NewMenu(conn rw, session *Session) (menu Menu, err errors.Error) {
	menu = Menu{
		Conn:    conn,
		Session: session,
		api:     session.User.Api,
		term:    terminal.NewTerminal(conn, "Opt[0]> "),
		gid:     0}
	//menu.term.SetSize(140, 40)
	return
}

// Menu manager
func (m *Menu) Manager() {
	defer fmt.Fprint(m.Conn, "Goodbye\r\n")
	for {
		command, err := m.term.ReadLine()
		if err == io.EOF {
			return
		} else if log.HandleErr("Menu Manager", err, "logout") {
			return
		}
		log.Debug("Menu Manager", "the command is: %s", command)
		if len(command) == 1 {
			switch strings.ToUpper(string(command[0])) {
			case "P":
				// 输入 P/p 显示您有权限的主机.
				m.Search("")
			case "G":
				// 输入 G/g 显示您有权限的主机组.
				m.GetHostGroups()
				//case "E":
				////  输入 E/e 批量执行命令.(未完成)
				//case "U":
				////  输入 U/u 批量上传文件.(未完成)
				//case "D":
				////  输入 D/d 批量下载文件.(未完成)
			case "H":
				//  输入 H/h 帮助.
				m.GetHelp()
			case "Q":
				//  输入 Q/q 退出.
				return
			default:
				fmt.Fprint(m.Conn, "TO BE CONTINUED\r\n")
			}
		} else if len(command) > 1 {
			switch strings.ToUpper(string(command[0])) {
			case "/":
				// 输入 / + IP, 主机名 or 备注 搜索. 如: /ip
				m.Search(command[1:])
			case "G":
				//  输入 G/g + 组ID 显示该组下主机. 如: g1
				m.gid, _ = strconv.Atoi(command[1:])
				m.term.SetPrompt("Opt[" + strconv.Itoa(m.gid) + "]> ")
				m.GetHostGroupList()
			case "E":
				if "xit" == command[1:] {
					return
				}
			default:
				// 输入 ID 直接登录 或 输入部分 IP,主机名,备注 进行搜索登录(如果唯一).
			}
		}
	}
	return
}

// 欢迎页
func (m *Menu) Welcome() {
	fmt.Fprintf(m.Conn, "\033[1;32m  %s, 欢迎使用Jumpserver开源跳板机系统  \033[0m\r\n", m.Session.Conn.User())
	m.GetHelp()
}

// 获取主机组列表
func (m *Menu) GetHostGroups() (err errors.Error) {
	MachineGroup, err := m.api.GetGroupList()
	if log.HandleErr("GetHostGroup", err, "") {
		return err
	}
	format := "[%-4d]\t%-16s\t%s\r\n"
	fmt.Fprintf(m.Conn, "[%-4s]\t%-16s\t%s\r\n", "GID", "Name", "Comment")

	for _, v := range MachineGroup {
		fmt.Fprintf(m.Conn, format, v.Gid, v.Name, v.Remark)
	}

	return
}

// 帮助页
func (m *Menu) GetHelp() {
	for k, v := range help {
		fmt.Fprintf(m.Conn, "    %d) %s\r\n", k, v)
	}
}

// 获取主机组列表
func (m *Menu) GetHostGroupList() (err errors.Error) {
	Machine, errs := m.api.GetList("", m.gid)
	if log.HandleErr("GetHostGroupList", errs, "") {
		return errs
	}
	m.printMachines(Machine)

	return
}

// 搜索主机
func (m *Menu) Search(q string) (err errors.Error) {
	log.Debug("Menu Search", "%v", q)
	m.Session.Machines, err = m.api.GetList(q, m.gid)
	if log.HandleErr("GetMachineList", err, "") {
		return err
	}
	m.printMachines(m.Session.Machines)
	return
}

func (m *Menu) printMachines(Machines []api.Machine) {
	format := "[%-4d]\t%-16s\t%-5d\t%s\t%-14s\t%s\r\n"
	fmt.Fprintf(m.Conn, "[%-4s]\t%-16s\t%-5s\t%s\t%-14s\t%s\r\n", "ID", "IP", "Port", "Hostname", "Username", "Comment")
	for _, v := range Machines {
		for _, u := range v.Users {
			fmt.Fprintf(m.Conn, format, v.Sid, v.Ip, v.Port, v.Name, u.Username, v.Remark)
			//remotes = append(remotes, api.Machine{
			//	Ip:     v.Ip,
			//	Port:   v.Port,
			//	Name:   v.Name,
			//	Sid:    v.Sid,
			//	Remark: v.Remark,
			//	Users:  []api.MachineUser{u},
			//})
		}
	}
	return
}
