#!/bin/bash
set -euo pipefail

########################################
#  此刻 (Cike) 一键部署脚本
#  域名: xuanting.kbti.fun
########################################

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目路径
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
FRONTEND_SRC="${PROJECT_DIR}/frontend"
BACKEND_SRC="${PROJECT_DIR}/backend"

# 部署路径
APP_DIR="/opt/cike"
FRONTEND_DEPLOY="${APP_DIR}/frontend"
BACKEND_DEPLOY="${APP_DIR}/backend"

# 域名配置
DOMAIN="xuanting.kbti.fun"

# 后端地址
SERVER_ADDR=":8080"

########################################
# 辅助函数
########################################

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_ok() {
    echo -e "${GREEN}[OK]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_banner() {
    echo ""
    echo "========================================"
    echo "     此刻 (Cike) 一键部署脚本"
    echo "     域名: ${DOMAIN}"
    echo "========================================"
    echo ""
}

check_command() {
    if ! command -v "$1" &> /dev/null; then
        log_error "未检测到 $1，请先安装: $2"
        exit 1
    fi
}

########################################
# 参数解析
########################################

show_help() {
    cat << 'EOF'
用法: ./deploy.sh [选项]

选项:
  -h, --help       显示帮助信息
  -s, --skip-deps  跳过依赖检查（加速重部署）
  -c, --clean      清理之前的构建产物后部署
  -d, --dry-run    仅检查环境，不执行部署

示例:
  ./deploy.sh              # 完整部署
  ./deploy.sh --skip-deps  # 跳过依赖检查，快速部署
  ./deploy.sh --clean      # 清理后重新部署
EOF
    exit 0
}

SKIP_DEPS=false
CLEAN_BUILD=false
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help) show_help ;;
        -s|--skip-deps) SKIP_DEPS=true; shift ;;
        -c|--clean) CLEAN_BUILD=true; shift ;;
        -d|--dry-run) DRY_RUN=true; shift ;;
        *) log_error "未知参数: $1"; show_help ;;
    esac
done

########################################
# 主逻辑
########################################

print_banner

# 依赖检查
if [ "$SKIP_DEPS" = false ]; then
    log_info "检查必要依赖..."
    check_command "go" "https://go.dev/dl/"
    check_command "node" "https://nodejs.org/"
    check_command "npm" "https://nodejs.org/"
    check_command "caddy" "https://caddyserver.com/docs/install"
    check_command "systemctl" "systemd"
    log_ok "依赖检查通过"
fi

# 版本检查
if [ "$SKIP_DEPS" = false ]; then
    echo "  Go 版本:     $(go version 2>/dev/null | awk '{print $3}')"
    echo "  Node.js 版本: $(node --version 2>/dev/null)"
    echo "  npm 版本:    $(npm --version 2>/dev/null)"
    echo "  Caddy 版本:  $(caddy version 2>/dev/null | head -n1 | awk '{print $1}')"
    echo ""
fi

# 干燥运行模式
if [ "$DRY_RUN" = true ]; then
    log_info "干燥运行模式 - 仅检查环境，不执行部署"
    log_info "项目目录: ${PROJECT_DIR}"
    log_info "前端源目录: ${FRONTEND_SRC}"
    log_info "后端源目录: ${BACKEND_SRC}"
    log_info "前端部署目录: ${FRONTEND_DEPLOY}"
    log_info "后端部署目录: ${BACKEND_DEPLOY}"
    log_info "域名: ${DOMAIN}"
    log_ok "环境检查完成，可以安全部署"
    exit 0
fi

# 清理构建产物
if [ "$CLEAN_BUILD" = true ]; then
    log_info "清理构建产物..."
    rm -rf "${BACKEND_SRC}/cike-server"
    rm -rf "${FRONTEND_SRC}/dist"
    log_ok "构建产物已清理"
fi

# 创建部署目录
log_info "创建部署目录..."
sudo mkdir -p "${FRONTEND_DEPLOY}"
sudo mkdir -p "${BACKEND_DEPLOY}"
sudo mkdir -p /var/lib/cike
sudo mkdir -p /var/log/cike
sudo mkdir -p /var/log/caddy
log_ok "部署目录已创建"

# ======================================
# 步骤 1: 构建并部署后端
# ======================================
echo ""
log_info "步骤 [1/5] 构建并部署后端..."

cd "${BACKEND_SRC}"

log_info "  → 整理 Go 依赖..."
go mod tidy

log_info "  → 编译后端二进制..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o cike-server -ldflags="-s -w" main.go

log_info "  → 复制到部署目录..."
# 如果服务正在运行，先停止再复制（避免 Text file busy）
if systemctl is-active --quiet cike-backend 2>/dev/null; then
    sudo systemctl stop cike-backend
fi
sudo cp cike-server "${BACKEND_DEPLOY}/"
sudo chmod +x "${BACKEND_DEPLOY}/cike-server"

log_ok "后端构建并部署完成"

# ======================================
# 步骤 2: 配置环境变量
# ======================================
echo ""
log_info "步骤 [2/5] 配置环境变量..."

if [ -f /etc/backend.env ]; then
    log_info "  → 已有环境变量文件，保留现有配置"
    log_warn "  → 如需重新生成，请删除 /etc/backend.env 后重新执行"
else
    # 生成随机密钥
    JWT_SECRET=$(openssl rand -base64 32 2>/dev/null || head -c 32 /dev/urandom | base64)
    REFRESH_SECRET=$(openssl rand -base64 32 2>/dev/null || head -c 32 /dev/urandom | base64)

    cat << EOF | sudo tee /etc/backend.env > /dev/null
# 服务端监听地址（Caddy 反向代理到此处）
SERVER_ADDR=${SERVER_ADDR}

# 数据库路径
DATABASE_URL=/var/lib/cike/cike.db

# JWT 密钥（生产环境务必替换为强随机字符串）
JWT_SECRET=${JWT_SECRET}

# Refresh 密钥
REFRESH_SECRET=${REFRESH_SECRET}
EOF
    sudo chmod 600 /etc/backend.env
    log_ok "环境变量文件已生成 (/etc/backend.env)"
fi

# ======================================
# 步骤 3: 构建并部署前端
# ======================================
echo ""
log_info "步骤 [3/5] 构建并部署前端..."

cd "${FRONTEND_SRC}"

log_info "  → 安装前端依赖..."
if ! npm ci; then
    log_warn "  → package-lock.json 与 package.json 不同步，尝试自动修复..."
    npm install
fi

log_info "  → 生产构建..."
npm run build

log_info "  → 复制到部署目录..."
sudo rm -rf "${FRONTEND_DEPLOY}"/*
sudo cp -r dist/* "${FRONTEND_DEPLOY}/"

log_ok "前端构建并部署完成"

# ======================================
# 步骤 4: 配置 systemd 服务
# ======================================
echo ""
log_info "步骤 [4/5] 配置 systemd 服务..."

cat << 'EOF' | sudo tee /etc/systemd/system/cike-backend.service > /dev/null
[Unit]
Description=Cike Backend Server
Documentation=https://github.com/suna0/cike
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root

# 工作目录
WorkingDirectory=/opt/cike/backend

# 启动命令
ExecStart=/opt/cike/backend/cike-server

# 环境变量文件
EnvironmentFile=/etc/backend.env

# 重启策略
Restart=on-failure
RestartSec=5
StartLimitInterval=60s
StartLimitBurst=3

# 资源限制
LimitNOFILE=65535

# 日志输出到 journal
StandardOutput=journal
StandardError=journal
SyslogIdentifier=cike-backend

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload

# 启用服务（首次部署时）
sudo systemctl enable cike-backend 2>/dev/null || true

# 重启服务
sudo systemctl restart cike-backend
log_ok "后端服务已启动"

# 检查服务状态
sleep 1
if systemctl is-active --quiet cike-backend; then
    log_ok "后端服务运行正常"
else
    log_error "后端服务启动失败，请检查日志:"
    log_error "  sudo journalctl -u cike-backend -n 50 --no-pager"
    exit 1
fi

# ======================================
# 步骤 5: 配置 Caddy（如需要）
# ======================================
echo ""
log_info "步骤 [5/5] 检查 Caddy 配置..."

CADDYFILE="/etc/caddy/Caddyfile"

if [ -f "$CADDYFILE" ] && grep -q "${DOMAIN}" "$CADDYFILE" 2>/dev/null; then
    log_info "  → Caddy 配置已存在且包含域名，跳过配置"
else
    log_info "  → 配置 Caddyfile..."

    cat << EOF | sudo tee "$CADDYFILE" > /dev/null
# 全局选项
{
    log {
        output file /var/log/caddy/access.log {
            roll_size 100mb
            roll_keep 10
            roll_keep_for 720h
        }
        format json
    }
}

# ${DOMAIN}
${DOMAIN} {
    # 前端静态文件托管
    root * /opt/cike/frontend
    file_server

    # API 反向代理到 Go 后端
    handle_path /api/* {
        reverse_proxy localhost:8080
    }

    # 可选：压缩
    encode gzip zstd

    # 安全响应头
    header {
        X-Content-Type-Options nosniff
        X-Frame-Options DENY
        X-XSS-Protection "1; mode=block"
        Referrer-Policy strict-origin-when-cross-origin
    }

    # SPA 路由支持
    handle {
        try_files {path} /index.html
    }
}
EOF

    log_ok "Caddy 配置已生成"
fi

# 验证并重启 Caddy
log_info "  → 验证 Caddy 配置..."
if sudo caddy validate --config "$CADDYFILE" 2>/dev/null; then
    log_ok "Caddy 配置验证通过"
else
    log_warn "Caddy 配置验证可能存在问题，继续尝试重载..."
fi

log_info "  → 重载 Caddy..."
sudo systemctl restart caddy || sudo systemctl start caddy

if systemctl is-active --quiet caddy; then
    log_ok "Caddy 运行正常"
else
    log_error "Caddy 启动失败，请检查日志:"
    log_error "  sudo journalctl -u caddy -n 50 --no-pager"
    exit 1
fi

# ======================================
# 部署完成
# ======================================
echo ""
echo "========================================"
log_ok "部署成功完成！"
echo "========================================"
echo ""
echo "  访问地址: https://${DOMAIN}"
echo "  API 地址: https://${DOMAIN}/api"
echo ""
echo "常用命令:"
echo "  查看后端状态: sudo systemctl status cike-backend"
echo "  查看后端日志: sudo journalctl -u cike-backend -f"
echo "  查看 Caddy 状态: sudo systemctl status caddy"
echo "  查看 Caddy 日志: sudo journalctl -u caddy -f"
echo "  重启后端: sudo systemctl restart cike-backend"
echo "  重载 Caddy: sudo systemctl reload caddy"
echo ""
echo "========================================"
