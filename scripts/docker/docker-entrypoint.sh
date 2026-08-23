#!/bin/bash

# 修正最大文件描述符数，部分docker版本给的默认值过高，会导致screen运行卡顿
ulimit -Sn 10000

# 获取传入的参数
steam_cmd_path='/app/steamcmd'
steam_dst_server='/app/dst-dedicated-server'
data_dir='/app/data'

# 确保持久化数据目录存在，并在其为空（例如首次挂载一个全新的空 volume）时
# 从镜像内置的默认值播种 dst_config 和管理员账户文件，这样这一步与
# /app/data 是否被挂载为 volume 无关，每次启动都会执行。
mkdir -p "$data_dir"
if [ ! -f "$data_dir/dst_config" ]; then
  cp /app/docker_dst_config.default "$data_dir/dst_config"
fi
if [ ! -f "$data_dir/password.txt" ]; then
  echo "username=admin" >> "$data_dir/password.txt"
  echo "password=123456" >> "$data_dir/password.txt"
  echo "displayName=admin" >> "$data_dir/password.txt"
  echo "photoURL=xxx" >> "$data_dir/password.txt"
fi
mkdir -p /app/backup
mkdir -p /app/mod
# 与 docker_dst_config 中的 persistent_storage_root=/app/data/save 和默认的
# cluster=Cluster_1 保持一致，方便直接把已有存档挂载到这个固定路径上，
# 而不需要在面板里修改 cluster 名称。
mkdir -p "$data_dir/save/DoNotStarveTogether/Cluster_1"

# 判断 steam_cmd_path 是否存在，不存在则创建
if [ ! -d "$steam_cmd_path" ]; then
  mkdir -p "$steam_cmd_path"
fi

# 进入 steam_cmd_path 目录
cd "$steam_cmd_path"

# 镜像构建时已经内置了 steamcmd 和游戏本体，正常情况下下面两个循环会
# 立即通过检查、不会真正下载。只有在构建时安装失败，或者把这两个目录
# 挂载成了空 volume 时，才会在这里补装。
retry=1
while [ ! -d "${steam_cmd_path}" ] || [ ! -e "${steam_cmd_path}/steamcmd.sh" ]; do
  if [ $retry -gt 3 ]; then
    echo "Download steamcmd failed after three times"
    exit -2
  fi
  echo "Not found steamcmd, start to installing steamcmd, try: ${retry}"
  wget http://media.steampowered.com/installer/steamcmd_linux.tar.gz -P $steam_cmd_path
  tar -zxvf $steam_cmd_path/steamcmd_linux.tar.gz -C $steam_cmd_path
  sleep 3
  ((retry++))
done

retry=1
while [ ! -e "${steam_dst_server}/bin/dontstarve_dedicated_server_nullrenderer" ]; do
  if [ $retry -gt 3 ]; then
    echo "Download Dont Starve Together Sever failed after three times"
    exit -2
  fi
  echo "Not found Dont Starve Together Sever, start to installing, try: ${retry}"
  bash $steam_cmd_path/steamcmd.sh +force_install_dir $steam_dst_server +login anonymous +app_update 343050 validate +quit
  sleep 3
  ((retry++))
done


# 运行其他命令，这里只是做示例
echo "SteamCMD installed at $steam_cmd_path"
echo "SteamDST server installed at $steam_dst_server"


cd /app
exec ./dst-admin-go
