FROM debian:buster-slim

LABEL maintainer="hujinbo23 jinbohu23@outlook.com"
LABEL description="DoNotStarveTogehter server panel written in golang.  github: https://github.com/hujinbo23/dst-admin-go"

# 切换到清华大学的 apt 镜像源
RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list \
    && sed -i 's|security.debian.org/debian-security|mirrors.tuna.tsinghua.edu.cn/debian-security|g' /etc/apt/sources.list

# 更新软件包索引
RUN apt-get update

# Install packages
RUN  dpkg --add-architecture i386 \
    && apt-get update \
    && apt-get install -y --no-install-recommends --no-install-suggests  \
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
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 拷贝程序二进制文件
COPY dst-admin-go /app/dst-admin-go
RUN chmod 755 /app/dst-admin-go

COPY run.sh /app/run.sh
RUN chmod 755 /app/run.sh

COPY config.yml /app/config.yml
COPY dst_config /app/dst_config
COPY dist /app/dist
COPY static /app/static

# 内嵌源配置信息

EXPOSE 8082/tcp
EXPOSE 10888/udp
EXPOSE 10998/udp
EXPOSE 10999/udp

# 运行命令
ENTRYPOINT ["./run.sh"]