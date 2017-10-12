package api

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"coco/util"
)

//初始化一个ApiServer
func New() *Server {
	server := Server{
		Url:     *util.JmsUrl,
		AppId:   *util.AppId,
		AppKey:  *util.AppKey,
		Ip:      *util.Ip,
		WsPort:  *util.WsPort,
		SshPort: *util.SshPort,
	}
	server.Action = Action{
		GetUserPubKey:     "mock.php?act=getpubkey",
		GetUserToken:      "mock.php?act=getusertoken",
		CheckMonitorToken: "mock.php?act=checkmonitortoken",
		GetMachineList:    "mock.php?act=machines",
		GetLoginCredit:    "mock.php?act=getcredit",
		ReportSession:     "mock.php?act=reportsession",
		Register:          "mock.php?act=register",
	}
	return &server
}

//发起HTTP请求
func (s *Server) Query(action string, data map[string]interface{}) ([]byte, error) {
	client := &http.Client{}
	dataJson, _ := json.Marshal(data)
	reqNew := bytes.NewBuffer(dataJson)
	uri := s.Url + "/" + action
	request, _ := http.NewRequest("POST", uri, reqNew)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Token", s.Token.Token)
	request.Header.Set("AppId", s.AppId)
	response, _ := client.Do(request)
	s.Token = UserToken{}
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return body, nil
	} else {
		return []byte{}, errors.New("Http Request Failed")
	}
}

//创建请求数据Map
func (s *Server) CreateQueryData() map[string]interface{} {
	return make(map[string]interface{})
}
