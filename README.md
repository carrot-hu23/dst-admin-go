# dst-admin-go
> 饥荒联机版管理后台Go版本
>
> Date: 2023/05/11

在线预览地址 http://1.12.223.51:8080/ （admin 123456）

使用go编写的饥荒管理面板,部署简单,占用内存少,界面美观,操作简单,提供可视化界面操作房间配置和模组在线配置,支持多房间管理

## 部署

**目前多房间版本还有些bug没有修复(等单房间功能稳定后在迁移过来)， 萌新勿用**

注意目录必须要有读写权限。
### 二进制部署
+ 1.1.6(单房间部署)
  
    [萌新部署教程](https://blog.csdn.net/Dig_hoof/article/details/131296762)

+ 2.0.0.beta

  点击查看 [多房间部署文档](docs/install.md)

### 一键部署
目前只支持 `centos`。感谢 [SubTel](https://github.com/SubTel)提供的脚本
[dst-admin-go一键部署脚本centos版本.sh](docs/dst-admin-go一键部署脚本centos版本.sh)
支持 systemctl 命令
```sh
wget https://github.com/hujinbo23/dst-admin-go/releases/download/1.1.6.hotfix/dst-admin-go.centos.sh
```
### docker 部署
此版本是单房间版本，安装后记得在页面`系统设置`把64位改成32启动，64位启动暂时还有问题
```
docker pull hujinbo23/dst-admin-go
docker run -d -p8082:8082 hujinbo23/dst-admin-go
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