package api

//jms服务器配置
type Server struct {
	Url    string
	AppId  string
	AppKey string
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
	PrivateKey string
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
