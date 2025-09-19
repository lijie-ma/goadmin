package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// global 全局日志实例
	global     Logger
	globalOnce sync.Once
)

// zapLogger 基于zap实现的日志记录器
type zapLogger struct {
	config     *Config
	zap        *zap.Logger
	sugar      *zap.SugaredLogger
	fields     []zapcore.Field
	zapOptions []zap.Option
}

// New 创建一个新的日志记录器
func New(opts ...Option) Logger {
	logger := &zapLogger{
		config:     DefaultConfig(),
		fields:     make([]zapcore.Field, 0),
		zapOptions: make([]zap.Option, 0),
	}

	// 应用选项
	for _, opt := range opts {
		opt(logger)
	}

	// 如果需要显示调用者信息
	if logger.config.ShowCaller {
		logger.zapOptions = append(logger.zapOptions, zap.AddCaller(), zap.AddCallerSkip(logger.config.CallerSkip))
	}

	// 构建zap logger
	logger.zap = zap.New(
		logger.config.BuildZapCore(),
		logger.zapOptions...,
	)

	// 创建sugar logger
	logger.sugar = logger.zap.Sugar()

	return logger
}

// Global 获取全局日志实例
func Global() Logger {
	globalOnce.Do(func() {
		global = New()
	})
	return global
}

// SetGlobal 设置全局日志实例
func SetGlobal(logger Logger) {
	global = logger
}

// With 创建一个带有固定字段的新日志记录器
func (l *zapLogger) With(fields ...Field) Logger {
	zapFields := make([]zapcore.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}

	newLogger := &zapLogger{
		config: l.config,
		// fields:     append(l.fields, zapFields...),
		zapOptions: l.zapOptions,
	}
	newLogger.zap = l.zap.With(zapFields...)
	newLogger.sugar = newLogger.zap.Sugar()
	return newLogger
}

// WithContext 从上下文创建一个新的日志记录器
func (l *zapLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}
	return l
}

// Debug 记录调试级别日志
func (l *zapLogger) Debug(msg string, fields ...Field) {
	zapFields := toZapFields(fields)
	l.zap.Debug(msg, append(l.fields, zapFields...)...)
}

// Info 记录信息级别日志
func (l *zapLogger) Info(msg string, fields ...Field) {
	zapFields := toZapFields(fields)
	l.zap.Info(msg, append(l.fields, zapFields...)...)
}

// Warn 记录警告级别日志
func (l *zapLogger) Warn(msg string, fields ...Field) {
	zapFields := toZapFields(fields)
	l.zap.Warn(msg, append(l.fields, zapFields...)...)
}

// Error 记录错误级别日志
func (l *zapLogger) Error(msg string, fields ...Field) {
	zapFields := toZapFields(fields)
	l.zap.Error(msg, append(l.fields, zapFields...)...)
}

// DPanic 记录开发环境恐慌级别日志
func (l *zapLogger) DPanic(msg string, fields ...Field) {
	zapFields := toZapFields(fields)
	l.zap.DPanic(msg, append(l.fields, zapFields...)...)
}

// Panic 记录恐慌级别日志
func (l *zapLogger) Panic(msg string, fields ...Field) {
	zapFields := toZapFields(fields)
	l.zap.Panic(msg, append(l.fields, zapFields...)...)
}

// Fatal 记录致命级别日志
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	zapFields := toZapFields(fields)
	l.zap.Fatal(msg, append(l.fields, zapFields...)...)
}

// Debugf 使用格式化字符串记录调试级别日志
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Infof 使用格式化字符串记录信息级别日志
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Warnf 使用格式化字符串记录警告级别日志
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Errorf 使用格式化字符串记录错误级别日志
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// DPanicf 使用格式化字符串记录开发环境恐慌级别日志
func (l *zapLogger) DPanicf(format string, args ...interface{}) {
	l.sugar.DPanicf(format, args...)
}

// Panicf 使用格式化字符串记录恐慌级别日志
func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugar.Panicf(format, args...)
}

// Fatalf 使用格式化字符串记录致命级别日志
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

// Sync 同步缓冲区
func (l *zapLogger) Sync() error {
	return l.zap.Sync()
}

// 工具函数: 将自定义Field转换为zap.Field
func toZapFields(fields []Field) []zapcore.Field {
	if len(fields) == 0 {
		return []zapcore.Field{}
	}

	zapFields := make([]zapcore.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

// 提供常用字段创建函数
func String(key string, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
