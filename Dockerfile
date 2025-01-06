# 使用官方的Ubuntu基础镜像
FROM ubuntu:20.04

LABEL maintainer="hujinbo23 jinbohu23@outlook.com"
LABEL description="DoNotStarveTogehter server panel written in golang.  github: https://github.com/hujinbo23/dst-admin-go"

# 更新并安装必要的软件包
RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y \
    curl \
    libcurl4-gnutls-dev:i386 \
    lib32gcc1 \
    lib32stdc++6 \
    libcurl4-gnutls-dev \
    libgcc1 \
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
COPY docker_dst_config /app/dst_config
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

# 运行命令
ENTRYPOINT ["./docker-entrypoint.sh"]
