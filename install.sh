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

uninstall_sublink() {
    echo -e "开始卸载 sublink"
    
    # Stop and disable the service
    systemctl stop sublink
    systemctl disable sublink

    # Remove the service file and the program
    rm -f /etc/systemd/system/sublink.service
    rm -f /usr/local/bin/sublink

    # Reload systemd
    systemctl daemon-reload

    echo -e "${green}sublink 卸载完成${plain}"
}

check_status() {
    if systemctl --quiet is-active sublink; then
        echo -e "${green}服务已启动${plain}"
    else
        echo -e "${red}服务未启动${plain}"
    fi
}

menu() {
    echo -e "1. 安装服务"
    echo -e "2. 卸载服务"
    echo -e "3. 启动服务"
    echo -e "4. 停止服务"
    echo -e "5. 查看服务状态"
    echo -e "6. 查看安装目录"
    echo -e "7. 查看运行状态"
    echo -e "8. 退出"
    read -p "请输入你的选择：" choice
    case "$choice" in
        1)
        install_sublink
        ;;
        2)
        uninstall_sublink
        ;;
        3)
        systemctl start sublink
        ;;
        4)
        systemctl stop sublink
        ;;
        5)
        systemctl status sublink
        ;;
        6)
        echo "/usr/local/bin/sublink"
        ;;
        7)
        check_status
        ;;
        8)
        exit 0
        ;;
        *)
        echo -e "${red}无效的选项！${plain}"
        ;;
    esac
}

# Create a symbolic link to this script
ln -s $(realpath \\\$0) /usr/local/bin/sublink

while true; do
    menu
done
