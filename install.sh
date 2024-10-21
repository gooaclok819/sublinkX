#!/bin/bash

while true; do
    # 获取服务状态
    status=$(systemctl is-active sublink)
    echo "当前版本: 1.6.1"
    
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
    echo "6. 修改端口"  # 新增选项
    echo "7. 重置账户密码"  # 新增选项
    echo "0. 退出"
    echo -n "请选择一个选项: "
    read option

    case $option in
        1)  
            echo "正在启动服务..."
            systemctl start sublink && echo "服务启动命令已执行" || echo "服务启动失败"
            sleep 1
            if [ "$(systemctl is-active sublink)" = "active" ]; then
                echo "服务已成功启动"
            else
                echo "服务启动失败，请检查日志"
            fi
            ;;
        2)
            echo "正在停止服务..."
            systemctl stop sublink && echo "服务已停止" || echo "停止服务失败"
            ;;
        3)
            echo "正在卸载服务..."
            systemctl stop sublink
            systemctl disable sublink
            rm /etc/systemd/system/sublink.service
            systemctl daemon-reload
            rm /usr/bin/sublink
            echo "服务已卸载"
            ;;
        4)
            echo "查看服务状态..."
            systemctl status sublink
            ;;
        5)
            echo "运行目录: /usr/local/bin/sublink"
            echo "需要备份的目录为db,template目录为模版文件可备份可不备份"
            ;;
        6)  # 处理修改端口的选项
            echo "正在停止服务..."
            systemctl stop sublink
            echo -n "请输入新的端口号: "
            read new_port
            echo "正在使用新端口启动服务..."
            systemctl start sublink run --port "$new_port" && echo "服务已使用新端口 $new_port 启动" || echo "启动失败"
            systemctl daemon-reload
            ;;
        7)  # 处理重置账户密码的选项
            echo "正在重置账户密码为默认值..."
            systemctl start sublink setting --username admin --password 123456 && echo "账户密码已重置为默认值：admin/123456" || echo "重置密码失败"
            systemctl daemon-reload
 
