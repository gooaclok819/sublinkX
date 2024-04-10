#!/bin/bash

# 检查是否以 root 权限运行
if [ "$EUID" -ne 0 ]
  then echo "请以 root 权限运行此脚本."
  exit
fi

# 获取机器型号
ARCH=$(uname -m)

# 根据机器型号下载对应的发行版本
if [ "$ARCH" = "x86_64" ]; then
    FILE_URL=$(curl -s https://api.github.com/repos/gooaclok819/sublinkX/releases/latest | jq -r '.assets[] | select(.name == "sublink_amd64") | .browser_download_url')
elif [ "$ARCH" = "aarch64" ]; then
    FILE_URL=$(curl -s https://api.github.com/repos/gooaclok819/sublinkX/releases/latest | jq -r '.assets[] | select(.name == "sublink_arm64") | .browser_download_url')
else
    echo "不支持的架构."
    exit 1
fi

wget $FILE_URL -O /usr/local/bin/sublink

# 给二进制文件赋予执行权限
chmod 777 /usr/local/bin/sublink

# 创建systemd服务
cat << EOF > /etc/systemd/system/sublink.service
[Unit]
Description=Sublink Service

[Service]
ExecStart=/usr/local/bin/sublink

[Install]
WantedBy=multi-user.target
EOF

# 刷新systemd，使之能够识别新的服务
systemctl daemon-reload

# 设置一个系统变量
echo "alias sublink='bash /usr/local/bin/menu.sh'" >> /root/.bashrc
source /root/.bashrc

# 创建菜单脚本
cat << EOF > /usr/local/bin/menu.sh
#!/bin/bash

function check_status {
    if systemctl -q is-active sublink; then
        echo "当前运行状态: 已运行"
    else
        echo "当前运行状态: 未运行"
    fi
}

check_status

echo "1. 安装并启动"
echo "2. 卸载并退出"
echo "3. 查看服务状态"
echo "4. 退出"

read -p "请输入你的选择: " choice

case $choice in
    1)
        systemctl start sublink
        echo "Sublink已启动."
        ;;
    2)
        systemctl stop sublink
        systemctl disable sublink
        rm /etc/systemd/system/sublink.service
        systemctl daemon-reload
        echo "Sublink已卸载."
        ;;
    3)
        systemctl status sublink
        ;;
    4)
        echo "正在退出..."
        ;;
    *)
        echo "无效的选择."
        ;;
esac

check_status
EOF

chmod +x /usr/local/bin/menu.sh
