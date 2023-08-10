# dst-admin-go
> Don't Starve Together Server management panel
>
> Date: 2023/05/11

Preview Demo: http://1.12.223.51:8080/ （admin 123456）

This is a management panel for Don't Starve Together, developed in Go. It offers simple deployment, low memory usage, an aesthetically pleasing interface, and user-friendly operations. The panel provides a visual interface for easily configuring game rooms and managing online mods. It also supports the management of multiple rooms. All of these features are designed to provide a smoother and more streamlined user experience.

这是一个使用 Go 编写的饥荒管理面板，它具有简单的部署流程、低内存占用、美观的界面和简洁的操作方式。该管理面板提供了直观的可视化界面，方便用户进行房间配置和在线模组配置，同时支持多房间的便捷管理。这一切旨在提供更加流畅的使用体验。

## Deployment/部署

**目前多房间版本还有些bug没有修复(等单房间功能稳定后在迁移过来)， 萌新勿用**

注意目录必须要有读写权限。


### 二进制部署

[部署教程](https://blog.csdn.net/Dig_hoof/article/details/131296762)

[视频教程](https://www.bilibili.com/read/cv25125509)

### docker

```
docker pull hujinbo23/dst-admin-go:1.1.8
docker run -d -p8082:8082 hujinbo23/dst-admin-go:1.1.8
```

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

## 请作者喝一杯咖啡
<img src="docs/image/alipay.jpg" alt="WeChat Pay" width="200" />
<img src="docs/image/wechatpay.png" alt="WeChat Pay" width="200" />

## QQ群 反馈交流
![首页效果](docs/image/饥荒开服面板交流issue群聊二维码.png)
