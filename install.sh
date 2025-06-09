#!/bin/bash

# 定义颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

INSTALL_DIR="/usr/local/bin/sublink"
BINARY_NAME="sublink"
SERVICE_FILE="/etc/systemd/system/sublink.service"

GO_VERSION="1.22.0" # <-- 指定的 Go 版本
GO_INSTALL_PATH="/usr/local/go" # Go 安装路径
GO_BIN_PATH="$GO_INSTALL_PATH/bin" # Go 可执行文件路径
REPO_URL="https://github.com/gooaclok819/sublinkX.git"
REPO_DIR="/tmp/sublinkX_repo"
GO_MAIN_PACKAGE_PATH="./" # 你的 main 包在仓库中的相对路径，通常是 "./" 或者 "./cmd/sublink"

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

# 获取最新的发行版标签 (此函数在编译模式下可能不再需要，因为我们直接从源编译)
get_latest_release() {
    # 在编译模式下，我们通常是基于代码库的最新 commit 或特定 tag 编译
    # 如果你仍然想获取release tag，你可以保留此函数，但其用途可能不同
    # 例如，你可以用git describe --tags --abbrev=0来获取最近的tag
    echo "latest_build" # 占位符，表示这是从最新代码编译的
    return 0
}

# 获取当前运行的二进制文件版本 (此函数需要根据编译模式调整)
get_current_binary_version() {
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        # 如果你的Go程序支持 --version 参数并能输出有意义的版本信息
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

# 根据机器类型获取二进制文件名 (此函数在编译模式下不再需要，因为我们直接编译出 sublink)
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

# --- 新增：安装 Go 编译器 ---
install_go() {
    echo -e "${YELLOW}正在安装 Go ${GO_VERSION}...${NC}"

    # 检查是否已安装正确的 Go 版本
    if command -v go &> /dev/null; then
        current_go_version=$(go version | awk '{print $3}' | sed 's/go//')
        if [[ "$current_go_version" == "$GO_VERSION"* ]]; then
            echo -e "${GREEN}Go ${GO_VERSION} 已安装。${NC}"
            return 0
        fi
    fi

    # 移除旧的 Go 安装 (如果存在)
    if [ -d "$GO_INSTALL_PATH" ]; then
        echo -e "${YELLOW}移除旧的 Go 安装...${NC}"
        execute_command sudo rm -rf "$GO_INSTALL_PATH"
    fi

    # 确定架构
    local arch=$(uname -m)
    local go_arch=""
    if [ "$arch" = "x86_64" ]; then
        go_arch="amd64"
    elif [ "$arch" = "aarch64" ]; then
        go_arch="arm64"
    else
        echo -e "${RED}不支持的架构: $arch，无法安装 Go。${NC}"
        return 1
    fi

    local go_tar_file="go${GO_VERSION}.linux-${go_arch}.tar.gz"
    local go_download_url="https://go.dev/dl/$go_tar_file"
    local temp_go_download_path="/tmp/$go_tar_file"

    echo -e "${YELLOW}下载 Go ${GO_VERSION} for ${go_arch}...${NC}"
    execute_command curl -L -o "$temp_go_download_path" "$go_download_url"
    if [ $? -ne 0 ]; then
        echo -e "${RED}下载 Go 失败，请检查网络或 URL。${NC}"
        return 1
    fi

    echo -e "${YELLOW}解压 Go 到 ${GO_INSTALL_PATH}...${NC}"
    execute_command sudo tar -C /usr/local -xzf "$temp_go_download_path"
    execute_command rm "$temp_go_download_path" # 清理临时文件

    # 配置环境变量 (仅为当前会话和未来会话)
    # 确保 /usr/local/go/bin 在 PATH 中
    if ! grep -q "export PATH=.*${GO_BIN_PATH}" /etc/profile; then
        echo -e "${YELLOW}配置 Go 环境变量...${NC}"
        echo "export PATH=\$PATH:${GO_BIN_PATH}" | sudo tee -a /etc/profile
        # 对于root用户，GOPATH通常不需要设置在/root下，但为了通用性可以保留
        echo "export GOPATH=\$HOME/go" | sudo tee -a /etc/profile # 建议设置GOPATH
        echo "export PATH=\$PATH:\$GOPATH/bin" | sudo tee -a /etc/profile # 将GOPATH/bin添加到PATH
    fi
    
    # 立即激活新的环境变量
    source /etc/profile >/dev/null 2>&1 || true # 尝试加载，忽略可能出现的错误

    echo -e "${GREEN}Go ${GO_VERSION} 安装完成。${NC}"
    go version # 再次确认版本
    return 0
}

# --- 菜单功能函数 (保持不变) ---
# ... (start_service, stop_service, view_status, view_run_dir, modify_port, reset_account, uninstall_service 保持不变) ...

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

    # Pay close attention to this line and its indentation.
    # Delete and re-type it if necessary.
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
# 更新服务 (修改为从源代码编译更新)
update_service() {
    echo -e "${GREEN}正在检查并更新 Sublink 服务 (从源代码编译)...${NC}"

    # 在这里，latest_release的含义不再是GitHub Release，而是指最新代码版本
    # 你可以根据需要修改get_latest_release函数来获取Git commit hash等
    # 或者直接假设我们将编译最新的master/main分支
    local latest_version="latest_code" # 编译最新代码

    current_version=$(get_current_binary_version)
    echo -e "${GREEN}当前版本: $current_version${NC}"
    echo -e "${GREEN}目标版本: $latest_version${NC}" # 这里只是个标记

    # 停止服务
    echo -e "${YELLOW}停止 Sublink 服务...${NC}"
    execute_command systemctl stop sublink || true # 允许停止失败

    # 克隆或更新仓库
    if [ -d "$REPO_DIR" ]; then
        echo -e "${YELLOW}进入仓库目录并拉取最新代码...${NC}"
        execute_command cd "$REPO_DIR"
        execute_command git pull
    else
        echo -e "${YELLOW}克隆 Sublink 仓库到 ${REPO_DIR}...${NC}"
        execute_command git clone "$REPO_URL" "$REPO_DIR"
        execute_command cd "$REPO_DIR"
    fi

    # 清理 Go 模块缓存 (重要，防止旧问题)
    echo -e "${YELLOW}执行 go clean -modcache...${NC}"
    execute_command "$GO_BIN_PATH/go" clean -modcache

    # 执行 go mod tidy
    echo -e "${YELLOW}执行 go mod tidy...${NC}"
    execute_command "$GO_BIN_PATH/go" mod tidy
    if [ $? -ne 0 ]; then
        echo -e "${RED}go mod tidy 失败，请检查 Go 模块。${NC}"
        echo -e "${NC}按任意键继续...${NC}"
        read -n 1 -s
        return 1
    fi

    # 编译
    echo -e "${YELLOW}编译 Sublink 二进制文件...${NC}"
    execute_command "$GO_BIN_PATH/go" build -o "$INSTALL_DIR/$BINARY_NAME" "$GO_MAIN_PACKAGE_PATH"
    if [ $? -ne 0 ]; then
        echo -e "${RED}Go 程序编译失败。${NC}"
        echo -e "${NC}按任意键继续...${NC}"
        read -n 1 -s
        return 1
    fi

    # 重新加载 systemd 守护进程
    echo -e "${YELLOW}重新加载 systemd 守护进程...${NC}"
    execute_command systemctl daemon-reload

    # 启动并启用 Sublink 服务
    echo -e "${YELLOW}启动并启用 Sublink 服务...${NC}"
    execute_command systemctl start sublink
    execute_command systemctl enable sublink

    echo -e "${GREEN}Sublink 服务更新完成。${NC}"
    echo -e "${NC}按任意键继续...${NC}"
    read -n 1 -s
    return 0
}


# 重置账号密码
reset_account() {
    read -p "$(echo -e "${YELLOW}请输入新的账号: ${NC}")" new_username
    read -p "$(echo -e "${YELLOW}请输入新的密码: ${NC}")" new_password

    if [ -z "$new_username" ] || [ -z "$new_password" ]; then
        echo -e "${RED}账号和密码不能为空。${NC}"
        echo -e "${NC}按任意键继续...${NC}"
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
    latest_release=$(get_latest_release) # 在此模式下，这可能只是一个占位符
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

# --- 安装逻辑 (修改为从源代码编译安装) ---

install_sublink() {
    check_root

    echo -e "${YELLOW}--- 正在安装 Sublink 服务 (从源代码编译) ---${NC}"

    # 安装 Go 编译器
    install_go
    if [ $? -ne 0 ]; then
        echo -e "${RED}Go 安装失败，无法继续安装 Sublink。${NC}"
        exit 1
    fi

    # 确保 Git 已安装
    if ! command -v git &> /dev/null; then
        echo -e "${RED}Git 未安装，请先安装 Git (sudo apt install git)。${NC}"
        exit 1
    fi

    # 创建程序目录
    if [ ! -d "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}创建安装目录: $INSTALL_DIR${NC}"
        execute_command mkdir -p "$INSTALL_DIR"
    fi

    # 克隆仓库
    echo -e "${YELLOW}克隆 Sublink 仓库到 ${REPO_DIR}...${NC}"
    execute_command git clone "$REPO_URL" "$REPO_DIR"
    if [ $? -ne 0 ]; then
        echo -e "${RED}Git 克隆失败，请检查仓库 URL 或网络。${NC}"
        exit 1
    fi

    # 进入仓库目录
    execute_command cd "$REPO_DIR"

    # 清理 Go 模块缓存 (重要，防止旧问题)
    echo -e "${YELLOW}执行 go clean -modcache...${NC}"
    execute_command "$GO_BIN_PATH/go" clean -modcache

    # 执行 go mod tidy
    echo -e "${YELLOW}执行 go mod tidy...${NC}"
    execute_command "$GO_BIN_PATH/go" mod tidy
    if [ $? -ne 0 ]; then
        echo -e "${RED}go mod tidy 失败，请检查 Go 模块。${NC}"
        exit 1
    fi

    # 编译
    echo -e "${YELLOW}编译 Sublink 二进制文件...${NC}"
    execute_command "$GO_BIN_PATH/go" build -o "$INSTALL_DIR/$BINARY_NAME" "$GO_MAIN_PACKAGE_PATH"
    if [ $? -ne 0 ]; then
        echo -e "${RED}Go 程序编译失败。${NC}"
        exit 1
    fi

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
