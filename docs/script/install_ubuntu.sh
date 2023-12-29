#!/bin/bash

sudo dpkg --add-architecture i386
sudo apt-get update
sudo apt-get install -y lib32gcc1
sudo apt-get install -y libcurl4-gnutls-dev:i386
sudo apt-get install -y screen
sudo apt-get install -y glibc

mkdir ~/steamcmd
cd ~/steamcmd

wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz
tar -xvzf steamcmd_linux.tar.gz
./steamcmd.sh +login anonymous +force_install_dir ~/dst-dedicated-server +app_update 343050 validate +quit

cp ~/steamcmd/linux32/libstdc++.so.6 ~/dst-dedicated-server/bin/lib32/

#Abandon the use of script execution, and change to execute directly through java code
#cd ~/dst-dedicated-server/bin
#echo ./dontstarve_dedicated_server_nullrenderer -console -cluster MyDediServer -shard Docker_M > overworld.sh
#echo ./dontstarve_dedicated_server_nullrenderer -console -cluster MyDediServer -shard Docker_C > cave.sh
#
#chmod +x overworld.sh
#chmod +x cave.sh

mkdir -p ~/.klei/DoNotStarveTogether/MyDediServer

cd ~



