package main

import (
	"./api"
	"log"
)

func main() {
	log.Print("获取用户PUBKEY")
	key, _ := api.GetUserPubKey("root")
	log.Print(key)
	log.Print("获取用户TOKEN：test")
	res, _ := api.GetLoginToken("root", key.Ticket)
	log.Print(res)
	log.Print("获取用户服务器列表：test")
	mlist, _ := api.GetList()
	log.Print(mlist)
	log.Print("获取Real server登陆凭证")
	mcredit, _ := api.GetLoginCredit(1, 1)
	log.Print(mcredit)
	log.Print("检测登陆凭证有效性")
	pb, _ := api.CheckUserToken()
	log.Print(pb)
	log.Print("获取是否有监控权限")
	pb, _ = api.CheckMonitorToken(1)
	log.Print(pb)
}
