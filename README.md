# 后台管理系统
## 二次开发基础架构
    - 实现基本的RBAC
    - 集成 goose 实现数据库迁移
    - 使用 jwt 做页面会话标识

## 使用说明
   - go build  -o goadmin main.go
   - goadmin migrate up   初始DB
   - goadmin control -c config/config.yaml 启动后台管理服务
