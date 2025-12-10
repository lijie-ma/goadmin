# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载 Go 模块依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o goadmin main.go

# 运行阶段
FROM alpine:3.23

# 使用国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 更新索引并安装依赖（包括 mysql-client 用于数据库初始化）
RUN apk update && apk --no-cache add ca-certificates tzdata mysql-client bash

# 设置时区为上海
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 创建非 root 用户
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/goadmin .

# 复制配置文件和必要的资源
COPY --from=builder /app/config ./config
COPY --from=builder /app/internal/i18n/locales ./internal/i18n/locales
COPY --from=builder /app/migrations ./migrations

# 复制并设置初始化脚本
COPY deploy/docker/init-db.sh /app/init-db.sh
RUN chmod +x /app/init-db.sh

# 创建日志目录
RUN mkdir -p /app/logs && \
    chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口（根据配置文件调整）
EXPOSE 8080

# 启动应用
CMD ["./goadmin", "control", "--config", "/app/config/config.yaml"]
