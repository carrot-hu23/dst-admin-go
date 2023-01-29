# dst-admin-go
> 饥荒联机版管理后台Go版本
>
> Date: 2022/12/18

## 运行

```
go run .\main.go port=:8888
```
port=: 可以指定端口启动

## 打包

window 下打包 Linux 二进制
```
打开 cmd
set GOARCH=amd64
set GOOS=linux

go build
```

## 全局处理异常
https://blog.csdn.net/u014155085/article/details/106733391

## byte[] to string
https://www.yisu.com/zixun/621470.html

## session
https://www.w3cschool.cn/yqbmht/ndc5uozt.html
https://studygolang.com/articles/34361


## 饥荒指令
https://www.bilibili.com/read/cv5536132/

## API

[-] backupController

[-] homeController

[-] loginController

[-] mainController

[x] playerController

[x] settingController

[-] systemController

[-] userController
