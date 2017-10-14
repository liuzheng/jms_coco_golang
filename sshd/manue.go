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
	welcome = map[int]string{
		1:"输入 ID",
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
