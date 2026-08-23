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

镜像已经在构建时内置了 steamcmd 和饥荒服务器本体，所以真正**必须**挂载的只有 `/app/data`
（面板设置、存档、管理员账号）。下面这条命令就是一个可以直接跑起来的最小示例：

```bash
mkdir -p ~/dstsave/data

docker run -d \
  --name dst-admin \
  -p 8082:8082 \
  -p 10888:10888/udp \
  -p 10998:10998/udp \
  -p 10999:10999/udp \
  -v ~/dstsave/data:/app/data \
  hujinbo23/dst-admin-go:latest
```

如果还想在容器重建/镜像升级时跳过重新安装 steamcmd 和游戏本体，或者想复用已有的
`/app/backup`、`/app/mod` 缓存，可以再加上这些**可选**挂载：

```bash
mkdir -p ~/dstsave/{back,mod,steamcmd,dst-dedicated-server}

docker run -d \
  --name dst-admin \
  -p 8082:8082 \
  -p 10888:10888/udp \
  -p 10998:10998/udp \
  -p 10999:10999/udp \
  -v ~/dstsave/data:/app/data \
  -v ~/dstsave/back:/app/backup \
  -v ~/dstsave/mod:/app/mod \
  -v ~/dstsave/steamcmd:/app/steamcmd \
  -v ~/dstsave/dst-dedicated-server:/app/dst-dedicated-server \
  hujinbo23/dst-admin-go:latest
```

**迁移已有存档**：默认的存档路径固定为
`/app/data/save/DoNotStarveTogether/Cluster_1`（与饥荒本体 `-cluster` 参数的默认值
`Cluster_1` 保持一致）。如果你已经有一份用其它工具/其它默认名称建立的存档，不需要在
面板里改 `cluster` 名称，直接把你的存档目录挂载到这个固定路径上即可，例如：

```bash
-v ~/my-old-backup:/app/data/save/DoNotStarveTogether/Cluster_1
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
| `/app/data` | 面板设置（`dst_config`）、游戏存档、管理员账号（`password.txt`）、SQLite 数据库、首次登录标记 —— 唯一真正需要挂载的目录 | ✅ 必须 |
| `/app/backup` | 存档备份目录 | 可选（不挂载也能正常使用，只是容器重建后备份历史会丢失） |
| `/app/mod` | MOD 缓存目录 | 可选（不挂载的话，容器重建后 mod 会重新下载） |
| `/app/steamcmd` | SteamCMD 安装目录 | 可选；不挂载的话容器每次都是全新安装（几十 MB，很快） |
| `/app/dst-dedicated-server` | 饥荒服务器文件 | 可选；不挂载的话容器每次都是全新安装（app 343050 完整安装约 4.5GB，看网络情况可能需要几分钟） |
| `/app/dst-admin-go.log` | 应用日志文件 | 可选 |
| `/app/config.yml` | 配置文件 | 可选 |

**特别说明**：
- 面板设置、存档、管理员账号、数据库、首次登录标记现在统一放在 `/app/data` 下，只挂载这一个目录即可持久化全部重要数据。
- `/app/data/dst_config`：镜像里已经内置了默认值（steamcmd、force_install_dir、cluster、backup、mod_download_path、persistent_storage_root）；如果 `/app/data` 是一个全新的空目录，容器启动时会自动从镜像内置的默认值播种这个文件，不需要手动创建。
- 官方发布的镜像**不会**内置 steamcmd/游戏本体本身（app 343050 完整安装约 4.5GB，内置会让镜像体积膨胀约 20 倍）；如果你想要一个开箱即用、完全不需要联网下载的镜像，可以用 `docker build --build-arg BAKE_GAME_SERVER=true ...` 自行构建。
- `/app/data/password.txt`：初始管理员账号信息（`admin` / `123456`），同样会在 `/app/data` 为空时自动创建；建议登录后立刻在面板里修改密码。
- 旧版本中单独挂载 `/root/.klei/DoNotStarveTogether`、`/app/password.txt`、`/app/first` 的方式仍然可以工作（未设置 `persistent_storage_root` 时会退回旧的默认路径），但不再推荐——迁移到 `/app/data` 之后挂载点更少、也更不容易出现"改了面板设置但重启后又变回默认值"的问题。

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
  -v ~/dstsave/data:/app/data \
  -v ~/dstsave/config.yml:/app/config.yml \
  hujinbo23/dst-admin-go:latest
```

## Docker Compose 示例

### 1. 创建数据目录

```bash
mkdir -p ~/dstsave/data
```

`~/dstsave/data` 是唯一必须提前创建的目录——首次启动时，容器会自动在里面播种
`dst_config`（面板设置默认值）和 `password.txt`（初始管理员账号 `admin` / `123456`），
不需要手动创建这些文件。

如果想让 `/app/backup`、`/app/mod`、`/app/steamcmd`、`/app/dst-dedicated-server`
也持久化（详见上面的[数据卷](#数据卷)说明，这几个都是可选的），可以一并创建：

```bash
mkdir -p ~/dstsave/{back,mod,steamcmd,dst-dedicated-server}
```

**目录结构**：
```
~/dstsave/
├── data/                         # 面板设置、存档、账号、数据库 —— 必须挂载
├── back/                         # 备份目录（可选）
├── mod/                          # MOD 缓存目录（可选）
├── steamcmd/                     # SteamCMD 安装目录（可选，避免每次重建容器都重新下载）
└── dst-dedicated-server/         # 饥荒服务器文件（可选，避免每次重建容器都重新下载）
```

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
      # 面板设置、存档、账号、数据库 —— 唯一必须挂载的目录
      - ${PWD}/dstsave/data:/app/data
      # 以下都是可选挂载
      - ${PWD}/dstsave/back:/app/backup
      - ${PWD}/dstsave/mod:/app/mod
      - ${PWD}/dstsave/steamcmd:/app/steamcmd
      - ${PWD}/dstsave/dst-dedicated-server:/app/dst-dedicated-server
    environment:
      - TZ=Asia/Shanghai
```

**迁移已有存档**：如果你已经有一份用其它工具或其它默认名称建立的存档，把它挂载到
`/app/data/save/DoNotStarveTogether/Cluster_1`（默认的 `cluster` 名称固定为
`Cluster_1`，与饥荒本体自身的默认值一致），就不需要在面板里改任何设置：

```yaml
      - ${PWD}/my-old-backup:/app/data/save/DoNotStarveTogether/Cluster_1
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

官方发布的镜像不内置游戏本体（约 4.5GB，内置会让镜像体积膨胀约 20 倍），所以首次启动
需要下载 SteamCMD 和饥荒服务器文件，国内网络可能较慢。可以考虑：
1. 挂载 `/app/steamcmd` 和 `/app/dst-dedicated-server` 为持久化目录，这样只有第一次
   启动才需要下载，之后重建容器都会复用已下载的内容
2. 自行用 `docker build --build-arg BAKE_GAME_SERVER=true ...` 构建一个把游戏本体
   直接内置到镜像里的版本，首次启动完全不需要联网下载（代价是镜像体积会大很多）
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
