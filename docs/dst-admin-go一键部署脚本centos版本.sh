#!/bin/bash
#Install Latest Stable 1Panel Release

osCheck=`uname -a`
if [[ $osCheck =~ 'x86_64' ]];then
    architecture="amd64"
elif [[ $osCheck =~ 'arm64' ]] || [[ $osCheck =~ 'aarch64' ]];then
    architecture="arm64"
else
    echo "暂不支持的系统架构，请参阅官方文档，选择受支持的系统。"
    exit 1
fi

# if [[ ! ${INSTALL_MODE} ]];then
# 	INSTALL_MODE="stable"
# else
#     if [[ ${INSTALL_MODE} != "dev" && ${INSTALL_MODE} != "stable" ]];then
#         echo "请输入正确的安装模式（dev or stable）"
#         exit 1
#     fi
# fi

# VERSION=$(curl -s https://resource.fit2cloud.com/1panel/package/${INSTALL_MODE}/latest)
yum install -y jq
# 使用curl下载版本信息
curl -o version.latest https://api.github.com/repos/hujinbo23/dst-admin-go/releases/latest

# 正则表达式匹配版本号
#VERSION=$(cat version.latest | grep -E 'tag_name\": \"[0-9]+\.[0-9]+\.[0-9]+\.[a-z]+' -o |head -n 1| tr -d 'tag_name\": \"')
VERSION=$(cat version.latest | jq -r .tag_name)

NAME=$(cat version.latest | jq -r .assets[].name)

if [[ "x${VERSION}" == "x" ]];then
    echo "获取最新版本失败，请稍候重试"
    exit 1
fi

echo "开始下载 dst-admin-go ${VERSION} 版本在线安装包"

package_file_name="dst-admin-go-${VERSION}-linux-${architecture}.tgz"

echo "安装包下载名称： ${package_file_name}"

package_download_url="https://github.com/hujinbo23/dst-admin-go/releases/download/${VERSION}/${NAME}"
#package_download_url="https://github.com/hujinbo23/dst-admin-go/releases/download/1.1.6.hotfix/dst-admin-go.1.16.hotfix.tgz"

echo "安装包下载地址： ${package_download_url}"

curl -Lk -o ${package_file_name} ${package_download_url}
#wget -O ${package_file_name} ${package_download_url}
# curl -sfL https://resource.fit2cloud.com/installation-log.sh | sh -s 1p install ${VERSION}
# if [ ! -f ${package_file_name} ];then
# 	echo "下载安装包失败，请稍候重试。"
# 	exit 1
# fi

tar -zxvf ${package_file_name}
if [ $? != 0 ];then
	echo "下载安装包失败，请稍候重试。"
	rm -f ${package_file_name}
	exit 1
fi
#cd dst-admin-go-${VERSION}-linux-${architecture}
DSTPATH=$(pwd)

echo "部署路径: ${DSTPATH}/dst-admin-go.${VERSION}"

cat > /usr/lib/systemd/system/dst-admin-go.service <<EOF
[Unit]
Description=dst-admin-go server daemon
After=network.target

[Service]
Type=simple
WorkingDirectory=${DSTPATH}/dst-admin-go.${VERSION}
ExecStart=${DSTPATH}/dst-admin-go.${VERSION}/dst-admin-go
ExecReload=/bin/kill -HUP \$MAINPID
ExecStop=/bin/kill -s TERM \$MAINPID
KillMode=process

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl start dst-admin-go.service
firewall-cmd --zone=public --add-port=8082/tcp --permanent
firewall-cmd --reload
systemctl enable dst-admin-go.service
systemctl status dst-admin-go.service

HOST_IP=$(ip a | grep inet | grep -v inet6 | grep -v '127.0.0.1' | awk '{print $2}' | awk -F / '{print$1}')
echo "

      请直接访问: http://$HOST_IP:8082 
      默认用户名：admin
      默认密码：123456
      
      
      "