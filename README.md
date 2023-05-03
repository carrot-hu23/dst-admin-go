# dst-admin-go
> 饥荒联机版管理后台Go版本
>
> Date: 2022/12/18


## 预览

在线预览地址 http://1.12.223.51:8082/

![首页效果](./doc/登录.png)
![统计效果](./doc/统计.png)
![面板效果](./doc/面板.png)
![日志效果](./doc/日志.png)
## 运行

**修改config.yml**
```
#端口
port: 8082
db: dst-db
```

**修改dst_config**（也可以通过页面修改）
```
# steamcmd 位置
steamcmd=/root/steamcmd/
# steamcmd 饥荒安装的位置
force_install_dir=/root/dst/
# 要启动的服务器
cluster=cluster2
# 游戏备份的路径
backup=C:\Users\xm\Desktop\饥荒配置文件和建家截图\饥荒存储备份\backup
# 游戏mod下载的路径
mod_download_path=/download_mod
```

运行
```
./dst-admin-go
```

## 打包

window 下打包 Linux 二进制 （由于sqlite受操作系统影响，Linux二进制请在Linux环境build）
```
打开 cmd
set GOARCH=amd64
set GOOS=linux

go build
```
