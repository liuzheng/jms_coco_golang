package main

import (
	"./api"
	"./util"
	"log"
)

func main() {
	//获取配置信息
	as := api.New()
	util.Config(as)
	log.Print("向JMS注册自己")
	resp, _ := as.Register()
	log.Print(resp)
	log.Print("获取用户PUBKEY")
	key, _ := as.GetUserPubKey("root")
	log.Print(key)
	log.Print("获取用户TOKEN：test")
	res, _ := as.GetLoginToken("root", key.Ticket)
	log.Print(res)
	log.Print("获取用户服务器列表：test")
	mlist, _ := as.GetList()
	log.Print(mlist)
	log.Print("获取Real server登陆凭证")
	mcredit, _ := as.GetLoginCredit(1, 1)
	log.Print(mcredit)
	log.Print("获取是否有监控权限")
	pb, _ := as.CheckMonitorToken(1)
	log.Print(pb)
	log.Print("上报开启Session")
	resp, _ = as.ReportSession(1, 1, 2, 1)
	log.Print(resp)
}
