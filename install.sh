#!/bin/bash

# 颜色设置
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查用户是否为root
if [ "$(id -u)" != "0" ]; then
    echo -e "${RED}该脚本必须以root身份运行。${NC}"
    exit 1
fi

# 定义安装路径和Go版本
INSTALL_APP_DIR="/usr/local/sublink"       # Sublink 应用程序的实际安装目录 (包含二进制文件)
INSTALL_CMD_PATH="/usr/local/bin/sublink"  # Sublink 快捷命令的路径 (菜单脚本)
REPO_URL="https://github.com/gooaclok819/sublinkX.git"
GO_VERSION="1.22.0"
GO_TAR_GZ="go${GO_VERSION}.linux-amd64.tar.gz"
GO_DOWNLOAD_URL="https://go.dev/dl/${GO_TAR_GZ}"

echo -e "${YELLOW}--- 开始 Sublink 自动化安装与菜单部署脚本 ---${NC}"

---

## 清理旧的安装和缓存

echo -e "${YELLOW}清理旧的Sublink安装目录和Go缓存...${NC}"

# 清理主应用程序目录
if [ -d "$INSTALL_APP_DIR" ]; then
    echo -e "${GREEN}删除旧的主安装目录：$INSTALL_APP_DIR${NC}"
    rm -rf "$INSTALL_APP_DIR"
fi

# 清理快捷命令路径，可能是文件也可能是目录
if [ -e "$INSTALL_CMD_PATH" ]; then # -e 检查文件或目录是否存在
    if [ -f "$INSTALL_CMD_PATH" ]; then # -f 检查是否是文件
        echo -e "${GREEN}删除旧的快捷命令文件：$INSTALL_CMD_PATH${NC}"
        rm -f "$INSTALL_CMD_PATH"
    elif [ -d "$INSTALL_CMD_PATH" ]; then # -d 检查是否是目录
        echo -e "${GREEN}删除旧的快捷命令目录：$INSTALL_CMD_PATH${NC}"
        rm -rf "$INSTALL_CMD_PATH"
    fi
fi

# 清理Go模块缓存 (使用 || true 避免在Go未安装时报错)
go clean -modcache || true
echo -e "${GREEN}已清理Go模块缓存（如果存在）${NC}"

---

## 安装或更新 Go 语言环境

echo -e "${YELLOW}安装或更新 Go ${GO_VERSION}...${NC}"

# 检查Go是否已安装并且版本符合要求
if command -v go &> /dev/null; then
    current_go_version=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$current_go_version" = "$GO_VERSION" ]; then
        echo -e "${GREEN}Go ${GO_VERSION} 已安装且版本正确。${NC}"
    else
        echo -e "${YELLOW}检测到Go版本 ${current_go_version}，正在安装 Go ${GO_VERSION}...${NC}"
        # 移除旧的Go安装
        rm -rf /usr/local/go
        # 下载并安装指定Go版本
        wget -q --show-progress "$GO_DOWNLOAD_URL" -O "/tmp/$GO_TAR_GZ"
        tar -C /usr/local -xzf "/tmp/$GO_TAR_GZ"
        rm "/tmp/$GO_TAR_GZ"
        echo -e "${GREEN}Go ${GO_VERSION} 安装完成。${NC}"
    fi
else
    echo -e "${YELLOW}未检测到Go环境，正在安装 Go ${GO_VERSION}...${NC}"
    wget -q --show-progress "$GO_DOWNLOAD_URL" -O "/tmp/$GO_TAR_GZ"
    tar -C /usr/local -xzf "/tmp/$GO_TAR_GZ"
    rm "/tmp/$GO_TAR_GZ"
    echo -e "${GREEN}Go ${GO_VERSION} 安装完成。${NC}"
fi

# 配置Go环境变量
export PATH=$PATH:/usr/local/go/bin
# 永久添加到 /etc/profile (系统范围)
if ! grep -q "/usr/local/go/bin" /etc/profile; then
    echo "export PATH=\$PATH:/usr/local/go/bin" >> /etc/profile
fi
# 永久添加到 ~/.bashrc (当前用户)
if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
fi
source /etc/profile # 立即生效环境变量
echo -e "${GREEN}Go环境变量已配置。${NC}"
go version

---

## 克隆并构建 Sublink

echo -e "${YELLOW}克隆 sublinkX 仓库并进行构建...${NC}"
mkdir -p "$INSTALL_APP_DIR"
if ! git clone "$REPO_URL" "$INSTALL_APP_DIR"; then
    echo -e "${RED}克隆仓库失败。请检查网络连接或仓库URL。${NC}"
    exit 1
fi

cd "$INSTALL_APP_DIR" || { echo -e "${RED}无法进入安装目录：$INSTALL_APP_DIR${NC}"; exit 1; }

echo -e "${YELLOW}初始化Go模块并构建Sublink...${NC}"
go mod tidy
if ! go build -o sublink; then
    echo -e "${RED}Sublink 构建失败。请检查依赖或代码问题。${NC}"
    exit 1
fi
echo -e "${GREEN}Sublink 构建成功。可执行文件位于：$INSTALL_APP_DIR/sublink${NC}"

---

## 创建 Systemd 服务

echo -e "${YELLOW}创建 systemctl 服务...${NC}"
# 使用 here-document 正确创建 systemd 服务文件
cat <<EOF | tee /etc/systemd/system/sublink.service
[Unit]
Description=Sublink Service
After=network.target

[Service]
ExecStart=$INSTALL_APP_DIR/sublink
WorkingDirectory=$INSTALL_APP_DIR
Restart=always
# 确保Go的可执行文件路径在服务启动时可用
Environment="PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

[Install]
WantedBy=multi-user.target
EOF

# 重新加载systemd守护进程
systemctl daemon-reload

# 停止、启动并启用服务
systemctl stop sublink &> /dev/null || true # 尝试停止，如果服务不存在则不报错
systemctl start sublink
systemctl enable sublink

echo -e "${GREEN}Sublink 服务已启动并已设置为开机启动。${NC}"
echo -e "${GREEN}默认账号 ${YELLOW}admin${NC}，密码 ${YELLOW}123456${NC}，默认端口 ${YELLOW}8000${NC}"

---

## 创建 Sublink 菜单命令脚本

echo -e "${YELLOW}创建 sublink 菜单命令脚本...${NC}"

# 写入菜单脚本到 /usr/local/bin/sublink
# 使用 'EOF' (带单引号) 来防止变量展开和内部语法问题
cat << 'EOF' | tee "$INSTALL_CMD_PATH"
#!/bin/bash

# 颜色设置
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 定义应用程序主目录 (与安装脚本中的 INSTALL_APP_DIR 保持一致)
APP_BASE_DIR="/usr/local/sublink"

# 定义更新函数
Up() {
    echo -e "${YELLOW}正在更新 Sublink...${NC}"
    
    # 停止服务
    systemctl stop sublink

    # 进入应用目录
    cd "$APP_BASE_DIR" || { echo -e "${RED}无法进入应用目录：$APP_BASE_DIR${NC}"; return 1; }

    # 拉取最新代码
    git pull origin main || { echo -e "${RED}拉取最新代码失败。${NC}"; return 1; }

    # 重新构建
    go mod tidy
    go build -o sublink || { echo -e "${RED}更新后构建失败。${NC}"; return 1; }

    # 重新加载systemd守护进程
    systemctl daemon-reload

    # 启动服务
    systemctl start sublink
    echo -e "${GREEN}Sublink 更新完成并已重新启动。${NC}"
}

# 定义菜单函数
Select() {
    while true; do
        clear # 清屏以保持菜单整洁
        echo -e "${BLUE}--- Sublink 管理菜单 ---${NC}"
        
        # 获取最新的发行版标签
        latest_release=$(curl --silent "https://api.github.com/repos/gooaclok819/sublinkX/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "获取失败")
        
        # 尝试获取当前版本（如果 Sublink 有一个查看版本号的命令）
        current_version="未知或待实现" # 默认为未知
        if [ -f "$APP_BASE_DIR/sublink" ]; then
            # 如果你的 sublink 二进制文件有 --version 或类似的命令，可以修改这里
            # 例如：current_version=$("$APP_BASE_DIR/sublink" --version | head -n 1)
            : # No-op
        fi


        # 获取服务状态
        status=$(systemctl is-active sublink)

        echo -e "${BLUE}--------------------------------${NC}"
        echo -e "最新版本: ${GREEN}$latest_release${NC}"
        echo -e "当前版本: ${GREEN}$current_version${NC}" # 这里需要根据实际情况获取
        if [ "$status" = "active" ]; then
            echo -e "当前运行状态: ${GREEN}已运行${NC}"
        else
            echo -e "当前运行状态: ${RED}未运行${NC}"
        fi
        echo -e "${BLUE}--------------------------------${NC}"
        echo "1. 启动服务"
        echo "2. 停止服务"
        echo "3. 卸载安装"
        echo "4. 查看服务状态"
        echo "5. 查看运行目录"
        echo "6. 修改端口"
        echo "7. 更新 (拉取最新代码并重新构建)" # 修正此行，确保括号在 heredoc 内被正确处理
        echo "0. 退出"
        echo -e "${BLUE}--------------------------------${NC}"
        echo -n "请选择一个选项: "
        read option

        case $option in
            1)
                echo -e "${YELLOW}正在启动 Sublink 服务...${NC}"
                systemctl start sublink
                systemctl daemon-reload
                echo -e "${GREEN}服务启动命令已发出。${NC}"
                sleep 2 # 等待服务启动
                ;;
            2)
                echo -e "${YELLOW}正在停止 Sublink 服务...${NC}"
                systemctl stop sublink
                systemctl daemon-reload
                echo -e "${GREEN}服务停止命令已发出。${NC}"
                sleep 2 # 等待服务停止
                ;;
            3)
                echo -e "${RED}警告：您选择了卸载安装。这将删除Sublink程序和服务。${NC}"
                read -p "确定要继续吗? (y/N): " confirm_uninstall
                if [[ "$confirm_uninstall" =~ ^[Yy]$ ]]; then
                    # 停止服务之前检查服务是否存在
                    if systemctl is-active --quiet sublink; then
                        echo -e "${YELLOW}停止 Sublink 服务...${NC}"
                        systemctl stop sublink
                    fi
                    if systemctl is-enabled --quiet sublink; then
                        echo -e "${YELLOW}禁用 Sublink 服务...${NC}"
                        systemctl disable sublink
                    fi
                    
                    # 删除服务文件
                    if [ -f /etc/systemd/system/sublink.service ]; then
                        echo -e "${YELLOW}删除 systemd 服务文件...${NC}"
                        rm /etc/systemd/system/sublink.service
                        systemctl daemon-reload
                    fi
                    
                    # 删除快捷命令文件
                    if [ -f "$INSTALL_CMD_PATH" ]; then
                        echo -e "${YELLOW}删除快捷命令文件：$INSTALL_CMD_PATH${NC}"
                        rm "$INSTALL_CMD_PATH"
                    fi

                    read -p "是否保留数据 (db, logs, template目录)? (y/n): " isDeleteData
                    if [[ "$isDeleteData" =~ ^[Nn]$ ]]; then # 如果选择 'n' (不保留)
                        echo -e "${YELLOW}删除 Sublink 数据目录...${NC}"
                        if [ -d "$APP_BASE_DIR/db" ]; then rm -rf "$APP_BASE_DIR/db"; fi
                        if [ -d "$APP_BASE_DIR/template" ]; then rm -rf "$APP_BASE_DIR/template"; fi
                        if [ -d "$APP_BASE_DIR/logs" ]; then rm -rf "$APP_BASE_DIR/logs"; fi
                    else # 默认保留，或者用户输入 'y'
                        echo -e "${GREEN}数据目录将被保留。${NC}"
                    fi

                    # 删除主应用程序目录
                    if [ -d "$APP_BASE_DIR" ]; then
                        echo -e "${YELLOW}删除主应用程序目录：$APP_BASE_DIR${NC}"
                        rm -rf "$APP_BASE_DIR"
                    fi
                    
                    echo -e "${GREEN}Sublink 卸载完成。${NC}"
                    exit 0 # 卸载完成后退出脚本
                else
                    echo -e "${GREEN}卸载操作已取消。${NC}"
                fi
                ;;
            4)
                echo -e "${YELLOW}正在查看 Sublink 服务状态...${NC}"
                systemctl status sublink
                read -p "按任意键继续..."
                ;;
            5)
                echo -e "${BLUE}Sublink 运行目录: ${GREEN}$APP_BASE_DIR${NC}"
                echo -e "${YELLOW}需要备份的目录通常为 ${GREEN}db${NC}。${GREEN}template${NC} 目录为模板文件，可选择备份。${NC}"
                echo -e "${YELLOW}当前目录已切换到：${GREEN}$APP_BASE_DIR${NC}"
                cd "$APP_BASE_DIR" || echo -e "${RED}切换目录失败。${NC}"
                read -p "按任意键继续..."
                ;;
            6)
                SERVICE_FILE="/etc/systemd/system/sublink.service"
                read -p "请输入新的端口号 (例如: 8000): " Port
                # 检查输入是否为有效数字
                if ! [[ "$Port" =~ ^[0-9]+$ ]] || [ "$Port" -lt 1 ] || [ "$Port" -gt 65535 ]; then
                    echo -e "${RED}无效的端口号。请输入1到65535之间的数字。${NC}"
                    read -p "按任意键继续..."
                    continue
                fi
                echo -e "新的端口号: ${GREEN}$Port${NC}"
                
                # 检查服务文件是否存在
                if [ ! -f "$SERVICE_FILE" ]; then
                    echo -e "${RED}服务文件不存在: $SERVICE_FILE。无法修改端口。${NC}"
                    read -p "按任意键继续..."
                    continue
                fi

                # 暂停服务
                systemctl stop sublink
                echo -e "${YELLOW}服务已暂停，正在修改端口配置...${NC}"
                
                # 构建新的 ExecStart 参数
                # 先移除旧的 -port 参数，然后添加新的
                sed -i -E "s/ -port [0-9]+//g" "$SERVICE_FILE" # 移除所有旧的 -port 参数
                sudo sed -i "/^ExecStart=/ s|\$| -port $Port|" "$SERVICE_FILE" # 在行尾添加新的 -port 参数
                
                echo -e "${GREEN}端口已修改为 $Port。${NC}"

                # 重新加载 systemd 守护进程
                sudo systemctl daemon-reload
                # 重启 sublink 服务
                sudo systemctl restart sublink
                echo -e "${GREEN}服务已重启。${NC}"
                read -p "按任意键继续..."
                ;;
            7)
                Up # 调用更新函数
                read -p "按任意键继续..."
                ;;
            0)
                echo -e "${BLUE}退出 Sublink 管理菜单。${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}无效的选项，请重新选择。${NC}"
                sleep 1
                ;;
        esac
    done
}

# 在脚本的最后调用菜单函数
Select
EOF

# 设置菜单脚本为可执行
chmod +x "$INSTALL_CMD_PATH"
echo -e "${GREEN}Sublink 菜单命令脚本已创建并设置为可执行：${INSTALL_CMD_PATH}${NC}"

echo -e "${GREEN}--- Sublink 安装与菜单部署完成！---${NC}"
echo -e "${YELLOW}请注意：为了确保 PATH 环境变量在当前会话中完全生效，您可能需要：${NC}"
echo -e "${YELLOW}1. 运行 'source ~/.bashrc' 或 'source /etc/profile'${NC}"
echo -e "${YELLOW}2. 或者，最简单和可靠的方法是重新登录您的SSH会话或打开一个新的终端。${NC}"
echo -e "${YELLOW}重新登录后，输入 ${GREEN}sublink${NC} 即可呼出菜单。${NC}"
