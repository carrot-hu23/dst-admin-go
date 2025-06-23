# Deployment/部署

注意目录必须要有读写权限

### 脚本一键部署

请加 QQ 群获取

### 二进制部署

请下载最新的 release 版本

[部署教程](https://blog.csdn.net/Dig_hoof/article/details/131296762)

[视频教程](https://www.bilibili.com/read/cv25125509)

### docker 部署

**第一次启动时会自动下载 steamcmd 和饥荒服务器，请耐心等待 10-20 分钟，你也可以使用挂载路径避免下载**

自己映射对应的端口

```
docker pull hujinbo23/dst-admin-go:1.3.1
docker run --name dst -d \
  -p 8084:8082 \
  -p 10999:10999/udp \
  -p 10998:10998/udp \
  -p 10888:10888/udp \
  -v /root/dstsave:/root/.klei/DoNotStarveTogether \
  -v /root/dstsave/backup:/app/backup \
  -v /root/steamcmd:/app/steamcmd \
  -v /root/dst-dedicated-server:/app/dst-dedicated-server \
  hujinbo23/dst-admin-go:1.3.1
```

**路径参考**

```
+ 容器存档启动路径: /root/.klei/DoNotStarveTogether
+ 容器存档备份路径: /app/backup
+ 容器存档模组路径: /app/mod
+ 容器玩家日志路径: /app/dst-db 这是一个文件
+ 容器服务日志路径: /app/dst-admin-go.log
+ 容器启动饥荒路径: /app/dst-dedicated-server
+ 容器启steamcmd：/app/steamcmd
```

#### 1.2.5 及其之前的版本

启动后记得 去页面 `系统设置页面` 改成这样

```
steamcmd安装路径
/app/steamcmd
饥荒服务器安装路径
/app/dst-dedicated-server
```
