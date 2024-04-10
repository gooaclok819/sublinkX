#!/bin/bash

# 获取机器型号
ARCH=$(uname -m)

# 根据机器型号下载对应的发行版本
if [ "$ARCH" = "x86_64" ]; then
    wget https://github.com/gooaclok819/sublinkX/releases/download/latest/sublink_amd64 -O /usr/local/bin/sublink
elif [ "$ARCH" = "aarch64" ]; then
    wget https://github.com/gooaclok819/sublinkX/releases/download/latest/sublink_arm64 -O /usr/local/bin/sublink
else
    echo "不支持的架构."
    exit 1
fi

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
echo "alias sublink='bash /usr/local/bin/menu.sh'" >> ~/.bashrc
source ~/.bashrc

# 创建菜单脚本
cat << EOF > /usr/local/bin/menu.sh
#!/bin/bash

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

if systemctl -q is-active sublink; then
    echo "当前运行状态: 已运行"
else
    echo "当前运行状态: 未运行"
fi
EOF

chmod +x /usr/local/bin/menu.sh
