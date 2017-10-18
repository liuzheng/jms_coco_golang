package sshd

import (
	"fmt"
	"io"
	"strconv"
	"coco/api"
	"coco/util/log"
	"coco/util/errors"
)

// DefaultInteractive is the default server selection prompt for users during
// session forward.
func DefaultInteractive(comm io.ReadWriter, session *Session) (api.Machine, errors.Error) {
	//remotes := session.Remotes
	count := 0

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
		fmt.Fprint(comm, "Please select remote server: ")
		var buf []byte
		b := make([]byte, 1)
		var (
			n int
			err error
		)
		for {
			if log.HandleErr("DefaultInteractive", err) {
				return api.Machine{}, errors.New(err.Error(), 200)
			}
			n, err = comm.Read(b)
			if n >= 0 {
				fmt.Fprintf(comm, "%s", b)
				switch b[0] {
				case '\r':
					fmt.Fprintln(comm, "")
					res, err := strconv.ParseInt(string(buf), 10, 64)
					if log.HandleErr("DefaultInteractive", err) {
						fmt.Fprintln(comm, "input not a valid integer. Please try again")
						continue loop
					}
					if int(res) >= count || res < 0 {
						fmt.Fprintln(comm, "No such server. Please try again")
						continue loop
					}
					return remotes[int(res)], nil
				case 0x03:
					fmt.Fprintln(comm, "\r\nGoodbye")
					return api.Machine{}, errors.New("user terminated session", 400)
				}
				buf = append(buf, b[0])
			}
		}
	}
}
