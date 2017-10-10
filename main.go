package main

import (
	"./api"
	"log"
)

func main() {
	log.Print("获取用户PUBKEY")
	key, _ := api.GetUserPubKey("test")
	log.Print(key)
	log.Print("LOGIN 验证用户TOKEN：testtoken")
	res, _ := api.CheckLoginToken("testtoken")
	log.Print(res)
	log.Print("Monitor 验证用户TOKEN：testtoken")
	res, _ = api.CheckMonitorToken("s1", "testtoken")
	log.Print(res)
}
