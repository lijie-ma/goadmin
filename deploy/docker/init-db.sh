#!/bin/bash

# 数据库初始化脚本
# 用于检查和初始化数据库

set -e

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting database initialization...${NC}"

# 等待数据库就绪
wait_for_db() {
    echo -e "${YELLOW}Waiting for MySQL to be ready...${NC}"
    until mysql -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" &>/dev/null; do
        echo -n "."
        sleep 2
    done
    echo -e "${GREEN}MySQL is ready!${NC}"
}

# 检查数据库是否已初始化
check_db_initialized() {
    local table_count=$(mysql -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" \
        -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='$DB_NAME'" \
        -sN 2>/dev/null || echo "0")

    if [ "$table_count" -gt "0" ]; then
        return 0  # 已初始化
    else
        return 1  # 未初始化
    fi
}

# 执行SQL文件
execute_sql_file() {
    local sql_file=$1
    echo -e "${YELLOW}Executing: $sql_file${NC}"
    mysql -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$sql_file"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Successfully executed: $sql_file${NC}"
    else
        echo -e "${RED}Failed to execute: $sql_file${NC}"
        return 1
    fi
}

# 执行数据库迁移
run_migrations() {
    echo -e "${GREEN}Running database migrations...${NC}"
    /app/goadmin migrate up
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Migrations completed successfully!${NC}"
    else
        echo -e "${RED}Migration failed!${NC}"
        return 1
    fi
}

# 主流程
main() {
    # 等待数据库就绪
    wait_for_db

    # 执行数据库迁移
    run_migrations

    echo -e "${GREEN}Database setup completed!${NC}"
}

# 设置环境变量默认值
: ${DB_HOST:=mysql}
: ${DB_PORT:=3306}
: ${DB_USER:=goadmin}
: ${DB_PASSWORD:=Qwert1234}
: ${DB_NAME:=goadmin}

# 执行主流程
main
