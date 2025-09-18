# 日志模块 (Logger Module)

基于 Zap 实现的高性能、灵活的日志模块，提供了多种配置选项和日志记录方法。

## 特性

- 高性能：基于 Uber 的 Zap 日志库实现
- 多级别日志：支持 Debug、Info、Warn、Error、DPanic、Panic、Fatal 等级别
- 灵活配置：支持 JSON/文本格式、控制台/文件输出、日志轮转等
- 结构化日志：支持添加固定字段和临时字段
- 上下文感知：可以从上下文中提取信息
- 全局实例：支持在应用的任何地方使用全局日志实例
- 选项模式：使用函数选项模式实现灵活配置
- 接口设计：通过接口抽象隐藏底层实现细节

## 快速开始

### 基本使用

```go
package main

import "goadmin/pkg/logger"

func main() {
    // 创建默认配置的日志记录器
    log := logger.New()

    // 记录不同级别的日志
    log.Debug("这是一条调试日志")
    log.Info("这是一条信息日志")
    log.Warn("这是一条警告日志")
    log.Error("这是一条错误日志")

    // 使用格式化日志
    log.Infof("用户 %s 登录系统，IP: %s", "admin", "127.0.0.1")
}
```

### 自定义配置

```go
package main

import "goadmin/pkg/logger"

func main() {
    // 使用选项创建日志记录器
    log := logger.New(
        logger.WithLevel("info"),           // 设置日志级别
        logger.WithConsole(true),           // 输出到控制台
        logger.WithJSON(true),              // 使用JSON格式
        logger.WithFilename("./logs/app.log"), // 设置日志文件
        logger.WithMaxSize(10),             // 单个日志文件最大10MB
        logger.WithMaxBackups(5),           // 最多保留5个旧日志文件
        logger.WithMaxAge(30),              // 最多保留30天
        logger.WithCompress(true),          // 压缩旧日志文件
        logger.WithField("service", "api"), // 添加固定字段
    )

    log.Info("服务启动")
}
```

### 结构化日志

```go
package main

import "goadmin/pkg/logger"

func main() {
    log := logger.New()

    // 添加临时字段
    log.Info("用户登录",
        logger.String("username", "admin"),
        logger.String("ip", "192.168.1.1"),
    )

    // 创建带有固定字段的日志记录器
    userLog := log.With(
        logger.String("user_id", "12345"),
        logger.String("username", "admin"),
    )

    userLog.Info("用户操作")
}
```

### 全局日志实例

```go
package main

import "goadmin/pkg/logger"

func main() {
    // 配置全局日志实例
    logger.SetGlobal(logger.New(
        logger.WithLevel("info"),
        logger.WithField("app", "goadmin"),
    ))

    // 在应用的任何地方使用
    log := logger.Global()
    log.Info("使用全局日志实例")
}
```

## 字段辅助函数

日志模块提供了一系列辅助函数来创建字段：

```go
// 字符串字段
logger.String("key", "value")

// 整数字段
logger.Int("count", 42)

// 布尔字段
logger.Bool("enabled", true)

// 浮点数字段
logger.Float64("price", 99.99)

// 错误字段
logger.Error(err)

// 任意类型字段
logger.Any("data", map[string]interface{}{"foo": "bar"})
```

## 完整示例

请运行测试查看完整示例：

```bash
go test -v pkg/logger/logger_test.go
```

## 配置选项

| 选项 | 说明 | 默认值 |
| --- | --- | --- |
| Level | 日志级别 | "info" |
| Console | 是否输出到控制台 | true |
| JSON | 是否使用JSON格式 | false |
| Filename | 日志文件路径 | "./logs/app.log" |
| MaxSize | 单个日志文件最大大小(MB) | 100 |
| MaxBackups | 最大保留的旧日志文件数量 | 3 |
| MaxAge | 保留旧日志文件的最大天数 | 30 |
| Compress | 是否压缩旧日志文件 | true |
| ShowCaller | 是否显示调用者信息 | true |
| CallerSkip | 调用者信息显示层级 | 1 |
| TimeFormat | 时间格式 | time.RFC3339 |

## 注意事项

- 在程序退出前，建议调用 `log.Sync()` 确保所有日志都被写入
- 对于高性能场景，避免使用格式化日志方法（如 `Infof`），而是使用字段方式（如 `Info` 搭配 `logger.String` 等）
- 在生产环境中，建议将日志级别设置为 "info" 或更高级别
