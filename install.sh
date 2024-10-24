#!/bin/bash
# 检查用户是否为root
if [ "$(id -u)" != "0" ]; then
    echo -e "${RED}该脚本必须以root身份运行。${NC}"
    exit 1
fi

# 创建一个程序目录
INSTALL_DIR="/usr/local/bin/sublink"

if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR"
fi

# 获取最新的发行版标签
latest_release=$(curl --silent "https://api.github.com/repos/gooaclok819/sublinkX/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
echo "最新版本: $latest_release"

# 检测机器类型
machine_type=$(uname -m)

if [ "$machine_type" = "x86_64" ]; then
    file_name="sublink_amd64"
elif [ "$machine_type" = "aarch64" ]; then
    file_name="sublink_arm64"
else
    echo "不支持的机器类型: $machine_type"
    exit 1
fi

# 下载文件
curl -LO "https://github.com/gooaclok819/sublinkX/releases/download/$latest_release/$file_name"

# 设置文件为可执行
chmod +x $file_name

# 移动文件到指定目录
mv $file_name "$INSTALL_DIR/sublink"

# 创建systemctl服务
echo "[Unit]
Description=Sublink Service

[Service]
ExecStart=$INSTALL_DIR/sublink
WorkingDirectory=$INSTALL_DIR
[Install]
WantedBy=multi-user.target" | tee /etc/systemd/system/sublink.service

# 重新加载systemd守护进程
systemctl daemon-reload

# 启动并启用服务
systemctl start sublink
systemctl enable sublink
echo "服务已启动并已设置为开机启动"
echo "默认账号admin密码123456 默认端口8000"
echo "安装完成已经启动输入sublink可以呼出菜单"


# 下载menu.sh并设置权限
curl -o "/usr/bin/sublink" https://raw.githubusercontent.com/gooaclok819/sublinkX/main/menu.sh
chmod 755 "/usr/bin/sublink"
