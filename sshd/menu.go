package sshd

import (
	"fmt"
	"coco/api"
	"coco/util/errors"
)

type Menu struct {
	Conn    rw
	Session *Session
	API     *api.Server
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
func NewMenu(conn rw, session *Session) (menu Menu, err error) {
	as := api.New()
	if !as.Login(session.User.Name) {
		return Menu{}, errors.New("Login fail", 403)
	}
	menu = Menu{Conn: conn, Session: session, API: as}
	return
}

// 欢迎页
func (m *Menu) Welcome() {
	fmt.Fprintf(m.Conn, "\033[1;32m  %s, 欢迎使用Jumpserver开源跳板机系统  \033[0m\r\n", m.Session.Conn.User())
	m.GetHelp()
}

// 获取主机列表
func (m *Menu) GetMachineList() {
	count := 0
	remotes := []api.Machine{}
	format := "[%-4d]\t%-16s\t%-5d\t%s\t%s\t%s\r\n"
	fmt.Fprintf(m.Conn, "[%-4s]\t%-16s\t%-5s\t%s\t%s\t%s\r\n", "ID", "IP", "Port", "Hostname", "Username", "Comment")
	for _, v := range m.Session.Machines {
		for _, u := range v.Users {
			fmt.Fprintf(m.Conn, format, count, v.Ip, v.Port, "hostname", u.Username, v.Remark)
			remotes = append(remotes, api.Machine{
				Ip:     v.Ip,
				Port:   v.Port,
				Name:   v.Name,
				Sid:    v.Sid,
				Remark: v.Remark,
				Users:  []api.MachineUser{u},
			})
			count++
		}
	}

}

// 获取主机组内主机列表
func (m *Menu) GetHostGroup() {

}

// 帮助页
func (m *Menu) GetHelp() {
	for k, v := range help {
		fmt.Fprintf(m.Conn, "    %d) %s\r\n", k, v)
	}
}

// 获取主机组列表
func (m *Menu) GetHostGroupList(id string) {

}

// 搜索主机
func (m *Menu) Search(q string) {

}
