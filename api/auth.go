package api

import (
	"errors"
)

func GetUserPubKey(username string) (pubkey string, err error) {
	if (username == "test") {
		return "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDZ0aTNRcRPOvJPN+kAsewJRCXiCGs7F9jELHaJuykwYwlV4/FuoQ+NK/4Q8BCLgF0gSZFFanv5WdGndxyDGJLAsevitCHKeWIivN2y48kdKFmrv1P+gIfBpo/rVVG+19wVKdXD4h+f2ZZhV2/edrEJb8xq6odsYwQ6VOLoaOai4MfkEmXnP5udef5lJ9NZD80qgMieqUW178+EN4b3eOPam0CT/DniYz7acJsvPv7z7DIvACzfSukf7HjtolKWrJp30Ypj4oNe9kpN6T6i9zKc+I21YTqwwGX3rxVuv+E7oU7AlQqxVfY526Xrmt1naKj+ufnOxIcZZ7X+zXLuxfVB xRain@xRain-MacBook.local", nil
	} else {
		return "", errors.New("error username")
	}
}

func CheckLoginToken(token string) (result bool, err error) {
	if (token == "testtoken") {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckMonitorToken(sessionid string, token string) (result bool, err error) {
	if (token == "testtoken" && sessionid == "s1") {
		return true, nil
	} else {
		return false, nil
	}
}
