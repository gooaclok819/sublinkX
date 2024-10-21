# 创建sublink_menu.sh脚本
echo '#!/bin/bash

while true; do
    # 获取服务状态
    status=$(systemctl is-active sublink)
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
    echo "6. 修改端口"
    echo "7. 重置账户密码"
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
            ;;
        6)
            read -p "请输入要修改的端口: " new_port
            /usr/local/bin/sublink/sublink run --port "$new_port"
            ;;
        7)
            read -p "请输入新的用户名 (默认：admin): " new_username
            new_username=${new_username:-admin}
            read -p "请输入新的密码 (默认：123456): " new_password
            new_password=${new_password:-123456}
            /usr/local/bin/sublink/sublink setting --username "$new_username" --password "$new_password"
            ;;
        0)
            exit 0
            ;;
        *)
            echo "无效的选项"
            ;;
    esac
done' > sublink_menu.sh
