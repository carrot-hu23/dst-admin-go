# dst-mac-arm64-env-install

安装基础依赖

```shell
apt update
apt install -y wget unzip tar
```



下载 DepotDownloader

```shell
cd /opt
wget https://github.com/SteamRE/DepotDownloader/releases/latest/download/DepotDownloader-linux-arm64.zip -O DepotDownloader.zip || \
wget https://github.com/SteamRE/DepotDownloader/releases/latest/download/DepotDownloader.zip -O DepotDownloader.zip
```

安装 .NET 运行时

```shell
wget https://packages.microsoft.com/config/ubuntu/22.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
dpkg -i packages-microsoft-prod.deb
apt update
apt install -y dotnet-runtime-8.0
```

下载 饥荒服务器

```shell
unzip DepotDownloader.zip -d DepotDownloader
cd DepotDownloader
./DepotDownloader -app 343050 -os linux -osarch 64 -dir /app/dst-dedicated-server -validate
```



```
apt update
apt install -y git cmake build-essential

# 克隆源码
git clone https://github.com/ptitSeb/box64.git /opt/box64
cd /opt/box64

# 创建构建目录
mkdir build && cd build

# 配置并启用 ARM 动态编译器（性能更好）
cmake .. -DARM_DYNAREC=ON -DCMAKE_BUILD_TYPE=RelWithDebInfo

# 编译
make -j$(nproc)

# 安装
make install

cp box64 /usr/local/bin/
```



```
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

```

