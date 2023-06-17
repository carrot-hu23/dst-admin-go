# dst-admin-go
> 饥荒联机版管理后台Go版本
>
> Date: 2023/05/11

使用go编写的饥荒管理面板,部署简单,占用内存少,界面美观,操作简单,提供可视化界面操作房间配置和模组在线配置,支持多房间管理

**有什么问题或者建议、bug 可以在issue里面提出来，有时间会做**

## 部署
注意目录必须要有读写权限。

## 说明
当前有两个版本，一个是可以开启多个房间，另一个是单房间的

+ 2.0.0.beta 为多房间版本（适合个人使用）
+ 1.1.6      为单房间版本（适合docker部署）

**docker 部署**
>此版本是单房间版本，
```text
docker pull 
registry.cn-hangzhou.aliyuncs.com/hujinbo23/dst-admin-go:1.1.6

docker run -d -p8082:8082 registry.cn-hangzhou.aliyuncs.com/hujinbo23/dst-admin-go:1.1.6

```

点击查看 [部署文档](./doc/install.md)

## 预览

在线预览地址 http://1.12.223.51:8082/
（admin 123456）
![首页效果](./doc/image/登录.png)
![首页效果](./doc/image/房间.png)
![首页效果](./doc/image/mod.png)
![首页效果](./doc/image/mod配置.png)
![统计效果](./doc/image/统计.png)
![面板效果](./doc/image/面板.png)
![日志效果](./doc/image/日志.png)
    

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
