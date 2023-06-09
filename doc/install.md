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

**记住 steamcmd 路径** 等下要用到


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
    database: dst-database
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

## 4. 在浏览器 输入 http://xxx:8082 进入页面
初始化用户信息,进入页面

点击右上角新建集群按钮,按照要求输入相应路径,
等待5~20分钟会自动创建饥荒服务和世界配置

### 创建集群时请不要使用 纯数字、中文、或者特殊字符，集群名称就是你存档的名称
错误集群名称示例
+ 1
+ ——
+ @@@
+ 1213

正确集群名称示例
+ caicai1
+ caicai2
+ caicai3