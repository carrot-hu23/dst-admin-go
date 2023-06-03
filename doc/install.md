# 安装部署环境

以 Ubuntu 系统为例

## 1. 安装 steamcmd （如果服务器之前已经安装过请跳过此步）
```
#!/bin/bash

#准备安装饥荒
sudo apt-get install -y lib32gcc1
sudo apt-get install -y libcurl4-gnutls-dev:i386
sudo apt-get install -y screen

mkdir ~/steamcmd
cd ~/steamcmd
if [[ ! -f 'steamcmd_linux.tar.gz' ]]; then
    wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz
else
    echo -e "steamcmd_linux.tar.gz 已下载"
fi

tar -xvzf steamcmd_linux.tar.gz
./steamcmd.sh +login anonymous +force_install_dir ~/dontstarve_dedicated_server +app_update 343050 validate +quit

cp ~/steamcmd/linux32/libstdc++.so.6 ~/dontstarve_dedicated_server/bin/lib32/
mkdir -p ~/.klei/DoNotStarveTogether/MyCluster1
cd ~
```

记住这两个路径

+ steamcmd  安装路径为 ~/steamcmd
+ 饥荒 安装的路径为 ~/dontstarve_dedicated_server


## 2. 从 release 下载 稳定的版本，并解压
1. 从release下载 dst-admin-go.tgz

2. 解压，上传到服务器

3. 修改config.yml 配置（端口）
    | 配置              | 解释                      | 是否必须|
    | ----------------- | ------------------------- | -------|
    | port          | 端口          | 是 |
    | db | 数据库名称（可以随便叫啥）              | 是 |

    **参考配置**
    ```yml
    port: 8082
    db: dst-database
    ```
## 3. 启动

```
chmod +x dst-admin-go
nohup ./dst-admin-go >dst-admin-go/log &
```
如果想要关掉服务
```
ps -ef | grep dst-admin-go
```
找到进程号 kill -9
