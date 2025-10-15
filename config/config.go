package config

import (
	"fmt"
	"goadmin/pkg/logger"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Logger   logger.Config  `yaml:"logger"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Debug   bool   `yaml:"debug"`
	Port    int    `yaml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Enable bool       `yaml:"enable"`
	Master DBConfig   `yaml:"master"`
	Slaves []DBConfig `yaml:"slaves"`
}

// DBConfig 数据库连接配置
type DBConfig struct {
	Driver          string        `yaml:"driver"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	Charset         string        `yaml:"charset"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	LogLevel        string        `yaml:"log_level"`
}

func (dbCfg *DBConfig) DSN() string {
	switch strings.ToLower(dbCfg.Driver) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			dbCfg.Username,
			dbCfg.Password,
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.Database,
			dbCfg.Charset,
		)
	case "postgres", "postgresql":
		return fmt.Sprintf(`postgresql://%s:%s@%s:%d/%s?sslmode=disable`,
			dbCfg.Username,
			dbCfg.Password,
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.Database)
	default:
		return ""
	}
}

// RedisConfig Redis配置
type RedisConfig struct {
	Enable       bool          `yaml:"enable"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolTimeout  time.Duration `yaml:"pool_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	MaxConnAge   time.Duration `yaml:"max_conn_age"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string        `yaml:"secret"`
	AccessExpire    time.Duration `yaml:"access_expire"`
	RefreshExpire   time.Duration `yaml:"refresh_expire"`
	Issuer          string        `yaml:"issuer"`
	RefreshTokenKey string        `yaml:"refresh_token_key"`
}

var (
	cfg Config
)

func Get() *Config {
	return &cfg
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("获取配置文件绝对路径失败: %w", err)
	}

	// 读取配置文件
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置文件
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 处理时间字段
	if err := parseTimeDurations(&cfg); err != nil {
		return nil, fmt.Errorf("解析时间字段失败: %w", err)
	}

	return &cfg, nil
}

// parseTimeDurations 解析配置中的时间字段
func parseTimeDurations(cfg *Config) error {
	// 数据库连接最大生命周期
	if err := parseDuration(&cfg.Database.Master.ConnMaxLifetime); err != nil {
		return err
	}

	// 从库连接最大生命周期
	for i := range cfg.Database.Slaves {
		if err := parseDuration(&cfg.Database.Slaves[i].ConnMaxLifetime); err != nil {
			return err
		}
	}

	// Redis相关超时设置
	if err := parseDuration(&cfg.Redis.DialTimeout); err != nil {
		return err
	}
	if err := parseDuration(&cfg.Redis.ReadTimeout); err != nil {
		return err
	}
	if err := parseDuration(&cfg.Redis.WriteTimeout); err != nil {
		return err
	}
	if err := parseDuration(&cfg.Redis.PoolTimeout); err != nil {
		return err
	}
	if err := parseDuration(&cfg.Redis.IdleTimeout); err != nil {
		return err
	}
	if err := parseDuration(&cfg.Redis.MaxConnAge); err != nil {
		return err
	}

	// JWT相关时间设置
	if err := parseDuration(&cfg.JWT.AccessExpire); err != nil {
		return err
	}
	if err := parseDuration(&cfg.JWT.RefreshExpire); err != nil {
		return err
	}

	return nil
}

// parseDuration 将字符串时间值转换为time.Duration
func parseDuration(d *time.Duration) error {
	if *d != 0 {
		return nil // 已经是time.Duration类型
	}

	// 获取原始字符串值
	dStr := fmt.Sprintf("%v", *d)
	if dStr == "0s" || dStr == "" {
		return nil // 默认值或空值
	}

	// 解析时间字符串
	parsed, err := time.ParseDuration(dStr)
	if err != nil {
		return fmt.Errorf("解析时间字段失败: %v -> %w", dStr, err)
	}

	*d = parsed
	return nil
}
