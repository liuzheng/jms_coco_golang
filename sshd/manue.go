package sshd

import (
	"fmt"
	"coco/api"
)

type Manue struct {
	Conn    rw
	Session *Session
}

var (
	welcome = []string{
		"输入 \033[32mID\033[0m 直接登录 或 输入\033[32m部分 IP,主机名,备注\033[0m 进行搜索登录(如果唯一).",
		"输入 \033[32m/\033[0m + \033[32mIP, 主机名 or 备注 \033[0m搜索. 如: /ip",
		"输入 \033[32mP/p\033[0m 显示您有权限的主机.",
		"输入 \033[32mG/g\033[0m 显示您有权限的主机组.",
		"输入 \033[32mG/g\033[0m\033[0m + \033[32m组ID\033[0m 显示该组下主机. 如: g1",
		"输入 \033[32mE/e\033[0m 批量执行命令.(未完成)",
		"输入 \033[32mU/u\033[0m 批量上传文件.(未完成)",
		"输入 \033[32mD/d\033[0m 批量下载文件.(未完成)",
		"输入 \033[32mH/h\033[0m 帮助.",
		"输入 \033[32mQ/q\033[0m 退出.",
	}
)

func NewManue(conn rw, session *Session) (manue Manue, err error) {
	manue = Manue{Conn:conn, Session:session}
	return
}

func (m *Manue)Welcome() {
	for k, v := range welcome {
		fmt.Fprintf(m.Conn, "    %d) %s\r\n", k, v)
	}
}
func (m *Manue)MachineList() {
	count := 0
	remotes := []api.Machine{}
	for _, v := range m.Session.Machines {
		for _, u := range v.Users {
			fmt.Fprintf(m.Conn, "    [%d] %s@%s:%d  %s\r\n", count, u.Username, v.Ip, v.Port, v.Remark)
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
