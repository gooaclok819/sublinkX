#!/bin/bash

while true; do
    echo "当前版本: 1.6.1"
    
    # 检查服务状态
    status=$(systemctl is-active sublink)
    if [ "$status" = "active" ]; then
        echo "当前运行状态: 已运行"
    else
        echo "当前运行状态: 未运行"
    fi

    # 菜单选项
    PS3="请选择一个选项: "
    options=(
        "启动服务"
        "停止服务"
        "卸载安装"
        "查看服务状态"
        "查看运行目录"
        "修改端口"
        "重置账户密码"
        "退出"
    )

    select opt in "${options[@]}"; do
        case $opt in
            "启动服务")
                echo "正在启动服务..."
                systemctl start sublink && echo "服务已启动" || echo "服务启动失败"
                break
                ;;
            "停止服务")
                echo "正在停止服务..."
                systemctl stop sublink && echo "服务已停止" || echo "停止服务失败"
                break
                ;;
            "卸载安装")
                echo "正在卸载服务..."
                systemctl stop sublink
                systemctl disable sublink
                rm /etc/systemd/system/sublink.service
                systemctl daemon-reload
                rm /usr/bin/sublink
                echo "服务已卸载"
                break
                ;;
            "查看服务状态")
                echo "查看服务状态..."
                systemctl status sublink
                break
                ;;
            "查看运行目录")
                echo "运行目录: /usr/local/bin/sublink"
                echo "需要备份的目录为 db, template 目录为模版文件可备份可不备份"
                break
                ;;
            "修改端口")
                echo -n "请输入新的端口号: "
                read -r new_port
                echo "正在使用新端口启动服务..."
                systemctl stop sublink
                systemctl start sublink run --port "$new_port" && echo "服务已使用新端口 $new_port 启动" || echo "启动失败"
                break
                ;;
            "重置账户密码")
                echo "正在重置账户密码为默认值..."
                systemctl start sublink setting --username admin --password 123456 && echo "账户密码已重置为默认值：admin/123456" || echo "重置密码失败"
                break
                ;;
            "退出")
                exit 0
                ;;
            *)
                echo "无效的选项，请重试"
                break
                ;;
        esac
    done
done
