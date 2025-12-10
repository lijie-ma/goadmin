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

# 执行数据库迁移
run_migrations() {
    echo -e "${GREEN}Running database migrations...${NC}"
    /app/goadmin migrate up -c /app/config/config.yaml
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Migrations completed successfully!${NC}"
    else
        echo -e "${RED}Migration failed!${NC}"
        return 1
    fi
}

# 主流程
main() {

    run_migrations

    echo -e "${GREEN}Database setup completed!${NC}"
}


# 执行主流程
main
