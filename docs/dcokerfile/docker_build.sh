#!/bin/bash

# 获取命令行参数
TAG=$1

# 构建镜像
docker build -t hujinbo23/dst-admin-go:$TAG .

# 推送镜像到Docker Hub
docker push hujinbo23/dst-admin-go:$TAG