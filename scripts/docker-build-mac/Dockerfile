FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV DST_DIR=/dst-server


# ===== 安装必要依赖 =====
RUN apt update && apt install -y --no-install-recommends \
    wget screen git unzip curl tar build-essential cmake ca-certificates \
    python3 python3-pip tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/lib/apt/lists/*


# ===== 安装 .NET 8 runtime（DepotDownloader 依赖） =====
RUN wget https://packages.microsoft.com/config/ubuntu/22.04/packages-microsoft-prod.deb -O /tmp/packages-microsoft-prod.deb \
    && dpkg -i /tmp/packages-microsoft-prod.deb \
    && apt update \
    && apt install -y dotnet-runtime-8.0 \
    && rm -rf /var/lib/apt/lists/* /tmp/packages-microsoft-prod.deb

# ===== 安装 box64 =====
RUN git clone https://github.com/ptitSeb/box64.git /opt/box64 \
    && mkdir /opt/box64/build && cd /opt/box64/build \
    && cmake .. -DARM_DYNAREC=ON -DCMAKE_BUILD_TYPE=RelWithDebInfo \
    && make -j$(nproc) \
    && make install \
    && rm -rf /opt/box64

# ===== 安装 DepotDownloader 最新 Release =====
RUN wget https://github.com/SteamRE/DepotDownloader/releases/download/DepotDownloader_3.4.0/DepotDownloader-linux-arm64.zip -O /opt/DepotDownloader.zip \
    && unzip /opt/DepotDownloader.zip -d /opt/DepotDownloader \
    && rm /opt/DepotDownloader.zip

# ===== 创建饥荒服务器目录 =====
RUN mkdir -p $DST_DIR

WORKDIR /app

ENV PATH=/usr/local/bin:$PATH

COPY dst-admin-go /app/dst-admin-go
RUN chmod 755 /app/dst-admin-go

COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod 755 /app/docker-entrypoint.sh

COPY config.yml /app/config.yml
COPY docker_dst_config /app/dst_config
COPY dist /app/dist
COPY static /app/static


ENTRYPOINT ["./docker-entrypoint.sh"]

