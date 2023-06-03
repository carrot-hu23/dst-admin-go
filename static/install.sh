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