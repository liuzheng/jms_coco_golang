package api

import (
	"encoding/json"
	"golang.org/x/crypto/ssh"
	log "github.com/liuzheng712/golog"
)

// 获取可用服务器列表
func (s *Server) GetList() ([]Machine, error) {
	data := s.CreateQueryData()
	res, _ := s.Query(s.Action.GetMachineList, data)
	var rd []Machine
	err := json.Unmarshal(res, &rd)
	return rd, err
}

// 获取服务器登陆凭证
func (s *Server) GetLoginCredit(serverId, userId int) (LoginCredit, error) {
	data := s.CreateQueryData()
	res, _ := s.Query(s.Action.GetLoginCredit, data)
	var rd LoginCredit
	err := json.Unmarshal(res, &rd)
	return rd, err
}

func (m *Machine) Signer() (signer ssh.Signer, err error) {
	signer, err = ssh.ParsePrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwC8ioFMB2pn7SVUMpfYyQyYiDaPqJpUDcqwdI6SZ/cF42qTv
FQ3/K42IUswgaL6LMKtzRcHBF8yHGXN/J52W1icDwi+8LF2qT9aohOuYYBLxGLiM
KSMrLTYqdqXfh+K6wxCG6FmOawPVaBnsQuRkih8MKsPoYnUD7ZSARI7oEEdXkEEy
LoUEF9h9863q4Mm8AFpV94z3MaaeCYLh/+kK4JMm0xZKx3Qjo7fHey43boh/wGwi
J/Qum3Kd7CXTdn/sMT/2lRQSfRtSP6U6GPHQc5ngQLLlTvGmAcu1Zn5w94RiT30D
IR1Pnen82X1VEWkooDa3XKlZaO8O6zLBORxrwQIDAQABAoIBAAK3J8gYJCiQkGpi
10DpD6C/Qy/qJv7QkDHomKjORJa/SxPbzL0Ba5+T6l8xHywbtR/E7GYJ7M6HR0lm
xl8q7NytxHHT9taxpBxorgziv05sNhfhU6mpI40J/OOcSGJFI9ppu15BFbARZ8yG
wT3YuErzUVqxVfE+MgYDWSSe2EySjzd7HmSDNW68umhSdgUmUavn9S2YBtUdcevj
hu+tA6QHoBajsytMM1lU8c/fbKe4sXHkiDtoobsE5xyXg4fKcgL/04RmBchXEPWc
yi69EE28UHYUWlz4vD3XN7K3Adb8PKw/133lwEbSqxFzpiGGjFKJHssu3MJExpDm
e1nqOgECgYEA33tb0nh0BG3UUp7tDq3tnDQzVjB6t2av7QpXvNXY4nBS+zUUfbFl
5lkWsEQ2hkO8cxf/N+TIBq3hNj6AaQzKU+qMORBzoA0qhBEl5Mq+J8I8FyXmsCl1
FeQMF02loK/YH/rqiRyxaTA8vh+QmGXeaqeJFZFbannnUSzjKp7iy9ECgYEA3CXx
kq5vj/12Tm1k+0UDJ3ne0gJz/1bGTRpp8xG/nwbYvhEScQa/i+9RJ5OD/JmPdYzJ
N6/bqxDHjfgd6uGy4oU59L9c8PRn4D1+1kPnnxKCOTzadgY4OAsZUUKsbU9keXzK
8hJ0JU1YrdRENoot5O3GJU3jPJUIncOM5HQwzPECgYEA0U2UYttjNR6zwymLNbtZ
lXkiN2+yDwCCdcvA/l9+AB1Y6mL9LPcvS3xWjoFkW30nCtgHtI51dCP3kDkbRMew
2sUJzQwGbnjGP+hbiBsF2vDEHJ3nn0dEnFr3o/+ZxpPKI7F5ZS7VTYHa8elYIBMQ
Ku17qxJ2/pLrPRCANVOXb3ECgYEAhOT64XsfqaJuKoXrMavFba0qD0if0YIGj0Dk
uHD8ZflyGbqgkU5rjwPWz7dDM9dXPLTSOyWTy76DWHZxSPsaS5f1FoP2jRZdJoa4
7IttZR99MSRFFH5Irbw1elJvWEVW9+eXc24kYuhAIh9DVlEvx0SqGpbcGBQ3cVU8
1iQ12WECgYBsP9cGMuQQKfF+cbzg4jU5avgT/uVkAlcpIiHOWfFtwlFR9yoqOynq
6ynjWVI0qZULhkkbWlx+N0vbNptvGW9jCaEsGy859p4qCjP3cCtn7WXUYF8KDAtR
lnwsTWISKeoSEuGNmJU1iALYS4eWD3FMgxTm3dTSMDphAxbXe5RRFw==
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		log.Error("newClient", "unable to parse private key: %v", err)
		//client.Success = false
		return
	}
	return signer, nil
}