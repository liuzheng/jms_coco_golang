package api

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
)

//初始化一个ApiServer
func New() *Server {
	apiServer := Server{}
	apiServer.Action = Action{
		GetUserPubKey:     "test.php?act=getpubkey",
		GetUserToken:      "test.php?act=getusertoken",
		CheckMonitorToken: "test.php?act=checkmonitortoken",
		GetMachineList:    "test.php?act=machines",
		GetLoginCredit:    "test.php?act=getcredit",
		ReportSession:     "test.php?act=reportsession",
		Register:          "test.php?act=register",
	}
	return &apiServer
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
