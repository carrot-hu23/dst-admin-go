# dst-admin-go
> 饥荒联机版管理后台

[English](README-EN.md)/[中文](README.md)

**新面板 [泰拉瑞亚面板](https://github.com/carrot-hu23/terraria-panel-app) 支持window,linux 一键启动，内置 1449 版本**

## 推广
感谢莱卡云赞助广告

[【莱卡云】热卖套餐配置低至32元/月起，镜像内置面板，一键开服，即刻畅玩，立享优惠！](https://www.lcayun.com/aff/OYXIWEQC)
![tengxunad1](docs/image/莱卡云游戏面板.png)


**现已支持 windows 和 Linux 平台**
> 低版本window server 请使用 1.2.8 之前的版本，高版本window使用最新的版本

使用go编写的饥荒管理面板,部署简单,占用内存少,界面美观,操作简单,提供可视化界面操作房间配置和模组在线配置,支持多房间管理，备份快照等功能

## 部署
注意目录必须要有读写权限。

点击查看 [部署文档](https://carrot-hu23.github.io/dst-admin-go-docs/)

## 预览

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


