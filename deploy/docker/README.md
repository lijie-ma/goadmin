# GoAdmin Docker 部署指南

本目录包含了用于部署 GoAdmin 项目的 Docker 配置文件。

## 项目结构

```
deploy/docker/
├── backend.Dockerfile      # 后端 Go 应用的 Dockerfile
├── frontend.Dockerfile     # 前端 Vue 应用的 Dockerfile
├── nginx.conf             # Nginx 配置文件
├── docker-compose.yml     # Docker Compose 编排文件
└── README.md             # 本文档
```

## 快速开始


### 1. 构建和启动服务

在项目根目录执行：
```bash
cd deploy/docker
docker-compose up -d --build
```

### 3. 查看服务状态

```bash
docker-compose ps
docker-compose logs -f
```

### 4. 停止服务

```bash
docker-compose down
# 如果需要删除数据卷
docker-compose down -v
```

## 服务访问

- **前端应用**: http://localhost
- **后端 API**: http://localhost:8080
- **MySQL 数据库**: localhost:3306
- **Redis**: localhost:6379

## 服务说明

### 前端服务 (frontend)

- 基于 Vue 3 + Vite
- 使用 Nginx 作为 Web 服务器
- 自动代理 `/api/` 请求到后端服务
- 支持 Vue Router History 模式

### 后端服务 (backend)

- Go 1.24 + Gin 框架
- 自动运行数据库迁移
- 支持健康检查
- 包含日志记录

### 数据库服务 (mysql)

- MySQL 8.0
- 字符集: utf8mb4
- 自动初始化数据库和表结构
```
    注意 mysql 默认 password plugin 为 caching_sha2_password, 请修改 mysql root 密码
    alter user 'goadmin'@'%' identified with mysql_native_password by 'Qwert1234';
    flush privileges;
```

### 缓存服务 (redis)

- Redis 7
- 开启持久化 (AOF)

## 生产环境部署建议

### 1. 安全配置

- 修改所有默认密码
- 使用强密码和随机的 JWT_SECRET
- 配置 HTTPS (可以使用 Traefik 或 Nginx SSL)
- 限制数据库和 Redis 的外部访问

### 2. 性能优化

- 根据实际需求调整容器资源限制
- 配置数据库连接池参数
- 开启 Redis 密码认证
- 使用 CDN 加速静态资源

### 3. 监控和日志

- 集成 Prometheus + Grafana 监控
- 使用 ELK Stack 收集日志
- 配置容器健康检查和自动重启

### 4. 备份策略

- 定期备份 MySQL 数据
- 备份 Redis 数据（如果包含重要缓存）
- 保存配置文件版本

## 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| MYSQL_ROOT_PASSWORD | MySQL root 密码 | rootpassword123 |
| MYSQL_DATABASE | 数据库名 | goadmin |
| MYSQL_USER | 数据库用户名 | goadmin |
| MYSQL_PASSWORD | 数据库密码 | goadmin123 |
| REDIS_PASSWORD | Redis 密码 | - |
| JWT_SECRET | JWT 签名密钥 | your-secret-key |
| JWT_EXPIRE | JWT 过期时间(秒) | 7200 |
| SERVER_PORT | 后端服务端口 | 8080 |
| SERVER_MODE | 服务运行模式 | release |
| LOG_LEVEL | 日志级别 | info |
| LOG_OUTPUT | 日志输出方式 | file |

## 目录挂载

### MySQL 数据

- `/var/lib/mysql`: 数据文件
- `/docker-entrypoint-initdb.d`: 初始化SQL脚本

### Redis 数据

- `/data`: Redis持久化数据

### 后端日志

- `/app/logs`: 应用日志目录

## 注意事项

1. 生产环境部署前请修改所有默认密码
2. 建议使用外部负载均衡器管理SSL证书
3. 根据实际需求调整容器资源限制
4. 生产环境中启用Redis密码认证
5. 定期备份数据和配置文件
