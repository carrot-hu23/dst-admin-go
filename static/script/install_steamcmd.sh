#!/bin/bash

# 获取传入的参数
steam_cmd_path=$1
steam_dst_server=$2

# 判断 steam_cmd_path 是否存在，不存在则创建
if [ ! -d "$steam_cmd_path" ]; then
  mkdir -p "$steam_cmd_path"
fi

# 进入 steam_cmd_path 目录
cd "$steam_cmd_path"

# 如果 $steam_dst_server 目录不存在，则下载并解压 SteamCMD 并安装游戏服务器
retry=1
while [ ! -d "${steam_cmd_path}/steamcmd" ] || [ ! -e "${steam_cmd_path}/steamcmd/steamcmd.sh" ]; do
  if [ $retry -gt 3 ]; then
    echo "Download steamcmd failed after three times"
    exit -2
  fi
  wget http://media.steampowered.com/installer/steamcmd_linux.tar.gz -P $steam_cmd_path/steamcmd
  tar -zxvf $steam_cmd_path/steamcmd/steamcmd_linux.tar.gz -C $steam_cmd_path/steamcmd
  sleep 3
  ((retry++))
done

# 如果 $steam_dst_server 目录不存在，则下载并解压 SteamCMD 并安装游戏服务器
retry=1
while [ ! -e "${steam_dst_server}/dst-dedicated-server/bin/dontstarve_dedicated_server_nullrenderer" ]; do
  if [ $retry -gt 3 ]; then
    echo "Download Dont Starve Together Sever failed after three times"
    exit -2
  fi
  echo "Not found Dont Starve Together Sever, start to installing, try: ${retry}"
  bash $steam_cmd_path/steamcmd/steamcmd.sh +force_install_dir $steam_dst_server/dst-dedicated-server +login anonymous +app_update 343050 validate +quit
  cp $steam_cmd_path/steamcmd/linux32/libstdc++.so.6 $steam_dst_server/dst-dedicated-server/bin/lib32/
  mkdir -p $HOME/.klei/DoNotStarveTogether/MyDediServer
  sleep 3
  ((retry++))
done
bash $steam_cmd_path/steamcmd/steamcmd.sh +force_install_dir $steam_dst_server/dst-dedicated-server +login anonymous +app_update 343050 validate +quit

# 运行其他命令，这里只是做示例
echo "SteamCMD installed at $steam_cmd_path"
echo "SteamDST server installed at $steam_dst_server"