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

#!/bin/bash

# 创建sublink.sh脚本内容并写入到 /usr/bin/sublink
echo '#!/bin/bash
function Select {
    # 获取最新的发行版标签
    latest_release=$(curl --silent "https://api.github.com/repos/gooaclok819/sublinkX/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    # 获取服务状态
    status=$(systemctl is-active sublink)
    echo "最新版本:$latest_release"
    echo "当前版本: 1.6"
    # 判断服务状态并打印
    if [ "$status" = "active" ]; then
        echo "当前运行状态: 已运行"
    else
        echo "当前运行状态: 未运行"
    fi
    echo "1. 启动服务"
    echo "2. 停止服务"
    echo "3. 卸载安装"
    echo "4. 查看服务状态"
    echo "5. 查看运行目录"
    echo "0. 退出"
    echo -n "请选择一个选项: "
    read option

    case $option in
        1)
            systemctl start sublink
            systemctl daemon-reload
            ;;
        2)
            systemctl stop sublink
            systemctl daemon-reload
            ;;
        3)
            systemctl stop sublink
            systemctl disable sublink
            rm /etc/systemd/system/sublink.service
            systemctl daemon-reload
            rm /usr/bin/sublink
            ;;
        4)
            systemctl status sublink
            ;;
        5)
            echo "运行目录: /usr/local/bin/sublink"
            echo "需要备份的目录为db,template目录为模版文件可备份可不备份"
            cd /usr/local/bin/sublink
            echo "已经切换到运行目录"
            ;;
        0)
            exit 0
            ;;
        *)
            echo "无效的选项,请重新选择"
            Select
            ;;
    esac
}

Select
' | sudo tee /usr/bin/sublink

# 将sublink脚本移动到 /usr/bin 并赋予执行权限
sudo chmod 755 /usr/bin/sublink
