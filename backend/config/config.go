package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Security SecurityConfig `mapstructure:"security"`
	File     FileConfig     `mapstructure:"file"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode        string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

func (d *DatabaseConfig) DSN() string {
	switch d.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			d.Host, d.Port, d.Username, d.Password, d.DBName, d.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Port, d.DBName)
	default:
		return ""
	}
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	SecretKey    string `mapstructure:"secret_key"`
	ExpirationMs int64  `mapstructure:"expiration_ms"`
	TokenPrefix  string `mapstructure:"token_prefix"`
}

type SecurityConfig struct {
	LoginFailLimit           int `mapstructure:"login_fail_limit"`
	LoginFailResetMinutes    int `mapstructure:"login_fail_reset_minutes"`
	AccountLockDurationMinutes int `mapstructure:"account_lock_duration_minutes"`
}

type FileConfig struct {
	UploadDir string `mapstructure:"upload_dir"`
}

type CORSConfig struct {
	AllowedOrigins string `mapstructure:"allowed_origins"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

var AppConfig *Config

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	AppConfig = &cfg
	return &cfg, nil
}

func InitLogger(cfg *LoggingConfig) *zap.Logger {
	level := zapcore.InfoLevel
	if cfg != nil && cfg.Level != "" {
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			level = zapcore.InfoLevel
		}
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var cores []zapcore.Core

	// Console core
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))

	// File core (if configured)
	if cfg != nil && cfg.File != "" {
		fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
		fileWriter, _, err := zap.Open(cfg.File)
		if err == nil {
			cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, level))
		}
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller())
}
