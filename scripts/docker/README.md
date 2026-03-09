# Docker 部署脚本

用于构建 DST Admin Go 的标准 Docker 镜像（Linux x86_64 架构）。

## 目录内容

- `Dockerfile` - Docker 镜像构建文件（基于 Ubuntu 20.04）
- `docker-entrypoint.sh` - 容器启动入口脚本
- `docker_build.sh` - 构建并推送镜像到 Docker Hub 的自动化脚本
- `docker_dst_config` - Docker 环境默认配置文件

## 快速开始

### 1. 构建镜像

```bash
# 首先在项目根目录构建 Linux 二进制文件
bash scripts/build_linux.sh

# 进入 docker 目录
cd scripts/docker

# 构建并推送镜像（需要先登录 Docker Hub）
bash docker_build.sh <version_tag>

# 示例
bash docker_build.sh 1.6.1
```

### 2. 运行容器

```bash
# 创建数据目录
mkdir -p ~/dstsave/{back,steamcmd,dst-dedicated-server}

# 运行容器
docker run -d \
  --name dst-admin \
  -p 8082:8082 \
  -p 10888:10888/udp \
  -p 10998:10998/udp \
  -p 10999:10999/udp \
  -v ~/dstsave:/root/.klei/DoNotStarveTogether \
  -v ~/dstsave/back:/app/backup \
  -v ~/dstsave/steamcmd:/app/steamcmd \
  -v ~/dstsave/dst-dedicated-server:/app/dst-dedicated-server \
  hujinbo23/dst-admin-go:latest
```

### 3. 访问管理面板

打开浏览器访问: http://localhost:8082

## 端口说明

| 端口 | 协议 | 用途 |
|-----|------|------|
| 8082 | TCP | 管理面板 Web 访问端口 |
| 10888 | UDP | 饥荒主世界（Master）通信端口 |
| 10998 | UDP | 饥荒洞穴世界（Caves）端口 |
| 10999 | UDP | 饥荒森林世界（Forest）端口 |

## 数据卷

容器内重要路径说明：

| 容器内路径 | 用途 | 是否推荐挂载 |
|-----------|------|-------------|
| `/root/.klei/DoNotStarveTogether` | 游戏存档目录 | ✅ 推荐 |
| `/app/backup` | 存档备份目录 | ✅ 推荐 |
| `/app/mod` | MOD 缓存目录 | 可选 |
| `/app/steamcmd` | SteamCMD 安装目录 | ✅ 推荐 |
| `/app/dst-dedicated-server` | 饥荒服务器文件 | ✅ 推荐 |
| `/app/dst-db` | SQLite 数据库文件 | ✅ 推荐 |
| `/app/password.txt` | 初始密码文件 | ✅ 推荐 |
| `/app/first` | 首次登录标记文件 | ✅ 推荐 |
| `/app/dst-admin-go.log` | 应用日志文件 | 可选 |
| `/app/config.yml` | 配置文件 | 可选 |

**特别说明**：
- `first` 文件：如果存在，启动时会跳过初始化界面，使用 `password.txt` 中的账号登录
- `dst-db` 文件：SQLite 数据库，包含所有配置和运行数据
- `password.txt` 文件：初始管理员账号信息，格式见 Docker Compose 示例

## 镜像特性

- **基础镜像**: Ubuntu 20.04
- **目标架构**: Linux x86_64 (amd64)
- **已安装组件**:
  - curl, wget - 网络工具
  - screen - 游戏进程管理
  - lib32gcc1, lib32stdc++6 - 32位运行库（饥荒服务器依赖）
  - libcurl4-gnutls-dev - cURL 开发库
  - procps, sudo, unzip - 系统工具

## 配置自定义

### 方法一：环境变量

```bash
docker run -d \
  -e BIND_ADDRESS="" \
  -e PORT=8082 \
  -e DATABASE=dst-db \
  hujinbo23/dst-admin-go:latest
```

### 方法二：挂载配置文件

```bash
docker run -d \
  -p 8082:8082 \
  -p 10888:10888/udp \
  -p 10998:10998/udp \
  -p 10999:10999/udp \
  -v ~/dstsave:/root/.klei/DoNotStarveTogether \
  -v ~/dstsave/back:/app/backup \
  -v ~/dstsave/steamcmd:/app/steamcmd \
  -v ~/dstsave/dst-dedicated-server:/app/dst-dedicated-server \
  -v ~/dstsave/config.yml:/app/config.yml \
  hujinbo23/dst-admin-go:latest
```

## Docker Compose 示例

### 1. 创建前置文件和目录

在使用 Docker Compose 之前，需要先创建必要的文件和目录：

```bash
# 创建所有必要的目录
mkdir -p ~/dstsave/back
mkdir -p ~/dstsave/steamcmd
mkdir -p ~/dstsave/dst-dedicated-server

# 创建 first 文件（标记非首次登录，避免进入初始化界面）
touch ~/dstsave/first

# 创建数据库文件
touch ~/dstsave/dst-db

# 创建初始密码文件
cat > ~/dstsave/password.txt << EOF
username = admin
password = 123456
displayName = admin
photoURL =
email = xxx
EOF
```

**目录结构**：
```
~/dstsave/
├── back/                        # 备份目录
├── steamcmd/                    # SteamCMD 安装目录
├── dst-dedicated-server/        # 饥荒服务器文件
├── dst-db                       # SQLite 数据库文件
├── password.txt                 # 初始密码文件
└── first                        # 首次登录标记文件
```

**说明**：
- `first` 文件：如果存在则跳过初始化界面，直接使用 `password.txt` 中的账号登录
- `dst-db` 文件：SQLite 数据库文件
- `password.txt` 文件：初始管理员账号信息

### 2. 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  dst-admin:
    image: hujinbo23/dst-admin-go:latest
    container_name: dst-admin
    restart: unless-stopped
    ports:
      - "8082:8082"
      - "10888:10888/udp"
      - "10998:10998/udp"
      - "10999:10999/udp"
    volumes:
      # 时区同步
      - /etc/localtime:/etc/localtime:ro
      - /etc/timezone:/etc/timezone:ro
      # 游戏存档目录
      - ${PWD}/dstsave:/root/.klei/DoNotStarveTogether
      # 备份目录
      - ${PWD}/dstsave/back:/app/backup
      # SteamCMD 目录
      - ${PWD}/dstsave/steamcmd:/app/steamcmd
      # 饥荒服务器目录
      - ${PWD}/dstsave/dst-dedicated-server:/app/dst-dedicated-server
      # 数据库文件
      - ${PWD}/dstsave/dst-db:/app/dst-db
      # 初始密码文件
      - ${PWD}/dstsave/password.txt:/app/password.txt
      # 首次登录标记文件
      - ${PWD}/dstsave/first:/app/first
    environment:
      - TZ=Asia/Shanghai
```

### 3. 启动容器

```bash
docker-compose up -d
```

### 4. 查看日志

```bash
# 查看容器日志
docker-compose logs -f

# 查看应用日志
docker exec -it dst-admin cat /app/dst-admin-go.log
```

## 常见问题

### 容器无法启动

查看日志排查问题：
```bash
docker logs dst-admin
```

### 游戏端口无法访问

1. 确保端口映射正确且使用 UDP 协议
2. 检查宿主机防火墙设置：
```bash
# CentOS/RHEL
firewall-cmd --add-port=10888/udp --permanent
firewall-cmd --reload

# Ubuntu/Debian
ufw allow 10888/udp
```

### 数据持久化失败

确保挂载目录有正确的权限：
```bash
chmod -R 755 ~/dstsave
```

### 游戏下载缓慢

首次启动需要下载 SteamCMD 和饥荒服务器文件（约 1-2GB），国内网络可能较慢。可以考虑：
1. 预先下载 SteamCMD 到 `~/dstsave/steamcmd` 目录
2. 预先使用 SteamCMD 下载游戏文件到 `~/dstsave/dst-dedicated-server` 目录
3. 使用代理加速 Steam 下载

## 性能建议

- **最低配置**: 2 核 CPU, 2GB 内存, 10GB 磁盘
- **推荐配置**: 4 核 CPU, 4GB 内存, 20GB 磁盘
- **生产环境**: 根据玩家数量和世界复杂度适当增加资源

## 注意事项

1. 生产环境建议使用固定版本标签，避免使用 `latest`
2. 定期备份 `~/dstsave` 目录，里面包含所有重要数据
3. 游戏端口必须使用 UDP 协议，TCP 无法正常工作
4. 容器重启后游戏进程需要手动启动（通过管理面板）
5. 首次启动会自动下载 SteamCMD 和饥荒服务器文件，需要一定时间
6. 所有数据统一存放在 `~/dstsave` 目录，便于管理和备份

## 相关链接

- [Docker Hub 镜像](https://hub.docker.com/r/hujinbo23/dst-admin-go)
- [GitHub 项目主页](https://github.com/hujinbo23/dst-admin-go)
- [饥荒联机版官方 Wiki](https://dontstarve.fandom.com/wiki/Don%27t_Starve_Together)
