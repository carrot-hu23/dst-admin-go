#!/bin/bash
set -e

# 修正最大文件描述符数
ulimit -Sn 10000

echo "Initializing data structure..."

DATA_DIR="/data"

# ===== 数据路径 =====
data_steamcmd="${DATA_DIR}/steamcmd"
data_dst_server="${DATA_DIR}/dst-dedicated-server"
data_backup="${DATA_DIR}/backup"
data_klei="${DATA_DIR}/klei"
data_db_file="${DATA_DIR}/dst-db"
password_file="${DATA_DIR}/password.txt"
first_file="${DATA_DIR}/first"

# ===== 基础目录 =====
mkdir -p "$DATA_DIR"
mkdir -p "$data_backup"
mkdir -p "$data_klei"

# ===== dst-db 文件（不存在则创建）=====
if [ ! -f "$data_db_file" ]; then
  echo "Creating empty dst-db file..."
  touch "$data_db_file"
fi

# ===== password.txt（不存在则初始化默认账号）=====
if [ ! -f "$password_file" ]; then
  echo "Initializing default admin account..."
  cat > "$password_file" <<EOF
username=admin
password=123456
displayName=admin
photoURL=xxx
EOF
fi

# ===== klei 目录映射 =====
mkdir -p /root/.klei
ln -sf "$data_klei" /root/.klei/DoNotStarveTogether

# ===== steamcmd 判断 =====
if [ -d /app/steamcmd ] && [ ! -L /app/steamcmd ]; then
  echo "Using user mounted steamcmd: /app/steamcmd"
  steam_cmd_path="/app/steamcmd"
else
  echo "Using data steamcmd: /data/steamcmd"
  mkdir -p "$data_steamcmd"
  ln -sf "$data_steamcmd" /app/steamcmd
  steam_cmd_path="$data_steamcmd"
fi

# ===== dst server 判断 =====
if [ -d /app/dst-dedicated-server ] && [ ! -L /app/dst-dedicated-server ]; then
  echo "Using user mounted dst server: /app/dst-dedicated-server"
  steam_dst_server="/app/dst-dedicated-server"
else
  echo "Using data dst server: /data/dst-dedicated-server"
  mkdir -p "$data_dst_server"
  ln -sf "$data_dst_server" /app/dst-dedicated-server
  steam_dst_server="$data_dst_server"
fi

# ===== 其他软链接 =====
ln -sf "$data_backup" /app/backup
ln -sf "$data_db_file" /app/dst-db
ln -sf "$password_file" /app/password.txt

# ⚠️ first 不自动创建（由程序初始化后生成）
ln -sf "$first_file" /app/first

# ============================================================
# 安装 SteamCMD（如果不存在）
# ============================================================

cd "$steam_cmd_path"

retry=1
while [ ! -e "${steam_cmd_path}/steamcmd.sh" ]; do
  if [ $retry -gt 3 ]; then
    echo "Download steamcmd failed after three times"
    exit -2
  fi

  echo "Installing steamcmd, try: ${retry}"
  wget http://media.steampowered.com/installer/steamcmd_linux.tar.gz -P "$steam_cmd_path"
  tar -zxvf "$steam_cmd_path/steamcmd_linux.tar.gz" -C "$steam_cmd_path"
  sleep 3
  ((retry++))
done

# ============================================================
# 安装 DST Dedicated Server（如果不存在）
# ============================================================

retry=1
while [ ! -e "${steam_dst_server}/bin/dontstarve_dedicated_server_nullrenderer" ]; do
  if [ $retry -gt 3 ]; then
    echo "Download DST server failed after three times"
    exit -2
  fi

  echo "Installing DST server, try: ${retry}"
  bash "$steam_cmd_path/steamcmd.sh" \
    +force_install_dir "$steam_dst_server" \
    +login anonymous \
    +app_update 343050 validate \
    +quit

  sleep 3
  ((retry++))
done

echo "SteamCMD ready at $steam_cmd_path"
echo "DST server ready at $steam_dst_server"

cd /app
exec ./dst-admin-go
