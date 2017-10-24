package api

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"coco/util"
	"coco/util/log"
	"coco/util/errors"
)

//初始化一个ApiServer
func New() *Server {
	server := Server{
		Url:     *util.JmsUrl,
		AppId:   *util.AppId,
		appKey:  *util.AppKey,
		Ip:      *util.Ip,
		WsPort:  *util.WsPort,
		SshPort: *util.SshPort,
	}
	server.Action = Action{
		GetUserPubKey:       "mock.php?act=getpubkey",
		GetUserToken:        "mock.php?act=getusertoken",
		CheckMonitorToken:   "mock.php?act=checkmonitortoken",
		GetMachineList:      "mock.php?act=machines",
		GetMachineGroupList: "mock.php?act=machine_groups",
		GetLoginCredit:      "mock.php?act=getcredit",
		ReportSession:       "mock.php?act=reportsession",
		ReportSessionClose:  "mock.php?act=reportsessionclose",
		Register:            "mock.php?act=register",
	}
	return &server
}

//发起HTTP请求
func (s *Server) Query(action string, data map[string]interface{}, ret interface{}) (aErrData errors.Error) {
	client := &http.Client{}
	dataJson, sErr := json.Marshal(data)
	log.Debug("REQUEST", "%v", string(dataJson))
	reqNew := bytes.NewBuffer(dataJson)
	uri := s.Url + "/" + action
	request, _ := http.NewRequest("POST", uri, reqNew)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Token", s.Token.Token)
	request.Header.Set("AppId", s.AppId)
	//发起HTTP请求
	if response, sErr := client.Do(request); sErr == nil {
		retBody, _ := ioutil.ReadAll(response.Body)
		log.Debug("RESPONSE", "%v", string(retBody))
		sErr = json.Unmarshal(retBody, &aErrData)
		sErr = json.Unmarshal(retBody, ret)
		log.Debug("Query", "%v", response.StatusCode)
		aErrData = errors.New("", response.StatusCode)
	}
	if sErr != nil {
		log.Error("APISYS", sErr.Error())
	}
	return
}

//创建请求数据Map
func (s *Server) CreateQueryData() map[string]interface{} {
	return make(map[string]interface{})
}
