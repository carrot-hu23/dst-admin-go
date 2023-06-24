#!/bin/bash
# dst-go.sh

# 下载并解压dst-admin-go
download() {
  if command -v wget > /dev/null
  then
    # 执行wget命令
    echo "Downloading dst-admin-go..."
    wget https://github.com/hujinbo23/dst-admin-go/releases/download/2.0.0.beta/dst-admin-go.tgz
    tar -xvf dst-admin-go.tgz
    cd dst-admin-go
    chmod +x dst-admin-go
  else
    echo "wget command not found."
  fi

}

# 检查dst-admin-go进程是否运行
check_status() {
  if pgrep dst-admin-go > /dev/null
  then
    echo "dst-admin-go is running."
  else
    echo "dst-admin-go is not running."
  fi
}

# 启动dst-admin-go进程
start() {
  echo "Starting dst-admin-go..."
  nohup ./dst-admin-go > /dev/null 2>&1 &
}

# 关闭dst-admin-go进程
stop() {
  echo "Stopping dst-admin-go..."
  pkill dst-admin-go
}

# 显示菜单
# 显示菜单
menu() {
  echo "Please select an option:"
  echo "0. Download dst-admin-go"
  echo "1. Check status"
  echo "2. Start"
  echo "3. Stop"
  read option
  case $option in
    0) download ;;
    1) check_status ;;
    2) start ;;
    3) stop ;;
    *) echo "Invalid option. Please try again." ;;
  esac
}

menu
