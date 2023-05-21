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
    db: dst-db
    ```

4. 修改dst_config 配置文件
   
    | 配置              | 解释                      | 是否必须|
    | ----------------- | ------------------------- | -------|
    | steamcmd          | steamcmd安装路径          | 是 |
    | force_install_dir | 饥荒安装路径              | 是 |
    | cluster           | 要启动房间的名称          | 是 |
    | backup            | 存档备份位置              | 是 |
    | mod_download_path | mod下载位置(路径要求存在,这个下载路径和饥荒mod无关) | 是 |
    
    
    **========== 萌新可以忽略这一行 ==========**
    
    如果你的启动路径不是 ` /$HOME/.klei/DoNotStarveTogether `, 请修改, 否则忽略
    （目前不兼容这种方式启动）

    | 配置                          | 解释 |
    | ----------------------------- | ---- |
    | persistent_storage_root       |      |
    | donot_starve_server_directory |      |
    
    [参考命令](https://dontstarve.fandom.com/zh/wiki/%E5%A4%9A%E4%BA%BA%E7%89%88%E9%A5%A5%E8%8D%92%E7%8B%AC%E7%AB%8B%E6%9C%8D%E5%8A%A1%E5%99%A8?variant=zh#%E5%90%AF%E5%8A%A8%E5%8F%82%E6%95%B0)
    
    ```sh
    dontstarve_dedicated_server_nullrenderer -console -persistent_storage_root " + persistent_storage_root + "-conf_dir " + donot_starve_server_directory + " -cluster " + cluster + " -shard " + DST_CAVES + " ;"
    ```
## 3. 启动

```
nohup ./dst-admin-go >dst-admin-go/log &
```
如果想要关掉服务
```
ps -ef | grep dst-admin-go
```
找到进程号 kill -9
