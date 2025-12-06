# 后台管理系统

## 二次开发基础架构
    - 实现基本的RBAC
    - 集成 goose 实现数据库迁移 支持mysql pg
    - 使用 jwt 做页面会话标识
    - 框架代码来自 ai + 手动改写

## 使用说明

### 定时任务
    参考 internal/cron/cron.go

### 守护
    实现 pkg/task/task.go 中的 Service 接口即可
        实例参考 cmd/server/cron.go
    注册到启动器中
```
    services := task.NewServiceManager()
    services.AddService(server.NewCronManager(), server.NewWebServer(cfg))
```

### 编译
   - go build  -o goadmin main.go
   - goadmin migrate up   初始DB
   - goadmin control -c config/config.yaml 启动后台管理服务

### 页面
```
cd public/src

npm install

# 启动开发服务器
npm run dev
```
    http://localhost:3000
    用户名：admin
    密码：123456

