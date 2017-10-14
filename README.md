## 目录结构说明

1. api 用于与jms通讯的相关接口
2. client 连接后段实际服务器相关
3. sshd ssh服务
4. websocket ws服务
5. util 工具类


## API调试说明

mock.php 为mockserver文件，直接运行

> php -S 0.0.0.0:8888 

即可启动，启动后，使用如下参数启动Coco

> go run main.go -jmsurl=http://127.0.0.1:8888 -appid=sdlfjskllu4234324

## Tip

### glide

```bash
rm -rf ~/.glide
mkdir -p ~/.glide
glide mirror set https://golang.org/x/mobile https://github.com/golang/mobile --vcs git
glide mirror set https://golang.org/x/crypto https://github.com/golang/crypto --vcs git
glide mirror set https://golang.org/x/net https://github.com/golang/net --vcs git
glide mirror set https://golang.org/x/tools https://github.com/golang/tools --vcs git
glide mirror set https://golang.org/x/text https://github.com/golang/text --vcs git
glide mirror set https://golang.org/x/image https://github.com/golang/image --vcs git
glide mirror set https://golang.org/x/sys https://github.com/golang/sys --vcs git
```


## 致谢：

github.com/joushou/sshmux