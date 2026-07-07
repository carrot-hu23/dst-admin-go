#!/bin/bash

# 修正最大文件描述符数，部分docker版本给的默认值过高，会导致screen运行卡顿
ulimit -Sn 10000

# 启用 amd64 架构
dpkg --add-architecture amd64

# 添加 amd64 源
echo "deb [arch=amd64] http://archive.ubuntu.com/ubuntu jammy main universe multiverse restricted
deb [arch=amd64] http://archive.ubuntu.com/ubuntu jammy-updates main universe multiverse restricted
deb [arch=amd64] http://archive.ubuntu.com/ubuntu jammy-security main universe multiverse restricted" > /etc/apt/sources.list.d/amd64.list

# 更新源
apt update

# 安装 x86_64 运行库
apt install -y libc6:amd64 libstdc++6:amd64

cd /opt/DepotDownloader
./DepotDownloader -app 343050 -os linux -osarch 64 -dir /app/dst-dedicated-server -validate

chmod +x /app/dst-dedicated-server/bin64/dontstarve_dedicated_server_nullrenderer_x64

cd /app
exec ./dst-admin-go
