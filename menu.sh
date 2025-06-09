#!/bin/bash
function Up {
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
    echo "更新完成"

}
function Select {
    # 获取最新的发行版标签
    latest_release=$(curl --silent "https://api.github.com/repos/gooaclok819/sublinkX/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    # 获取服务状态
    cd /usr/local/bin/sublink # 进入sublink目录
    status=$(systemctl is-active sublink)
    version=$(./sublink --version)
    echo "最新版本:$latest_release"
    echo "当前版本:$version"
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
    echo "7. 更新"
    echo "8. 重置账号密码"
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
            sudo rm -r /usr/local/bin/sublink/sublink
            sudo rm -r /usr/bin/sublink
            read -p "是否删除模板文件和数据库(y/n): " isDelete
            if [ "$isDelete" = "y" ]; then
                sudo rm -r /usr/local/bin/sublink/db
                sudo rm -r /usr/local/bin/sublink/template
                sudo rm -r /usr/local/bin/sublink/logs
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
            ;;
        6)
            SERVICE_FILE="/etc/systemd/system/sublink.service"
            read -p "请输入新的端口号: " Port
            echo "新的端口号: $Port"
            PARAMETER="run --port $Port"
            # 检查服务文件是否存在
            if [ ! -f "$SERVICE_FILE" ]; then
                echo "服务文件不存在: $SERVICE_FILE"
                exit 1
            fi

            # 检查 ExecStart 是否已经包含该参数
            if grep -q "run --port" "$SERVICE_FILE"; then
                echo "参数已存在，正在替换..."
                # 使用 sed 替换 ExecStart 行中的 -port 参数
                sudo sed -i "s/-port [0-9]\+/-port $Port/" "$SERVICE_FILE"
            else
                # 如果没有 -port 参数，添加新参数
                # 使用 sed 替换 ExecStart 行，添加启动参数
                sudo sed -i "/^ExecStart=/ s|$| $PARAMETER|" "$SERVICE_FILE"
                echo "参数已添加到 ExecStart 行: $PARAMETER"
            fi

            # 重新加载 systemd 守护进程
            sudo systemctl daemon-reload
            # 重启 sublink 服务
            sudo systemctl restart sublink

            echo "服务已重启。"

            ;;
        7)
            # 停止服务之前检查服务是否存在
            if systemctl is-active --quiet sublink; then
                systemctl stop sublink
            fi
            # 检查是否为最新版本
            if [[ $version = $latest_release ]]; then
                echo "当前已经是最新版本"
            else
                Up
            fi
            ;;
        8)
            read -p "请输入新的账号: " User
            read -p "请输入新的密码: " Password
            # 运行二进制文件并传递启动参数，放在后台运行
            cd /usr/local/bin/sublink
            ./sublink setting --username "$User" --password "$Password" &
            # 获取该程序的PID
            pid=$!
            # 等待程序完成
            wait $pid
            # 如果需要可以在此处进行清理
            systemctl restart sublink
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
