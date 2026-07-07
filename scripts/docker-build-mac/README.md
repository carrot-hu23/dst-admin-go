# Docker 部署脚本（Mac ARM64）

用于在 Mac ARM64（Apple Silicon: M1/M2/M3）平台上构建和运行 DST Admin Go 的特殊 Docker 镜像。

## 背景说明

饥荒联机版（Don't Starve Together）服务器原生只支持 x86_64 架构，无法直接在 ARM64 设备上运行。本镜像通过 **Box64** 动态翻译技术实现 x86_64 程序在 ARM64 平台上执行，使得 Apple Silicon Mac 也能运行饥荒服务器。

## 技术架构

```
ARM64 主机（Apple Silicon Mac）
  └─> Docker 容器（ARM64）
       ├─> dst-admin-go（原生 ARM64 Go 程序）
       ├─> Box64（x86_64 → ARM64 动态翻译层）
       │    └─> DST Server（x86_64 游戏服务器）
       └─> DepotDownloader（ARM64 版本，用于下载游戏文件）
```

## 目录内容

- `Dockerfile` - ARM64 优化的镜像构建文件（Ubuntu 22.04）
- `docker-entrypoint.sh` - 容器启动脚本
- `docker_dst_config` - 默认配置文件
- `dst-mac-arm64-env-install.md` - 手动安装环境的详细步骤文档（非 Docker 部署）

## 核心组件

| 组件 | 版本 | 用途 |
|-----|------|------|
| Box64 | 最新版 | x86_64 到 ARM64 的动态二进制翻译器（启用 ARM_DYNAREC 优化） |
| .NET Runtime | 8.0 | DepotDownloader 运行依赖 |
| DepotDownloader | 3.4.0 | Steam 内容下载工具（ARM64 版本） |

## 快速开始

### 1. 构建镜像

```bash
# 在 Mac ARM64 机器上执行
cd scripts/docker-build-mac

# 构建镜像
docker build --platform linux/arm64 -t dst-admin-go-arm64:latest .
```

### 2. 运行容器

```bash
# 创建数据目录
mkdir -p ~/dstsave/{back,dst-dedicated-server}

# 运行容器
docker run -d \
  --name dst-admin-arm64 \
  --platform linux/arm64 \
  -p 8082:8082 \
  -p 10888:10888/udp \
  -p 10998:10998/udp \
  -p 10999:10999/udp \
  -v ~/dstsave:/root/.klei/DoNotStarveTogether \
  -v ~/dstsave/back:/app/backup \
  -v ~/dstsave/dst-dedicated-server:/app/dst-dedicated-server \
  dst-admin-go-arm64:latest
```

### 3. 访问管理面板

打开浏览器访问: http://localhost:8082

## 性能对比

### 性能表现

| 指标 | x86_64 原生 | ARM64 + Box64 |
|-----|------------|---------------|
| CPU 性能 | 100% | 60-80% |
| 内存占用 | 基准 | +30-40% |
| 启动速度 | 快 | 中等 |
| 稳定性 | 完美 | 良好（偶尔崩溃） |

### 适用场景

✅ **推荐用于**:
- 开发和测试环境
- 小型私服（<10 人）
- 个人学习和实验

❌ **不推荐用于**:
- 大型公开服务器
- 高并发生产环境
- 对性能要求严格的场景

## 镜像特性

- **基础镜像**: Ubuntu 22.04
- **目标架构**: ARM64 (aarch64)
- **时区设置**: Asia/Shanghai
- **已安装组件**:
  - Box64（从源码编译，启用 ARM_DYNAREC 优化）
  - .NET 8.0 Runtime
  - DepotDownloader（ARM64 版本）
  - screen, wget, curl, git
  - Python 3, pip
  - build-essential, cmake（用于编译 Box64）

## 环境变量

| 变量名 | 默认值 | 说明 |
|-------|--------|------|
| `DST_DIR` | `/dst-server` | 饥荒服务器安装目录 |
| `DEBIAN_FRONTEND` | `noninteractive` | 非交互式安装模式 |

## Docker Compose 示例

### 1. 创建前置文件和目录

在使用 Docker Compose 之前，需要先创建必要的文件和目录：

```bash
# 创建所有必要的目录
mkdir -p ~/dstsave/back
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
├── dst-dedicated-server/        # 饥荒服务器文件
├── dst-db                       # SQLite 数据库文件
├── password.txt                 # 初始密码文件
└── first                        # 首次登录标记文件
```

**说明**：
- `first` 文件：如果存在则跳过初始化界面，直接使用 `password.txt` 中的账号登录
- `dst-db` 文件：SQLite 数据库文件
- `password.txt` 文件：初始管理员账号信息
- ARM64 版本不使用 SteamCMD，而是通过 DepotDownloader 下载游戏文件

### 2. 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  dst-admin-arm64:
    image: dst-admin-go-arm64:latest
    container_name: dst-admin-arm64
    platform: linux/arm64
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
      # 饥荒服务器目录（ARM64 版本使用 DepotDownloader 下载）
      - ${PWD}/dstsave/dst-dedicated-server:/app/dst-dedicated-server
      # 数据库文件
      - ${PWD}/dstsave/dst-db:/app/dst-db
      # 初始密码文件
      - ${PWD}/dstsave/password.txt:/app/password.txt
      # 首次登录标记文件
      - ${PWD}/dstsave/first:/app/first
    environment:
      - TZ=Asia/Shanghai
    # 为 ARM64 翻译层提供更多资源
    deploy:
      resources:
        limits:
          memory: 4G
        reservations:
          memory: 2G
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
docker exec -it dst-admin-arm64 cat /app/dst-admin-go.log
```

## 手动安装（非 Docker）

如果需要在 ARM64 Linux 系统上手动部署，请参考 `dst-mac-arm64-env-install.md` 文档。

关键步骤概述：
1. 安装 .NET 8 Runtime
2. 从源码编译安装 Box64（启用 ARM_DYNAREC）
3. 下载 DepotDownloader（ARM64 版本）
4. 使用 DepotDownloader 下载饥荒服务器
5. 配置 Box64 环境变量

## 故障排查

### Box64 启动失败

检查 Box64 是否正确安装：
```bash
docker exec -it dst-admin-arm64 box64 --version
```

查看 Box64 日志（如果有）：
```bash
docker exec -it dst-admin-arm64 cat /var/log/box64.log
```

### 游戏服务器崩溃

ARM64 翻译层可能遇到兼容性问题，查看容器日志：
```bash
# 查看容器日志
docker logs dst-admin-arm64

# 进入容器查看 screen 会话
docker exec -it dst-admin-arm64 screen -ls
docker exec -it dst-admin-arm64 screen -r <session_name>
```

### 性能不足

优化建议：
1. 增加容器内存限制: `--memory=4g`
2. 减少游戏世界大小和复杂度
3. 减少加载的 MOD 数量
4. 降低玩家数量上限

### 游戏下载失败

首次启动需要通过 DepotDownloader 下载游戏文件：
```bash
# 手动下载游戏文件
docker exec -it dst-admin-arm64 /opt/DepotDownloader/DepotDownloader \
  -app 343050 -os linux -osarch 64 -dir /dst-server -validate
```

## 性能优化建议

### 1. 资源配置

```bash
docker run -d \
  --cpus="4" \
  --memory="4g" \
  --memory-swap="6g" \
  ... # 其他参数
```

### 2. Box64 优化

Box64 在构建时已启用以下优化：
- `ARM_DYNAREC=ON` - 启用 ARM 动态重编译（显著提升性能）
- `CMAKE_BUILD_TYPE=RelWithDebInfo` - 优化构建模式

### 3. 游戏配置优化

- 选择较小的世界大小（Small/Medium）
- 减少世界生成的生物数量
- 避免使用性能密集型 MOD
- 限制玩家数量（建议 ≤ 8 人）

## 与标准版对比

| 特性 | 标准版（x86_64） | ARM64 版 |
|-----|----------------|----------|
| 架构 | x86_64 | ARM64 + Box64 翻译 |
| 部署难度 | 简单 | 中等 |
| 性能 | 100% | 60-80% |
| 内存占用 | 标准 | +30-40% |
| 兼容性 | 完美 | 良好（偶尔问题） |
| 适用场景 | 生产环境 | 开发/测试/小型服务器 |
| 镜像大小 | ~800MB | ~1.2GB |

## 注意事项

1. **仅适用于 Mac ARM64**: M1/M2/M3/M4 系列芯片
2. **首次启动耗时长**: 需要下载约 1-2GB 的游戏文件
3. **不建议生产环境**: 性能和稳定性不如原生 x86_64
4. **内存需求高**: 建议至少 4GB 可用内存
5. **可能的兼容性问题**: 某些 MOD 可能在 Box64 下无法正常工作
6. **崩溃处理**: 游戏崩溃后需要手动重启（通过管理面板）

## 已知限制

- 某些 CPU 密集型 MOD 可能导致性能下降
- 大型世界（Huge）可能出现卡顿
- 多世界（Master + Caves）同时运行时内存压力大
- Box64 翻译层偶尔可能遇到兼容性问题导致崩溃

## 参考资源

- [Box64 项目](https://github.com/ptitSeb/box64) - x86_64 到 ARM64 翻译器
- [DepotDownloader](https://github.com/SteamRE/DepotDownloader) - Steam 内容下载工具
- [.NET Runtime 下载](https://dotnet.microsoft.com/download/dotnet/8.0) - .NET 8 运行时
- [DST 专用服务器 Wiki](https://dontstarve.fandom.com/wiki/Guides/Don%E2%80%99t_Starve_Together_Dedicated_Servers)

## 支持与反馈

如果在 ARM64 环境下遇到问题，请在 GitHub Issues 中注明：
- 设备型号（M1/M2/M3 等）
- macOS 版本
- Docker 版本
- 具体的错误日志

由于 ARM64 是实验性支持，某些问题可能无法完全解决。生产环境建议使用标准的 x86_64 版本。
