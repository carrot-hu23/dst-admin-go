# dst-admin-go
> dst-admin-go manage web

[English](README-EN.md)/[中文](README.md)

Don't Starve management panel written in go, easy to deploy, small memory usage, beautiful interface, simple operation, provides visual interface to operate room configuration and module online configuration, supports multi-room management, backup snapshot and other functions

## Preview

![首页效果](docs/image/登录.png)
![首页效果](docs/image/房间.png)
![首页效果](docs/image/mod.png)
![首页效果](docs/image/mod配置.png)
![统计效果](docs/image/统计.png)
![面板效果](docs/image/面板.png)
![日志效果](docs/image/日志.png)



## Run

**config.yml**
```
#端口
port: 8082
database: dst-db
```


run
```
go mod tidy
go run main.go
```

## Build


```
打开 cmd
set GOARCH=amd64
set GOOS=linux

go build
```

## QQ Group
![QQ 群](docs/image/饥荒开服面板交流issue群聊二维码.png)