package api

//用户登陆TOKEN
type UserToken struct {
	Uid     int
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
