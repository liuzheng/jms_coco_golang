package util

import (
	"flag"
	"net"
	"coco/util/log"
)

const Version = "coco-0.5"

var (
	Ip      = flag.String("ip", GetIp(), "对外提供服务的IP地址")
	AppId   = flag.String("appid", "", "Jumpserver Core中添加完Coco后获得的AppId")
	JmsUrl  = flag.String("jmsurl", "", "Jumpserver Core地址")
	AppKey  = flag.String("appkey", "", "Jumpserver Core中添加完Coco后获得的AppKey")
	WsPort  = flag.Int("wsport", 7871, "对外提供WS服务的地址")
	SshPort = flag.Int("sshport", 7822, "对外提供SSH服务的地址")
	Hostkey = flag.String("hostkey", "host_key", "SSH Server的私钥文件")
)

func CheckConfig() {
	if *JmsUrl == "" {
		log.Fatal("Config", "JmsUrl 未指定，Jumpserver Core地址 -jmsurl")
	} else {
		log.Info("Config", "读取配置 JmsUrl：%s", *JmsUrl)
	}
	if *AppId == "" {
		log.Fatal("Config", "AppId 未指定，Jumpserver Core中添加完Coco后获得的AppId -appid")
	} else {
		log.Info("Config", "读取配置 AppId：%s", *AppId)
	}
	if *AppKey == "" {
		log.Fatal("Config", "AppKey 未指定，Jumpserver Core中添加完Coco后获得的AppKey -appkey")
	} else {
		log.Info("Config", "读取配置 AppKey：%s", log.Password(*AppKey))
	}
	log.Info("Config", "读取配置 HostKeyFile：%s", *Hostkey)
	log.Info("Config", "读取配置 服务Ip：%s", *Ip)
	log.Info("Config", "读取配置 WsPort：%d", *WsPort)
	log.Info("Config", "读取配置 SshPort：%d", *SshPort)
}

func GetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("GetIp", "无法获取本机IP，请使用-ip参数指定IP")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
