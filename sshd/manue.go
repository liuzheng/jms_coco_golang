package sshd

import (
	"fmt"
	"strings"
	"coco/api"
	"coco/util/errors"
	"coco/util/log"
	"bytes"
)

type Manue struct {
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
func NewMenu(conn rw, session *Session) (manue Manue, err error) {
	as := api.New()
	if !as.Login(session.User.Name) {
		return Manue{}, errors.New("Login fail", 403)
	}
	manue = Manue{Conn: conn, Session: session, API: as}
	return
}

// Menu manager
func (m *Manue) Manager() {
	for {
		command, err := m.Command()
		if log.HandleErr("Manue Manager", err, "logout") {
			return
		}
		log.Debug("Manue Manager", "the command is: %s", command)
		if len(command) == 1 {
			fmt.Fprint(m.Conn, "\r\n")
			switch strings.ToUpper(string(command[0])) {
			case "P":
				// 输入 P/p 显示您有权限的主机.
				m.GetMachineList()
			case "G":
				// 输入 G/g 显示您有权限的主机组.
				m.GetHostGroup()
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
				fmt.Fprint(m.Conn, "Goodbye\r\n")
				return
			default:
				fmt.Fprint(m.Conn, "TO BE CONTINUED")
			}
		} else if len(command) > 1 {
			fmt.Fprint(m.Conn, "\r\n")
			switch strings.ToUpper(string(command[0])) {
			case "/":
				// 输入 / + IP, 主机名 or 备注 搜索. 如: /ip
				m.Search(command[1:])
			case "G":
				//  输入 G/g + 组ID 显示该组下主机. 如: g1
				m.GetHostGroupList(command[1:])
			default:
				// 输入 ID 直接登录 或 输入部分 IP,主机名,备注 进行搜索登录(如果唯一).
			}
		}
	}
	return
}

// Command: got the command user input, just the Opt menu
func (m *Manue) Command() (command string, err errors.Error) {
	fmt.Fprint(m.Conn, "\r\nOpt>")
	var left, right []byte
	b := make([]byte, 3)
	for {
		n, err := m.Conn.Read(b)
		if err != nil {
			log.Error("Server", "%v", err)
			fmt.Fprint(m.Conn, "^D")
			fmt.Fprint(m.Conn, "\r\nGoodbye\r\n")
			return "Exit", nil
		}
		if n >= 0 {
			log.Debug("Opt loop", "the key %v", b)
			switch {
			case 0x0A == b[0] || 0x0D == b[0]: // 换行键 or 回车键
				return string(append(left, right...)), nil
			case 0x03 == b[0]: // ctrl-c
				fmt.Fprint(m.Conn, "^C")
				return "", nil
			case 0x04 == b[0]: // ctrl-d
				fmt.Fprint(m.Conn, "^D\r\nGoodbye\r\n")
				return "Exit", nil
			case 0x7f == b[0] || 0x08 == b[0]: // delete or backspace
				if l := len(left); l > 0 {
					left = left[:l-1]
					fmt.Fprintf(m.Conn, "\b \b%s ", string(right))
					for i := -1; i < len(right); i++ {
						fmt.Fprintf(m.Conn, "%s", []byte{27, 91, 68})
					}
				}
				continue
			case 0x20 <= b[0] && b[0] <= 0x7E:
				fmt.Fprintf(m.Conn, "%s%s", b, string(right))
				for i := 0; i < len(right); i++ {
					fmt.Fprintf(m.Conn, "%s", []byte{27, 91, 68})
				}
			case bytes.Compare([]byte{27, 91, 65}, b) == 0: // 方向键(↑)
				log.Debug("Menu Command", "方向键(↑)")
			case bytes.Compare([]byte{27, 91, 66}, b) == 0: // 方向键(↓)
				log.Debug("Menu Command", "方向键(↓)")
			case bytes.Compare([]byte{27, 91, 67}, b) == 0: // 方向键(→)
				log.Debug("Menu Command", "方向键(→)")
				if len(right) > 0 {
					fmt.Fprintf(m.Conn, "%s", b)
					left = append(left, right[0])
					right = right[1:]
				}
			case bytes.Compare([]byte{27, 91, 68}, b) == 0: // 方向键(←)
				log.Debug("Menu Command", "方向键(←)")
				if len(left) > 0 {
					fmt.Fprintf(m.Conn, "%s", b)
					right = append([]byte{left[len(left)-1]}, right...)
					left = left[:len(left)-1]
				}
			case bytes.Compare([]byte{27, 0, 0}, b) == 0: // esc
				log.Debug("Menu Command", "ESC")
			default:
				fmt.Fprintf(m.Conn, "%s", "")
				continue
			}
			if 0x1b != b[0] {
				left = append(left, b[0])
			}
			b = make([]byte, 3)
			log.Debug("Menu Command", " %v %v", left, right)
		}
	}
}

// 欢迎页
func (m *Manue) Welcome() {
	fmt.Fprintf(m.Conn, "\033[1;32m  %s, 欢迎使用Jumpserver开源跳板机系统  \033[0m\r\n", m.Session.Conn.User())
	m.GetHelp()
}

// 获取主机列表
func (m *Manue) GetMachineList() {
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
func (m *Manue) GetHostGroup() {

}

// 帮助页
func (m *Manue) GetHelp() {
	for k, v := range help {
		fmt.Fprintf(m.Conn, "    %d) %s\r\n", k, v)
	}
}

// 获取主机组列表
func (m *Manue) GetHostGroupList(id string) {
	log.Debug("Menu GetHostGroupList", "%v", id)

}

// 搜索主机
func (m *Manue) Search(q string) {
	log.Debug("Menu Search", "%v", q)
}
