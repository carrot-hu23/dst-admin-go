FROM debian:buster-slim

LABEL maintainer="hujinbo23 jinbohu23@outlook.com"
LABEL description="DoNotStarveTogehter server panel written in golang.  github: https://github.com/hujinbo23/dst-admin-go"

# Install packages
RUN echo "deb http://mirrors.aliyun.com/debian/ buster main contrib non-free" > /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian/ buster-updates main contrib non-free" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian/ buster-backports main contrib non-free" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian-security buster/updates main contrib non-free" >> /etc/apt/sources.list && \
   # dpkg --add-architecture i386 && \
   # apt-get update && \
   # apt-get install -y lib32stdc++6 libgcc1 libcurl4-gnutls-dev:i386 wget screen procps
    dpkg --add-architecture i386 \
    && apt-get update \
    && apt-get install -y --no-install-recommends --no-install-suggests  \
        libstdc++6:i386 \
        libgcc1:i386 \
        lib32gcc1 \
        lib32stdc++6 \
        libcurl4-gnutls-dev:i386 \
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