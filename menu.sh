#!/bin/bash
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
    echo "6. 修改端口"
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
            # 停止服务之前检查服务是否存在
            if systemctl is-active --quiet sublink; then
                systemctl stop sublink
            fi
            if systemctl is-enabled --quiet sublink; then
                systemctl disable sublink
            fi
            # 删除服务文件
            if [ -f /etc/systemd/system/sublink.service ]; then
                sudo rm /etc/systemd/system/sublink.service
            fi
            # 删除相关文件和目录
            if [ -d /usr/local/bin/sublink ]; then
                sudo rm -r /usr/local/bin/sublink/*
                sudo rm -r /usr/local/bin/sublink/template
                sudo rm -r /usr/local/bin/sublink/logs
            fi
            if [ -d /usr/bin/sublink ]; then
                sudo rm -r /usr/bin/sublink
            fi
            echo "卸载完成"
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
        6)
            SERVICE_FILE="/etc/systemd/system/sublink.service"
            read -p "请输入新的端口号: " Port
            echo "新的端口号: $Port"
            PARAMETER="run -port $Port"
            # 检查服务文件是否存在
            if [ ! -f "$SERVICE_FILE" ]; then
                echo "服务文件不存在: $SERVICE_FILE"
                exit 1
            fi
            # 检查 ExecStart 是否已经包含该参数
            if grep -q "$PARAMETER" "$SERVICE_FILE"; then
                echo "参数已存在，无需修改。"
            else
                #暂停服务
                systemctl stop sublink
                # 使用 sed 替换 ExecStart 行，添加启动参数
                sudo sed -i "/^ExecStart=/ s|$| $PARAMETER|" "$SERVICE_FILE"
                echo "参数已添加到 ExecStart 行: $PARAMETER"
                
                # 重新加载 systemd 守护进程
                sudo systemctl daemon-reload
                # 重启 sublink 服务
                sudo systemctl restart sublink

                echo "服务已重启。"
            fi
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