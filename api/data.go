package api

//jms服务器配置
type Server struct {
	Url     string
	AppId   string
	AppKey  string
	WsPort  int
	SshPort int
	Ip      string
	Token   UserToken
	Action  Action
}

type UserAuth struct {
	Username   string
	Server     UserServer
	Action     Action
	UserPubKey UserPubKey
	UserToken  UserToken
}
type UserServer interface {
	Query(action string, data map[string]interface{}) ([]byte, error)
	CreateQueryData() map[string]interface{}
}

//api操作对应URL
type Action struct {
	GetUserPubKey     string
	GetUserToken      string
	GetMachineList    string
	GetLoginCredit    string
	CheckMonitorToken string
	Register          string
	ReportSession     string
}

//用户登陆TOKEN
type UserToken struct {
	Token   string
	Expired int
}

//服务器登陆凭证
type LoginCredit struct {
	Sid        int
	Username   string
	PrivateKey string `json:"private_key"`
}

//服务器列表单项
type Machine struct {
	Sid    int
	Name   string
	Ip     string
	Port   int
	Remark string
	Users  []MachineUser
}

//服务器列表可用用户
type MachineUser struct {
	Uid      int
	Username string
}

//用户pubkey返回
type UserPubKey struct {
	Ticket string
	Key    string
}

//监控权限返回
type ResponsePass struct {
	Pass bool
}
