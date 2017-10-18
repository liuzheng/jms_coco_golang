package main

import (
	"coco/api"
	"coco/util/log"
	"coco/sshd"
	"coco/websocket"
	"coco/util"
	"fmt"
)

func main() {
	util.Initial()
	log.Info("BOOT", "1-系统配置加载")
	util.CheckConfig()
	log.Info("BOOT", "2-初始化API模块")
	as := api.New()
	log.Info("BOOT", "3-向JMS服务器注册")
	as.Register()
	log.Info("BOOT", "4-启动SSH服务器")
	go sshd.Run()
	log.Info("BOOT", "5-启动WebSocket服务器")
	go websocket.Run()
	log.Info("BOOT", "6-系统启动完成")

	log.Info("main", "向JMS注册自己")
	resp, _ := as.Register()
	log.Info("main", "%v", resp)
	log.Info("main", "获取用户PUBKEY")
	key, _ := as.GetUserPubKey("root")
	log.Info("main", "%v", key)
	log.Info("main", "获取用户TOKEN")
	res, _ := as.GetLoginToken("root", "sdflkdjflsdjflk")
	log.Info("main", "%v", res)
	log.Info("main", "获取用户服务器列表")
	mlist, _ := as.GetList("", 0)
	log.Info("main", "%v", mlist)
	log.Info("main", "获取用户服务器组列表")
	mglist, _ := as.GetGroupList()
	log.Info("main", "%v", mglist)
	log.Info("main", "获取Real server登陆凭证")
	mcredit, _ := as.GetLoginCredit(1, 1)
	log.Info("main", "%v", mcredit)
	log.Info("main", "获取是否有监控权限")
	pb, _ := as.CheckMonitorToken(1)
	log.Info("main", "%v", pb)
	log.Info("main", "上报开启Session")
	resp, _ = as.ReportSession("UUID-1", 1, 2, 1)
	log.Info("main", "%v", resp)
	log.Info("main", "上报关闭Session")
	resp, _ = as.ReportSessionClose("UUID-2")
	log.Info("main", "%v", resp)
	for {
		fmt.Scanln()
	}
}
