package logger

// Option 日志选项函数类型
type Option func(*zapLogger)

// WithConfig 设置日志配置
func WithConfig(config *Config) Option {
	return func(l *zapLogger) {
		l.config = config
	}
}

// WithLevel 设置日志级别
func WithLevel(level string) Option {
	return func(l *zapLogger) {
		l.config.Level = level
	}
}

// WithConsole 设置是否输出到控制台
func WithConsole(console bool) Option {
	return func(l *zapLogger) {
		l.config.Console = console
	}
}

// WithJSON 设置是否使用JSON格式
func WithJSON(json bool) Option {
	return func(l *zapLogger) {
		l.config.JSON = json
	}
}

// WithFilename 设置日志文件路径
func WithFilename(filename string) Option {
	return func(l *zapLogger) {
		l.config.Filename = filename
	}
}

// WithMaxSize 设置单个日志文件最大大小
func WithMaxSize(maxSize int) Option {
	return func(l *zapLogger) {
		l.config.MaxSize = maxSize
	}
}

// WithMaxBackups 设置最大保留的旧日志文件数量
func WithMaxBackups(maxBackups int) Option {
	return func(l *zapLogger) {
		l.config.MaxBackups = maxBackups
	}
}

// WithMaxAge 设置保留旧日志文件的最大天数
func WithMaxAge(maxAge int) Option {
	return func(l *zapLogger) {
		l.config.MaxAge = maxAge
	}
}

// WithCompress 设置是否压缩旧日志文件
func WithCompress(compress bool) Option {
	return func(l *zapLogger) {
		l.config.Compress = compress
	}
}

// WithCallerSkip 设置调用者信息显示层级
func WithCallerSkip(callerSkip int) Option {
	return func(l *zapLogger) {
		l.config.CallerSkip = callerSkip
	}
}

// WithShowCaller 设置是否显示调用者信息
func WithShowCaller(showCaller bool) Option {
	return func(l *zapLogger) {
		l.config.ShowCaller = showCaller
	}
}

// WithTimeFormat 设置时间格式
func WithTimeFormat(timeFormat string) Option {
	return func(l *zapLogger) {
		l.config.TimeFormat = timeFormat
	}
}

// WithField 添加单个固定字段
func WithField(key string, value interface{}) Option {
	return func(l *zapLogger) {
		l.fields = append(l.fields, toZapFields([]Field{{Key: key, Value: value}})...)
	}
}
