package util

import (
	"flag"
	"net"
	"log"
	"../api"
)

func Config(apiServer *api.Server) {
	appId := flag.String("appid", "", "Jumpserver Core中添加完Coco后获得的AppId")
	jmsUrl := flag.String("apiurl", "", "Jumpserver Core地址")
	appKey := flag.String("appkey", "", "Jumpserver Core中添加完Coco后获得的AppKey")
	wsPort := flag.Int("wsport", 7871, "对外提供WS服务的地址")
	sshPort := flag.Int("sshport", 7822, "对外提供SSH服务的地址")
	ip := flag.String("ip", GetIp(), "对外提供服务的IP地址")
	flag.Parse()
	apiServer.Url = *jmsUrl
	apiServer.AppId = *appId
	apiServer.AppKey = *appKey
	apiServer.Ip = *ip
	apiServer.WsPort = *wsPort
	apiServer.SshPort = *sshPort
}

func GetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("无法获取本机IP，请使用-ip参数指定IP")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String();
			}
		}
	}
	return "127.0.0.1"
}
