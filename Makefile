# Makefile
# ===============================
# 项目 Docker Compose 管理脚本
# ===============================

# compose 文件路径
COMPOSE_FILE := deploy/docker/docker-compose.yml

# 容器前缀（可选）
PROJECT_NAME := goadmin

# Docker 命令
DOCKER_COMPOSE := docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --progress plain

# -------------------------------
# 1️⃣ 启动服务
# -------------------------------
.PHONY: up
up:
	@echo "🚀 启动 Docker Compose 服务..."
	$(DOCKER_COMPOSE) up -d
	@echo "✅ 服务已启动"

# -------------------------------
# 2️⃣ 查看日志
# -------------------------------
.PHONY: logs
logs:
	@$(DOCKER_COMPOSE) logs -f

# -------------------------------
# 3️⃣ 停止服务
# -------------------------------
.PHONY: down
down:
	@echo "🧹 停止并移除容器..."
	$(DOCKER_COMPOSE) down
	@echo "✅ 容器已停止"

# -------------------------------
# 4️⃣ 清理所有镜像（慎用）
# -------------------------------
.PHONY: clean
clean: down
	@echo "🔥 清理相关镜像..."
	docker image prune -f
	@echo "✅ 镜像清理完成"

# -------------------------------
# 5️⃣ 清理彻底（包括卷、网络）
# -------------------------------
.PHONY: clean-all
clean-all: down
	@echo "🔥 深度清理：镜像 + 卷 + 网络..."
	docker system prune -a --volumes -f
	@echo "✅ 全部清理完成"


# ==============================================
# 删除 前后端 容器、镜像
# ==============================================
.PHONY: remove
remove: down
	docker rmi -f goadmin_backend:latest goadmin_frontend:latest
	@echo "🧹 清理前后端容器、镜像完成..."

# -------------------------------
# 6️⃣ 状态检查
# -------------------------------
.PHONY: ps
ps:
	$(DOCKER_COMPOSE) ps
