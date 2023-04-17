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
https://dontstarve.fandom.com/zh/wiki/%E6%8E%A7%E5%88%B6%E5%8F%B0/%E5%A4%9A%E4%BA%BA%E7%89%88%E9%A5%91%E8%8D%92%E4%B8%AD%E7%9A%84%E5%91%BD%E4%BB%A4?variant=zh

## API

[-] backupController

[-] homeController

[-] loginController

[-] mainController

[x] playerController

[x] settingController

[-] systemController

[-] userController
