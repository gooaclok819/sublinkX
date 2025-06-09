#!/bin/bash

# 定义颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

INSTALL_DIR="/usr/local/bin/sublink"
BINARY_NAME="sublink"
SERVICE_FILE="/etc/systemd/system/sublink.service"

# --- 辅助函数 ---

# 检查用户是否为root
check_root() {
    if [ "$(id -u)" != "0" ]; then
        echo -e "${RED}该脚本必须以root身份运行。${NC}"
        exit 1
    fi
}

# 执行命令并检查错误
execute_command() {
    local cmd="$1"
    shift
    echo -e "${YELLOW}执行: $cmd $@${NC}"
    output=$("$cmd" "$@" 2>&1)
    if [ $? -ne 0 ]; then
        echo -e "${RED}错误: 命令 '$cmd $@' 执行失败。${NC}"
        echo -e "${RED}输出: ${output}${NC}"
        return 1
    fi
    echo "$output"
    return 0
}

# 获取最新的发行版标签
get_latest_release() {
    curl_output=$(curl --silent "https://api.github.com/repos/gooaclok819/sublinkX/releases/latest")
    if [ $? -ne 0 ]; then
        echo ""
        return 1
    fi
    tag_name=$(echo "$curl_output" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    echo "$tag_name"
    return 0
}

# 获取当前运行的二进制文件版本
get_current_binary_version() {
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        version_output=$(execute_command "$INSTALL_DIR/$BINARY_NAME" "--version")
        if [ $? -eq 0 ]; then
            echo "$version_output"
        else
            echo "未知"
        fi
    else
        echo "未安装"
    fi
}

# 根据机器类型获取二进制文件名
get_binary_file_name() {
    machine_type=$(uname -m)
    if [ "$machine_type" = "x86_64" ]; then
        echo "sublink_amd64"
    elif [ "$machine_type" = "aarch64" ]; then
        echo "sublink_arm64"
    else
        echo ""
        return 1
    fi
    return 0
}

# --- 菜单功能函数 ---

# 启动服务
start_service() {
    echo -e "${GREEN}正在启动服务...${NC}"
    execute_command systemctl start sublink
    execute_command systemctl daemon-reload
    echo -e "${GREEN}服务已启动。${NC}"
}

# 停止服务
stop_service() {
    echo -e "${YELLOW}正在停止服务...${NC}"
    execute_command systemctl stop sublink
    execute_command systemctl daemon-reload
    echo -e "${YELLOW}服务已停止。${NC}"
}

# 查看服务状态
view_status() {
    echo -e "${YELLOW}正在查看服务状态...${NC}"
    execute_command systemctl status sublink
    echo -e "${NC}按任意键继续...${NC}"
    read -n 1 -s
}

# 查看运行目录
view_run_dir() {
    echo -e "${GREEN}运行目录: ${INSTALL_DIR}${NC}"
    echo -e "${YELLOW}需要备份的目录为db,template目录为模版文件可备份可不备份。${NC}"
    echo -e "${NC}按任意键继续...${NC}"
    read -n 1 -s
}

# 修改端口
modify_port() {
    read -p "$(echo -e "${YELLOW}请输入新的端口号: ${NC}")" new_port

    if ! [[ "$new_port" =~ ^[0-9]+$ ]]; then
        echo -e "${RED}无效的端口号，请输入数字。${NC}"
        return
    fi

    echo -e "${YELLOW}新的端口号: $new_port${NC}"

    if [ ! -f "$SERVICE_FILE" ]; then
        echo -e "${RED}服务文件不存在: $SERVICE_FILE。请先运行安装程序。${NC}"
        return
    fi

    # 替换或添加 ExecStart 行中的 --port 参数
    # 查找 ExecStart 行，如果已经有 --port 参数，则替换它；否则，添加它
    if grep -q "ExecStart=.* --port [0-9]\+" "$SERVICE_FILE"; then
        # 已经存在 --port 参数，替换它
        execute_command sed -i -E "s/(--port )[0-9]+/\1$new_port/" "$SERVICE_FILE"
    else
        # 不存在 --port 参数，在 ExecStart 后添加它
        execute_command sed -i -E "s|(ExecStart=.+\s?)|\1run --port $new_port |" "$SERVICE_FILE"
    fi

    echo -e "${YELLOW}重新加载 systemd 守护进程...${NC}"
    execute_command systemctl daemon-reload
    echo -e "${YELLOW}重启 Sublink 服务...${NC}"
    execute_command systemctl restart sublink

    echo -e "${GREEN}端口修改完成，服务已重启。${NC}"
    echo -e "${NC}按任意键继续...${NC}"
    read -n 1 -s
}

# 更新服务
update_service() {
    echo -e "${GREEN}正在检查并更新 Sublink 服务...${NC}"

    latest_release=$(get_latest_release)
    if [ -z "$latest_release" ]; then
        echo -e "${RED}无法获取最新版本信息，更新失败。${NC}"
        echo -e "${NC}按任意键继续...${NC}"
        read -n 1 -s
        return
    fi
    echo -e "${GREEN}最新版本: $latest_release${NC}"

    current_version=$(get_current_binary_version)
    echo -e "${GREEN}当前版本: $current_version${NC}"

    if [ "$current_version" = "$latest_release" ]; then
        echo -e "${GREEN}当前已经是最新版本，无需更新。${NC}"
        echo -e "${NC}按任意键继续...${NC}"
        read -n 1 -s
        return
    fi

    file_name=$(get_binary_file_name)
    if [ -z "$file_name" ]; then
        echo -e "${RED}不支持的机器类型，更新失败。${NC}"
        echo -e "${NC}按任意键继续...${NC}"
        read -n 1 -s
        return
    fi

    echo -e "${YELLOW}停止 Sublink 服务...${NC}"
    execute_command systemctl stop sublink || true # 允许停止失败，因为可能未运行

    echo -e "${YELLOW}下载新版本 $latest_release/$file_name...${NC}"
    temp_file="/tmp/$file_name"
    execute_command curl -L -o "$temp_file" "https://github.com/gooaclok819/sublinkX/releases/download/$latest_release/$file_name"
    if [ $? -ne 0 ]; then
        echo -e "${RED}下载失败，请检查网络连接。${NC}"
        echo -e "${NC}按任意键继续...${NC}"
        read -n 1 -s
        return
    fi

    echo -e "${YELLOW}设置文件可执行权限...${NC}"
    execute_command chmod +x "$temp_file"

    echo -e "${YELLOW}移动文件到 ${INSTALL_DIR}/${BINARY_NAME}...${NC}"
    execute_command mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"

    echo -e "${YELLOW}重新加载 systemd 守护进程...${NC}"
    execute_command systemctl daemon-reload

    echo -e "${YELLOW}启动并启用 Sublink 服务...${NC}"
    execute_command systemctl start sublink
    execute_command systemctl enable sublink

    echo -e "${GREEN}Sublink 服务更新完成。${NC}"
    echo -e "${NC}按任意键继续...${NC}"
    read -n 1 -s
}

# 重置账号密码
reset_account() {
    read -p "$(echo -e "${YELLOW}请输入新的账号: ${NC}")" new_username
    read -p "$(echo -e "${YELLOW}请输入新的密码: ${NC}")" new_password

    if [ -z "$new_username" ] || [ -z "$new_password" ]; then
        echo -e "${RED}账号和密码不能为空。${NC}"
        echo -e "${NC}按任意键继续...${NC}"
        read -n 1 -s
        return
    fi

    echo -e "${YELLOW}正在重置账号密码...${NC}"
    # 调用 Go 程序的 setting 命令
    execute_command "$INSTALL_DIR/$BINARY_NAME" setting --username "$new_username" --password "$new_password"
    
    echo -e "${YELLOW}重启 Sublink 服务以应用更改...${NC}"
    execute_command systemctl restart sublink

    echo -e "${GREEN}账号密码重置成功。${NC}"
    echo -e "${NC}按任意键继续...${NC}"
    read -n 1 -s
}

# 卸载安装
uninstall_service() {
    read -p "$(echo -e "${YELLOW}你是否要卸载 Sublink 服务? (y/n): ${NC}")" confirm_uninstall
    if [ ! "$confirm_uninstall" = "y" ]; then
        echo -e "${NC}取消卸载。${NC}"
        return
    fi

    echo -e "${YELLOW}正在卸载 Sublink 服务...${NC}"

    # 停止服务之前检查服务是否存在
    if systemctl is-active --quiet sublink; then
        echo -e "${YELLOW}停止服务...${NC}"
        execute_command systemctl stop sublink
    fi
    if systemctl is-enabled --quiet sublink; then
        echo -e "${YELLOW}禁用服务...${NC}"
        execute_command systemctl disable sublink
    fi

    read -p "$(echo -e "${YELLOW}是否删除 systemd 服务文件 (包含端口设置)? (y/n): ${NC}")" is_del_systemd
    if [ "$is_del_systemd" = "y" ]; then
        echo -e "${YELLOW}删除 systemd 服务文件...${NC}"
        execute_command rm "$SERVICE_FILE"
        execute_command systemctl daemon-reload # 重新加载以清除服务
    fi

    echo -e "${YELLOW}删除相关文件和目录...${NC}"
    execute_command rm -rf "$INSTALL_DIR/$BINARY_NAME" # 删除二进制文件
    execute_command rm -f "/usr/bin/sublink_installer.sh" # 删除自身脚本

    read -p "$(echo -e "${YELLOW}是否删除模板文件和数据库? (y/n): ${NC}")" is_delete_data
    if [ "$is_delete_data" = "y" ]; then
        echo -e "${YELLOW}删除数据目录...${NC}"
        execute_command rm -rf "$INSTALL_DIR/db"
        echo -e "${YELLOW}删除模板目录...${NC}"
        execute_command rm -rf "$INSTALL_DIR/template"
        echo -e "${YELLOW}删除日志目录...${NC}"
        execute_command rm -rf "$INSTALL_DIR/logs"
    fi

    echo -e "${GREEN}卸载完成。${NC}"
    echo -e "${NC}按任意键继续...${NC}"
    read -n 1 -s
    exit 0 # 卸载后直接退出
}

# --- 主菜单逻辑 ---

show_menu() {
    clear
    current_version=$(get_current_binary_version)
    latest_release=$(get_latest_release)
    if [ -z "$latest_release" ]; then
        latest_release="无法获取"
    fi

    # 获取服务状态
    service_status=$(systemctl is-active sublink 2>/dev/null)
    if [ "$service_status" = "active" ]; then
        display_status="${GREEN}已运行${NC}"
    else
        display_status="${RED}未运行${NC}"
    fi

    echo -e "${YELLOW}--- Sublink 管理菜单 ---${NC}"
    echo -e "最新版本: ${GREEN}${latest_release}${NC}"
    echo -e "当前版本: ${GREEN}${current_version}${NC}"
    echo -e "当前运行状态: ${display_status}"
    echo -e "${GREEN}1. 启动服务${NC}"
    echo -e "${GREEN}2. 停止服务${NC}"
    echo -e "${GREEN}3. 卸载安装${NC}"
    echo -e "${GREEN}4. 查看服务状态${NC}"
    echo -e "${GREEN}5. 查看运行目录${NC}"
    echo -e "${GREEN}6. 修改端口${NC}"
    echo -e "${GREEN}7. 更新${NC}"
    echo -e "${GREEN}8. 重置账号密码${NC}"
    echo -e "${GREEN}0. 退出${NC}"
    echo -n -e "${YELLOW}请选择一个选项: ${NC}"
}

handle_menu_option() {
    read option

    case "$option" in
        1) start_service ;;
        2) stop_service ;;
        3) uninstall_service ;;
        4) view_status ;;
        5) view_run_dir ;;
        6) modify_port ;;
        7) update_service ;;
        8) reset_account ;;
        0) echo -e "${GREEN}退出。${NC}"; exit 0 ;;
        *) echo -e "${RED}无效的选项，请重新选择。${NC}"; sleep 1 ;;
    esac
}

# --- 安装逻辑 ---

install_sublink() {
    check_root

    echo -e "${YELLOW}--- 正在安装 Sublink 服务 ---${NC}"

    latest_release=$(get_latest_release)
    if [ -z "$latest_release" ]; then
        echo -e "${RED}无法获取最新版本信息，安装失败。${NC}"
        exit 1
    fi
    echo -e "${GREEN}最新版本: $latest_release${NC}"

    file_to_download=$(get_binary_file_name)
    if [ -z "$file_to_download" ]; then
        echo -e "${RED}不支持的机器类型: $(uname -m)。安装失败。${NC}"
        exit 1
    fi

    # 创建程序目录
    if [ ! -d "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}创建安装目录: $INSTALL_DIR${NC}"
        execute_command mkdir -p "$INSTALL_DIR"
    fi

    # 下载文件
    echo -e "${YELLOW}下载 Sublink $latest_release 版本到 /tmp/$file_to_download...${NC}"
    temp_download_path="/tmp/$file_to_download"
    execute_command curl -L -o "$temp_download_path" "https://github.com/gooaclok819/sublinkX/releases/download/$latest_release/$file_to_download"
    if [ $? -ne 0 ]; then
        echo -e "${RED}下载失败，请检查网络连接。${NC}"
        exit 1
    fi

    # 设置文件为可执行
    echo -e "${YELLOW}设置可执行权限...${NC}"
    execute_command chmod +x "$temp_download_path"

    # 移动文件到指定目录
    echo -e "${YELLOW}移动二进制文件到 ${INSTALL_DIR}/${BINARY_NAME}...${NC}"
    execute_command mv "$temp_download_path" "$INSTALL_DIR/$BINARY_NAME"

    # 创建 systemctl 服务文件
    echo -e "${YELLOW}创建 systemd 服务文件: $SERVICE_FILE${NC}"
    cat << EOF | tee "$SERVICE_FILE"
[Unit]
Description=Sublink Service
After=network.target

[Service]
ExecStart=$INSTALL_DIR/$BINARY_NAME run --port 8000
WorkingDirectory=$INSTALL_DIR
Restart=always
User=root
Group=root
# AmbientCapabilities=CAP_NET_BIND_SERVICE # 如果ExecStart使用非root用户，可以考虑添加

[Install]
WantedBy=multi-user.target
EOF

    # 重新加载systemd守护进程
    echo -e "${YELLOW}重新加载 systemd 守护进程...${NC}"
    execute_command systemctl daemon-reload

    # 启动并启用服务
    echo -e "${YELLOW}启动并启用 Sublink 服务...${NC}"
    execute_command systemctl start sublink
    execute_command systemctl enable sublink

    # 将此脚本自身复制到 /usr/bin，以便可以全局调用菜单
    echo -e "${YELLOW}将脚本自身复制到 /usr/bin/sublink...${NC}"
    execute_command cp "$0" "/usr/bin/sublink"
    execute_command chmod +x "/usr/bin/sublink"

    echo -e "${GREEN}--- 安装完成 ---${NC}"
    echo -e "${GREEN}默认账号admin, 密码123456 (如果数据库是新的)。${NC}"
    echo -e "${GREEN}默认端口8000。${NC}"
    echo -e "${GREEN}服务已启动并已设置为开机启动。${NC}"
    echo -e "${GREEN}现在输入 'sublink' 可以呼出菜单进行管理。${NC}"
    echo -e "${NC}按任意键进入菜单...${NC}"
    read -n 1 -s
}

# --- 脚本入口点 ---

# 如果没有安装 Go 程序，则执行安装
if [ ! -f "$INSTALL_DIR/$BINARY_NAME" ] || [ ! -f "$SERVICE_FILE" ]; then
    install_sublink
fi

# 循环显示菜单
while true; do
    show_menu
    handle_menu_option
done