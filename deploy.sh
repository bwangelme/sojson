#!/bin/bash

# SoJSON 部署脚本
# 用于构建、部署和重启 SoJSON 服务

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 项目信息
PROJECT_NAME="sojson"
PROJECT_DIR="$(pwd)"
BINARY_NAME="sojson"
SUPERVISOR_CONFIG_DIR="/opt/homebrew/etc/supervisor.d"
SUPERVISOR_CONFIG_FILE="sojson.ini"

echo "========================================"
echo "🚀 SoJSON 部署脚本"
echo "========================================"

# 检查是否在正确的目录
if [ ! -f "main.go" ]; then
    log_error "未找到 main.go 文件，请在项目根目录运行此脚本"
    exit 1
fi

# 1. 构建项目
log_info "开始构建项目..."
if make build; then
    log_success "项目构建完成"
else
    log_error "项目构建失败"
    exit 1
fi

# 检查二进制文件是否存在
if [ ! -f "$BINARY_NAME" ]; then
    log_error "二进制文件 $BINARY_NAME 不存在"
    exit 1
fi

# 2. 创建日志目录
log_info "创建日志目录..."
mkdir -p log
log_success "日志目录已创建"

# 3. 检查 supervisor 配置文件是否存在
log_info "检查 supervisor 配置文件..."
if [ ! -f "$SUPERVISOR_CONFIG_FILE" ]; then
    log_error "supervisor 配置文件 $SUPERVISOR_CONFIG_FILE 不存在"
    log_info "请确保配置文件存在后再运行部署脚本"
    exit 1
fi

# 4. 检查 supervisor 目录是否存在
if [ ! -d "$SUPERVISOR_CONFIG_DIR" ]; then
    log_error "Supervisor 配置目录不存在: $SUPERVISOR_CONFIG_DIR"
    log_info "请确保已安装 supervisor: brew install supervisor"
    exit 1
fi

# 5. 复制配置文件到 supervisor 目录
log_info "复制配置文件到 supervisor 目录..."
if cp "$SUPERVISOR_CONFIG_FILE" "$SUPERVISOR_CONFIG_DIR/"; then
    log_success "配置文件已复制到 $SUPERVISOR_CONFIG_DIR/"
else
    log_error "配置文件复制失败"
    exit 1
fi

# 6. 重新读取 supervisor 配置
log_info "重新读取 supervisor 配置..."
if supervisorctl reread; then
    log_success "supervisor 配置已重新读取"
else
    log_error "supervisor 配置重新读取失败"
    exit 1
fi

# 7. 更新 supervisor 配置
log_info "更新 supervisor 配置..."
if supervisorctl update; then
    log_success "supervisor 配置已更新"
else
    log_error "supervisor 配置更新失败"
    exit 1
fi

# 8. 重启服务
log_info "重启 $PROJECT_NAME 服务..."
if supervisorctl restart "$PROJECT_NAME" 2>/dev/null; then
    log_success "$PROJECT_NAME 服务已重启"
elif supervisorctl start "$PROJECT_NAME" 2>/dev/null; then
    log_success "$PROJECT_NAME 服务已启动"
else
    log_error "$PROJECT_NAME 服务启动失败"
    log_info "查看服务状态:"
    supervisorctl status "$PROJECT_NAME" || true
    log_info "查看错误日志:"
    supervisorctl tail "$PROJECT_NAME" stderr || true
    exit 1
fi

# 9. 检查服务状态
log_info "检查服务状态..."
sleep 2
if supervisorctl status "$PROJECT_NAME" | grep RUNNING > /dev/null; then
    log_success "$PROJECT_NAME 服务运行正常"
    
    # 显示服务信息
    echo ""
    echo "========================================"
    echo "🎉 部署完成！"
    echo "========================================"
    echo "📋 服务状态:"
    supervisorctl status "$PROJECT_NAME"
    echo ""
    echo "🌐 访问地址:"
    echo "   主页: http://localhost:8080"
    echo "   API:  http://localhost:8080/api"
    echo ""
    echo "📝 日志文件:"
    echo "   输出日志: $PROJECT_DIR/log/$PROJECT_NAME.out.log"
    echo "   错误日志: $PROJECT_DIR/log/$PROJECT_NAME.err.log"
    echo "========================================"
else
    log_error "$PROJECT_NAME 服务启动失败"
    log_info "服务状态:"
    supervisorctl status "$PROJECT_NAME" || true
    exit 1
fi
