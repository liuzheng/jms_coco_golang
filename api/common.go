package api

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"coco/util"
	"fmt"
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
func (s *Server) Query(action string, data map[string]interface{}, ret interface{}) (error) {
	client := &http.Client{}
	dataJson, _ := json.Marshal(data)
	reqNew := bytes.NewBuffer(dataJson)
	uri := s.Url + "/" + action
	request, _ := http.NewRequest("POST", uri, reqNew)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Token", s.Token.Token)
	request.Header.Set("AppId", s.AppId)
	var retErr error
	//发起HTTP请求
	if response, err := client.Do(request); err == nil {
		retBody, _ := ioutil.ReadAll(response.Body)
		if response.StatusCode != 200 {
			ret = RespErrorJson{}
		}
		if err := json.Unmarshal(retBody, ret); err != nil {
			retErr = &RespError{
				Code: -500,
				Msg:  "Json返回解析发生错误",
				Raw:  err.Error(),
			}
		}
		if response.StatusCode != 200 {
			retErr = &RespError{
				Code: response.StatusCode,
				Msg:  ret.(RespErrorJson).Error,
				Raw:  response.Status,
			}
		}
	} else {
		retErr = &RespError{
			Code: -500,
			Msg:  "发起HTTP请求发生错误",
			Raw:  err.Error(),
		}
	}
	return retErr
}

//创建请求数据Map
func (s *Server) CreateQueryData() map[string]interface{} {
	return make(map[string]interface{})
}

//API的错误处理
func (re *RespError) Error() string {
	return fmt.Sprintf("API请求错误，代码：%v ，错误信息：%v", re.Code, re.Msg)
}
