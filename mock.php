<?php
/**
 * Created by IntelliJ IDEA.
 * User: xrain
 * Date: 2017/10/11
 * Time: 上午11:30
 */
switch ($_GET['act']) {
    case 'getpubkey':
        echo json_encode([
            'key' => "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDALyKgUwHamftJVQyl9jJDJiINo+omlQNyrB0jpJn9wXjapO8VDf8rjYhSzCBovoswq3NFwcEXzIcZc38nnZbWJwPCL7wsXapP1qiE65hgEvEYuIwpIystNip2pd+H4rrDEIboWY5rA9VoGexC5GSKHwwqw+hidQPtlIBEjugQR1eQQTIuhQQX2H3zrergybwAWlX3jPcxpp4JguH/6QrgkybTFkrHdCOjt8d7LjduiH/AbCIn9C6bcp3sJdN2f+wxP/aVFBJ9G1I/pToY8dBzmeBAsuVO8aYBy7VmfnD3hGJPfQMhHU+d6fzZfVURaSigNrdcqVlo7w7rMsE5HGvB root@bab0e381d0d8",
            'ticket' => md5(time() + time() * rand(1111, 89999))
        ]);
        break;
    case 'checkmonitortoken':
        echo json_encode([
            'pass' => false
        ]);
        break;
    case 'machines':
        echo json_encode([
            [
                'sid' => 1,
                'name' => "测试服务器1",
                'ip' => "112.74.170.194",
                'port' => 9011,
                'remark' => "这是一台测试服务器1号",
                'users' => [
                    [
                        'uid' => 1,
                        'username' => 'root'
                    ],
                    [
                        'uid' => 2,
                        'username' => 'test'
                    ],
                ]
            ],
            [
                'sid' => 2,
                'name' => "测试服务器2",
                'ip' => "112.74.170.194",
                'port' => 9011,
                'remark' => "这是一台测试服务器1号",
                'users' => [
                    [
                        'uid' => 4,
                        'username' => 'root'
                    ],
                    [
                        'uid' => 8,
                        'username' => 'test'
                    ],
                ]
            ]
        ]);
        break;
    case 'getusertoken':
        echo json_encode([
            "token" => 'DWSSe2EySjzd7HmSDNW68umhSdgUmUavn9S2YBt',
            'expired' => time() + 3600
        ]);
        break;
    case 'getcredit':
        echo json_encode([
            'ip' => "112.74.170.194",
            'port' => 9011,
            'username' => 'root',
            'private_key' => <<<EOF
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
EOF
        ]);
        break;
    default:
        echo json_encode([
            'pass' => true
        ]);
        break;
}
