package sshd

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"coco/api"
)

// DefaultInteractive is the default server selection prompt for users during
// session forward.
func DefaultInteractive(comm io.ReadWriter, session *Session) (api.Machine, error) {
	//remotes := session.Remotes
	count := 0
	fmt.Fprintf(comm, "%s, 欢迎使用Jumpserver开源跳板机系统\r\n", session.Conn.User())
	remotes := []api.Machine{}
	for _, v := range session.Machines {
		for _, u := range v.Users {
			fmt.Fprintf(comm, "    [%d] %s@%s:%d  %s\r\n", count, u.Username, v.Ip, v.Port, v.Remark)
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

	// Beware, nasty input parsing loop
loop:
	for {
		fmt.Fprintf(comm, "Please select remote server: ")
		var buf []byte
		b := make([]byte, 1)
		var (
			n   int
			err error
		)
		for {
			if err != nil {
				return api.Machine{}, err
			}
			n, err = comm.Read(b)
			if n >= 0 {
				fmt.Fprintf(comm, "%s", b)
				switch b[0] {
				case '\r':
					fmt.Fprintf(comm, "\r\n")
					res, err := strconv.ParseInt(string(buf), 10, 64)
					if err != nil {
						fmt.Fprintf(comm, "input not a valid integer. Please try again\r\n")
						continue loop
					}
					if int(res) >= count || res < 0 {
						fmt.Fprintf(comm, "No such server. Please try again\r\n")
						continue loop
					}

					return remotes[int(res)], nil
				case 0x03:
					fmt.Fprintf(comm, "\r\nGoodbye\r\n")
					return api.Machine{}, errors.New("user terminated session")
				}

				buf = append(buf, b[0])
			}
		}
	}
}
