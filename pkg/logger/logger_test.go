package logger

import (
	"testing"
)

// TestBasicUsage 测试基本日志功能
func TestBasicUsage(t *testing.T) {
	// 创建一个默认配置的日志记录器
	log := New()

	// 记录不同级别的日志
	log.Debug("这是一条调试日志")
	log.Info("这是一条信息日志")
	log.Warn("这是一条警告日志")
	log.Error("这是一条错误日志")

	// 使用格式化日志
	log.Debugf("格式化的调试日志: %s", "调试信息")
	log.Infof("用户 %s 登录系统，IP: %s", "admin", "127.0.0.1")
	log.Warnf("服务器负载达到 %.2f%%", 85.5)
	log.Errorf("数据库连接失败: %v", "连接超时")
}

// TestCustomConfig 测试自定义配置
func TestCustomConfig(t *testing.T) {
	// 自定义配置
	config := &Config{
		Level:      "info",           // 只输出info及以上级别的日志
		Console:    true,             // 输出到控制台
		JSON:       true,             // 使用JSON格式
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    10,               // 单个日志文件最大10MB
		MaxBackups: 5,                // 最多保留5个旧日志文件
		MaxAge:     30,               // 最多保留30天
		Compress:   true,             // 压缩旧日志文件
		ShowCaller: true,             // 显示调用者信息
		CallerSkip: 1,                // 调用者信息显示层级
	}

	// 使用自定义配置创建日志记录器
	log := New(WithConfig(config))

	log.Info("使用自定义配置的日志记录器")
}

// TestWithFields 测试字段注入功能
func TestWithFields(t *testing.T) {
	// 创建一个带有固定字段的日志记录器
	log := New(
		WithField("service", "user-service"),
		WithField("version", "1.0.0"),
	)

	// 记录日志时会自动带上这些固定字段
	log.Info("用户服务启动")

	// 添加临时字段
	log.Info("用户登录",
		String("username", "admin"),
		String("ip", "192.168.1.1"),
	)

	// 创建一个新的派生日志记录器，带有额外的固定字段
	userLog := log.With(
		String("user_id", "12345"),
		String("username", "admin"),
	)

	// 这条日志会包含所有固定字段
	userLog.Info("用户操作日志")
}

// TestGlobalLogger 测试全局日志实例
func TestGlobalLogger(t *testing.T) {
	// 配置全局日志实例
	SetGlobal(New(
		WithLevel("info"),
		WithField("app", "goadmin"),
	))

	// 在应用的任何地方获取并使用全局日志实例
	log := Global()
	log.Info("使用全局日志实例")
}

// TestContextualLogger 测试上下文日志
func TestContextualLogger(t *testing.T) {
	log := New()

	// 假设从context中提取了一些信息
	requestLog := log.With(
		String("request_id", "req-123"),
		String("user_id", "user-456"),
	)

	requestLog.Info("处理请求")
}

// TestLoggerOptions 测试日志选项
func TestLoggerOptions(t *testing.T) {
	// 使用多个选项创建日志记录器
	log := New(
		WithLevel("debug"),
		WithConsole(true),
		WithJSON(false),
		WithFilename("./logs/custom.log"),
		WithMaxSize(50),
		WithMaxBackups(3),
		WithMaxAge(7),
		WithCompress(true),
		WithShowCaller(true),
		WithCallerSkip(1),
		WithField("environment", "production"),
	)

	log.Info("使用多个选项配置的日志记录器")
}

// TestCustomFormat 测试自定义格式
func TestCustomFormat(t *testing.T) {
	config := DefaultConfig()
	config.JSON = false                       // 使用控制台格式而非JSON
	config.TimeFormat = "2006-01-02 15:04:05" // 自定义时间格式

	log := New(WithConfig(config))
	log.Info("使用自定义格式的日志记录器")
}

// TestLoggerSync 测试日志同步
func TestLoggerSync(t *testing.T) {
	log := New(
		WithFilename("./logs/sync_test.log"),
	)

	log.Info("测试日志同步")
	if err := log.Sync(); err != nil {
		t.Errorf("日志同步失败: %v", err)
	}
}

// TestLogLevels 测试所有日志级别
func TestLogLevels(t *testing.T) {
	log := New(WithLevel("debug"))

	log.Debug("调试级别日志")
	log.Info("信息级别日志")
	log.Warn("警告级别日志")
	log.Error("错误级别日志")
	// 以下级别会导致程序退出，仅作示例
	// log.DPanic("开发环境恐慌级别日志")
	// log.Panic("恐慌级别日志")
	// log.Fatal("致命级别日志")
}

// TestFieldHelpers 测试字段辅助函数
func TestFieldHelpers(t *testing.T) {
	log := New()

	// 测试各种字段类型
	log.Info("测试字段辅助函数",
		String("string_field", "字符串"),
		Int("int_field", 42),
		Bool("bool_field", true),
		Float64("float_field", 3.14),
		Any("map_field", map[string]string{"key": "value"}),
	)

	// 测试错误字段
	err := &testError{"测试错误"}
	log.Info("错误日志", Error(err))
}

// 实现一个简单的错误类型用于测试
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
