package api

import (
	"errors"
)

//根据用户名获取pubkey
func GetUserPubKey(username string) (UserPubKey, error) {
	if (username == "root") {
		return UserPubKey{
			Ticket: "4b4a3ea8ceb95ae4db7d5f3e333109fe42c",
			Key:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDALyKgUwHamftJVQyl9jJDJiINo+omlQNyrB0jpJn9wXjapO8VDf8rjYhSzCBovoswq3NFwcEXzIcZc38nnZbWJwPCL7wsXapP1qiE65hgEvEYuIwpIystNip2pd+H4rrDEIboWY5rA9VoGexC5GSKHwwqw+hidQPtlIBEjugQR1eQQTIuhQQX2H3zrergybwAWlX3jPcxpp4JguH/6QrgkybTFkrHdCOjt8d7LjduiH/AbCIn9C6bcp3sJdN2f+wxP/aVFBJ9G1I/pToY8dBzmeBAsuVO8aYBy7VmfnD3hGJPfQMhHU+d6fzZfVURaSigNrdcqVlo7w7rMsE5HGvB root@bab0e381d0d8",
		}, nil
	} else {
		return UserPubKey{}, errors.New("error username")
	}
}

//根据pubkey和username获取登陆TOKEN
func GetLoginToken(username string, ticket string) (UserToken, error) {
	return UserToken{
		Token:   "4b4a3ea8ceb95ae4db7d5f3e109fe42c",
		Expired: 1507675673,
	}, nil
}

//检查用户TOKEN有效性
func CheckUserToken() (bool, error) {
	return true, nil
}

//检查用户能否开启监控SHELL
func CheckMonitorToken(sessionid int) (bool, error) {
	return true, nil
}
