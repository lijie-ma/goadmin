#!/usr/bin/env bash
set -Eeuo pipefail

# ===========================
# 彩色输出 & 时间戳函数
# ===========================
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

log() {
    local level="$1"
    local color="$2"
    shift 2
    local msg="$*"
    printf "%b[%s] [%s]%b %s\n" "$color" "$(date '+%Y-%m-%d %H:%M:%S')" "$level" "$NC" "$msg"
}

# ===========================
# 数据迁移函数
# ===========================
run_migrations() {
    log "INFO" "$GREEN" "Running database migrations..."
    if /app/goadmin migrate up -c /app/config/config.yaml; then
        log "INFO" "$GREEN" "Migrations completed successfully!"
    else
        log "ERROR" "$RED" "Migration failed!"
        return 1
    fi
}

# ===========================
# 主流程
# ===========================
log "INIT" "$YELLOW" "Starting initialization sequence..."


# 执行数据库迁移
if ! run_migrations; then
    log "WARN" "$YELLOW" "Database migration failed. Continuing anyway..."
fi

# 启动主程序
log "START" "$GREEN" "Launching goadmin service..."
exec /app/goadmin control --config /app/config/config.yaml