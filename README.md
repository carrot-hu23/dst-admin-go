# dst-admin-go
> 饥荒联机版管理后台

**新面板 [泰拉瑞亚面板](https://github.com/carrot-hu23/terraria-panel-app) 支持window,linux 一键启动，内置 1449 版本**

## 推广
[【腾讯云】热卖套餐配置低至32元/月起，助您一键开服，即刻畅玩，立享优惠！](https://cloud.tencent.com/act/cps/redirect?redirect=5878&cps_key=8478a20880d339923787a350f9f8cbf5&from=console)
![tengxunad1](docs/image/tengxunad1.png)

**Now，Support Windows and Linux  platform**

**现已支持 windows 和 Linux 平台**

This is a management panel for Don't Starve Together, developed in Go. It offers simple deployment, low memory usage, an aesthetically pleasing interface, and user-friendly operations. The panel provides a visual interface for easily configuring game rooms and managing online mods. It also supports the management of multiple rooms. All of these features are designed to provide a smoother and more streamlined user experience.

使用go编写的饥荒管理面板,部署简单,占用内存少,界面美观,操作简单,提供可视化界面操作房间配置和模组在线配置,支持多房间管理，备份快照等功能

新增 **暗黑主题**，**国际化**，支持**多层世界**，支持更大屏幕显示

## 部署
注意目录必须要有读写权限。

点击查看 [部署文档](docs/install.md)

## 预览

在线预览地址 http://1.12.223.51:8082/
（admin 123456）
![首页效果](docs/image/登录.png)
![首页效果](docs/image/房间.png)
![首页效果](docs/image/mod.png)
![首页效果](docs/image/mod配置.png)
![统计效果](docs/image/统计.png)
![面板效果](docs/image/面板.png)
![日志效果](docs/image/日志.png)


## 运行

**修改config.yml**
```
#端口
port: 8082
database: dst-db
```


运行
```
go mod tidy
go run main.go
```

## 打包


### window 打包

window 下打包 Linux 二进制

```
打开 cmd
set GOARCH=amd64
set GOOS=linux

go build
```

## QQ 群
![QQ 群](docs/image/饥荒开服面板交流issue群聊二维码.png)


