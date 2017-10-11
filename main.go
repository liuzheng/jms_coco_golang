package main

import (
	"./api"
	"./util"
	log "github.com/liuzheng712/golog"
)

func main() {
	log.Logs("", "DEBUG", "DEBUG")
	//获取配置信息
	as := api.New()
	util.Config(as)
	log.Info("main", "向JMS注册自己")
	resp, _ := as.Register()
	log.Debug("main", "%v", resp)
	log.Info("main", "获取用户PUBKEY")
	user, _ := as.Auth()

	key, _ := user.GetUserPubKey()
	log.Debug("main", "%v", key)
	log.Info("main", "获取用户TOKEN：test")
	res, _ := user.GetLoginToken()
	log.Debug("main", "%v", res)
	log.Info("main", "获取用户服务器列表：test")
	mlist, _ := as.GetList() // TODO: change to user
	log.Debug("main", "%v", mlist)
	log.Info("main", "获取Real server登陆凭证")
	mcredit, _ := as.GetLoginCredit(1, 1) // TODO: change to user
	log.Debug("main", "%v", mcredit)
	log.Info("main", "获取是否有监控权限")
	pb, _ := user.CheckMonitorToken(1)
	log.Debug("main", "%v", pb)
	log.Info("main", "上报开启Session")
	resp, _ = as.ReportSession(1, 1, 2, 1)
	log.Debug("main", "%v", resp)
}
