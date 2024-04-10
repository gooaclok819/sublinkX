#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
plain='\033[0m'

# check root
[[ $EUID -ne 0 ]] && echo -e "${red}错误：${plain} 必须使用root用户运行此脚本！\n" && exit 1

# check os and architecture
arch=$(uname -m)
if [[ $arch == "x86_64" ]]; then
    arch="amd64"
elif [[ $arch == "aarch64" ]]; then
    arch="arm64"
else
    echo -e "${red}不支持的架构：${arch}${plain}\n" && exit 1
fi

install_sublink() {
    echo -e "开始安装 sublink"
    
    # Get the latest version
    latest_version=$(curl --silent "https://api.github.com/repos/gooaclok819/sublinkX/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    download_url="https://github.com/gooaclok819/sublinkX/releases/download/${latest_version}/sublink_${arch}"
    
    wget -N --no-check-certificate -O /usr/local/bin/sublink ${download_url}
    if [[ $? -ne 0 ]]; then
        echo -e "${red}下载 sublink 失败，请确保你的服务器能够下载该文件${plain}"
        exit 1
    fi

    chmod +x /usr/local/bin/sublink
    echo -e "${green}sublink 安装完成${plain}"

    # Create a systemd service file
    cat > /etc/systemd/system/sublink.service << EOF
[Unit]
Description=Sublink Service
After=network.target

[Service]
ExecStart=/usr/local/bin/sublink
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

    # Reload systemd, enable and start the service
    systemctl daemon-reload
    systemctl enable sublink
    systemctl start sublink
}

menu() {
    echo -e "1. 启动服务"
    echo -e "2. 停止服务"
    echo -e "3. 查看安装目录"
    echo -e "4. 退出"
    read -p "请输入你的选择：" choice
    case "$choice" in
        1)
        systemctl start sublink
        ;;
        2)
        systemctl stop sublink
        ;;
        3)
        echo "/usr/local/bin/sublink"
        ;;
        4)
        exit 0
        ;;
        *)
        echo -e "${red}无效的选项！${plain}"
        ;;
    esac
}

echo -e "${green}开始安装${plain}"
install_sublink

# Create a symbolic link to this script
ln -s $(realpath \$0) /usr/local/bin/sublink

while true; do
    menu
done
