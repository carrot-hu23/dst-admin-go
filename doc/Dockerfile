FROM debian:stable-slim

RUN sed -i 's|deb.debian.org|mirrors.aliyun.com|g' /etc/apt/sources.list && \
    sed -i 's|security.debian.org|mirrors.aliyun.com/debian-security|g' /etc/apt/sources.list && \
    dpkg --add-architecture i386

RUN apt-get update && apt-get install -y --no-install-recommends \
    lib32stdc++6 \
    libgcc1 \
    libcurl4-gnutls-dev:i386 \
    wget \
    screen \
    procps && \
    rm -rf /var/lib/apt/lists/*


# 设置工作目录
WORKDIR /app

# 拷贝程序二进制文件
COPY dst-admin-go /app/dst-admin-go
RUN chmod 755 /app/dst-admin-go

COPY run.sh /app/run.sh
RUN chmod 755 /app/run.sh

COPY ../config.yml /app/config.yml
COPY ../dst_config /app/dst_config
COPY dist /app/dist
COPY ../static /app/static

# 内嵌源配置信息

EXPOSE 8082/tcp
EXPOSE 10888/udp
EXPOSE 10998/udp
EXPOSE 10999/udp

# 运行命令
ENTRYPOINT ["./run.sh"]