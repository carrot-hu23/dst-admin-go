# 使用官方的Debian基础镜像
FROM debian:bookworm-slim

LABEL maintainer="hujinbo23 jinbohu23@outlook.com"
LABEL description="DoNotStarveTogehter server panel written in golang.  github: https://github.com/hujinbo23/dst-admin-go"

# 更新并安装必要的软件包
# 注意：Debian bookworm 将 lib32gcc1/libgcc1 重命名为 lib32gcc-s1/libgcc-s1
RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    curl \
    libcurl4-gnutls-dev:i386 \
    lib32gcc-s1 \
    lib32stdc++6 \
    libcurl4-gnutls-dev \
    libgcc-s1 \
    libstdc++6 \
    wget \
    ca-certificates \
    screen \
    procps \
    sudo \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 拷贝程序二进制文件
COPY dst-admin-go /app/dst-admin-go
RUN chmod 755 /app/dst-admin-go

COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod 755 /app/docker-entrypoint.sh

COPY config.yml /app/config.yml

# Optionally bake steamcmd and the DST dedicated server into the image, so a
# plain `docker run` with no volumes at all works immediately, with no
# first-run download. This is opt-in (default off) because app 343050 alone
# is ~4.5GB -- baking it in by default would balloon the published
# ghcr.io/<owner>/dst-admin-go image ~20x for every user, most of whom
# already mount /app/steamcmd and /app/dst-dedicated-server as volumes.
# Build with `--build-arg BAKE_GAME_SERVER=true` for a fully self-contained
# image instead. Either way, docker-entrypoint.sh installs these at
# container start if they're missing, so mounting them as volumes remains
# entirely optional in both cases -- baking just decides whether "missing"
# is already true on first boot or only after they're volume-mounted away.
ARG BAKE_GAME_SERVER=false
RUN if [ "$BAKE_GAME_SERVER" = "true" ]; then \
      mkdir -p /app/steamcmd && \
      cd /app/steamcmd && \
      wget -q http://media.steampowered.com/installer/steamcmd_linux.tar.gz && \
      tar -zxf steamcmd_linux.tar.gz && \
      rm steamcmd_linux.tar.gz && \
      mkdir -p /app/dst-dedicated-server && \
      ( bash /app/steamcmd/steamcmd.sh +force_install_dir /app/dst-dedicated-server +login anonymous +app_update 343050 validate +quit \
        || echo "Warning: could not bake the DST dedicated server into the image; docker-entrypoint.sh will install it at container start instead" ); \
    fi

# docker_dst_config seeds the panel's Settings defaults (steamcmd,
# force_install_dir, cluster, backup, mod_download_path, persistent_storage_root)
# on first run. It's baked at this stable path -- NOT copied directly into
# /app/data -- because /app/data is meant to be volume-mounted for
# persistence, and an empty mounted volume would otherwise shadow it
# entirely. docker-entrypoint.sh copies it into /app/data/dst_config only if
# that file doesn't already exist there.
COPY docker_dst_config /app/docker_dst_config.default
COPY dist /app/dist
COPY static /app/static

# 内嵌源配置信息
# 控制面板访问的端口
EXPOSE 8082/tcp
# 饥荒世界通信的端口
EXPOSE 10888/udp
# 饥荒洞穴世界的端口
EXPOSE 10998/udp
# 饥荒森林世界的端口
EXPOSE 10999/udp

# Marks this image as a container build so the panel can lock the Settings
# fields that correspond to fixed image paths (steamcmd, force_install_dir,
# cluster) instead of letting them silently drift out of sync with what
# actually exists on disk -- see internal/service/dstConfig/dst_config.go.
ENV DST_ADMIN_CONTAINER_MODE=true

# 运行命令
ENTRYPOINT ["./docker-entrypoint.sh"]