# Deployment/部署
注意目录必须要有读写权限

### 脚本一键部署
[点击下载文件](docs/install/一键部署脚本v4.6版本.sh.x)

### 二进制部署
请下载最新的release版本

[部署教程](https://blog.csdn.net/Dig_hoof/article/details/131296762)

[视频教程](https://www.bilibili.com/read/cv25125509)

### docker 部署

**第一次启动时会自动下载steamcmd和饥荒服务器，请耐心等待10-20分钟，你也可以使用挂载路径避免下载**

```
docker pull hujinbo23/dst-admin-go:1.2.0
docker run -d -p8082:8082 hujinbo23/dst-admin-go:1.2.0

**路径挂载参考**
docker 存档目录挂载 命令参考

docker run -d --name dst-admin-go \
  -p8080:8082 \
    -v /root/dstsave:/root/.klei/DoNotStarveTogether \
    -v /root/dstbackup:/app/backup \
    -v /root/dstmod:/app/mod \
    hujinbo23/dst-admin-go:1.2.0

容器存档启动路径: /root/.klei/DoNotStarveTogether
容器存档备份路径: /app/backup
容器存档模组路径: /app/mod
容器玩家日志路径: /app/dst-db
容器服务日志路径: /app/dst-admin-go.log
容器steamcmd路径: /app/steamcmd
容器饥荒服务器路径: /app/dst-dedicated-server
```

**多房间版本 萌新勿用**

目前多房间版本，默认以**32位**启动，不提供steamcmd和饥荒服务器安装，需要自己手动安装
同时相比较普通版本缺少宕机自动恢复，自动更新，定时任务等（有时间在迁移过来）

+ steamcmd 请填写 /app/steamcmd
+ 饥荒路径填写 /app/dst-dedicated-server

```
docker pull hujinbo23/dst-admin-go:home.1.0.2
docker run -d -p8082:8083 hujinbo23/dst-admin-go:home.1.0.2
```
