package logger

import (
	"os"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 日志配置
type Config struct {
	// 日志级别 ("debug", "info", "warn", "error", "dpanic", "panic", "fatal")
	Level string `json:"level" yaml:"level"`
	// 是否输出到控制台
	Console bool `json:"console" yaml:"console"`
	// 是否使用JSON格式
	JSON bool `json:"json" yaml:"json"`
	// 日志文件路径
	Filename string `json:"filename" yaml:"filename"`
	// 单个日志文件最大大小，单位MB
	MaxSize int `json:"max_size" yaml:"max_size"`
	// 最大保留的旧日志文件数量
	MaxBackups int `json:"max_backups" yaml:"max_backups"`
	// 保留旧日志文件的最大天数
	MaxAge int `json:"max_age" yaml:"max_age"`
	// 是否压缩旧日志文件
	Compress bool `json:"compress" yaml:"compress"`
	// 调用者信息显示层级
	CallerSkip int `json:"caller_skip" yaml:"caller_skip"`
	// 是否显示调用者信息
	ShowCaller bool `json:"show_caller" yaml:"show_caller"`
	// 时间格式
	TimeFormat string `json:"time_format" yaml:"time_format"`
}

// DefaultConfig 返回默认的日志配置
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Console:    true,
		JSON:       false,
		Filename:   "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
		CallerSkip: 1,
		ShowCaller: true,
		TimeFormat: time.RFC3339,
	}
}

// GetZapLevel 根据字符串获取zap日志级别
func (c *Config) GetZapLevel() zapcore.Level {
	switch c.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// BuildZapCore 构建zap核心
func (c *Config) BuildZapCore() zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { enc.AppendString(t.Format(c.TimeFormat)) }),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志输出位置
	var writeSyncer zapcore.WriteSyncer
	if c.Filename != "" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   c.Filename,
			MaxSize:    c.MaxSize,
			MaxBackups: c.MaxBackups,
			MaxAge:     c.MaxAge,
			Compress:   c.Compress,
		}
		writeSyncer = zapcore.AddSync(lumberJackLogger)
	}

	// 如果需要同时输出到控制台
	if c.Console {
		if writeSyncer != nil {
			writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(zapcore.Lock(os.Stdout)))
		} else {
			writeSyncer = zapcore.AddSync(zapcore.Lock(os.Stdout))
		}
	}

	// 设置编码器
	var encoder zapcore.Encoder
	if c.JSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建Core
	core := zapcore.NewCore(encoder, writeSyncer, c.GetZapLevel())
	return core
}
