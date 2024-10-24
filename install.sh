#!/bin/bash
# 检查用户是否为root
if [ "$(id -u)" != "0" ]; then
    echo -e "${RED}该脚本必须以root身份运行。${NC}"
    exit 1
fi

#创建一个程序目录
cd /usr/local/bin
mkdir sublink
cd sublink

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

# 移动文件到/usr/local/bin
sudo mv $file_name /usr/local/bin/sublink/sublink

# 创建systemctl服务
echo "[Unit]
Description=Sublink Service

[Service]
ExecStart=/usr/local/bin/sublink/sublink
WorkingDirectory=/usr/local/bin/sublink
[Install]
WantedBy=multi-user.target" | sudo tee /etc/systemd/system/sublink.service

# 启动并启用服务
sudo systemctl start sublink
sudo systemctl enable sublink
sudo systemctl daemon-reload
echo "服务已启动并已设置为开机启动"
echo "默认账号admin密码123456 端口8000"
echo "安装完成已经启动输入sudo sublink可以呼出菜单"

sudo curl -o /usr/bin/sublink/menu.sh https://raw.githubusercontent.com/gooaclok819/sublinkX/main/menu.sh
<<<<<<< HEAD
sudo chmod 755 /usr/bin/sublink/menu.sh
=======
sudo chmod 755 /usr/bin/sublink/menu.sh
>>>>>>> 0d0299d75f1e75a282db738f4ae6051ce1949de6
