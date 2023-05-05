# 基础镜像
FROM ubuntu:20.04

# 设置维护者信息
LABEL maintainer="your-name <your-email>"

# 更新 apt 软件包索引并安装基础依赖
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# 安装需要的依赖
RUN apt-get update && \
    apt-get install -y libstdc++6:i386 libgcc1:i386 lib32gcc1 lib32stdc++6 libcurl4-gnutls-dev:i386 screen sudo && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 拷贝程序二进制文件
COPY dst-admin-go /app/dst-admin-go

# 拷贝配置文件和静态文件
COPY config.yml /app/config.yml
COPY dst_config /app/dst_config
COPY dist /app/dist

# 暴露端口
EXPOSE 8082/tcp
EXPOSE 10888/udp
EXPOSE 10998/udp
EXPOSE 10999/udp

# 运行命令
ENTRYPOINT ["/app/dst-admin-go"]