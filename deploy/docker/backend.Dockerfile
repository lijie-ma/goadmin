# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app
RUN apk add --no-cache git bash findutils

COPY go.mod go.sum ./
RUN go mod download

# 安装 wire 工具
RUN go install github.com/google/wire/cmd/wire@latest

COPY . .

# 执行 wire 自动生成
RUN find . -type f -name "wire.go" -execdir wire gen \;

# 构建二进制
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o goadmin main.go

# 运行阶段
FROM alpine:3.23

# 使用国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 更新索引并安装依赖（包括 mysql-client 用于数据库初始化）
RUN apk update && \
    apk --no-cache add ca-certificates tzdata mysql-client bash busybox-extras curl

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
COPY deploy/docker/run.sh /app/run.sh

RUN chmod +x /app/run.sh && mkdir -p /app/logs && chown -R appuser:appuser /app
USER appuser
EXPOSE 8080

CMD ["/app/run.sh"]