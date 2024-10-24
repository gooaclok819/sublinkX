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
            systemctl stop sublink
            systemctl disable sublink
            rm /etc/systemd/system/sublink.service
            systemctl daemon-reload
            rm /usr/bin/sublink
            rm -r /usr/local/bin/sublink/template
            rm -r /usr/local/bin/sublink/logs
            rm /usr/local/bin/sublink/sublink
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
            PARAMETER="-port $Port"
            # 检查服务文件是否存在
            if [ ! -f "$SERVICE_FILE" ]; then
                echo "服务文件不存在: $SERVICE_FILE"
                exit 1
            fi

            # 检查 ExecStart 是否已经包含该参数
            if grep -q "$PARAMETER" "$SERVICE_FILE"; then
                echo "参数已存在，无需修改。"
            else
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