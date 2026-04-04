package config

import (
	"log"
	"time"

	"github.com/go-god/setting"
)

// Config 项目配置结构体
type Config struct {
	AppEnv       string        `mapstructure:"app_env"`
	AppDebug     string        `mapstructure:"app_debug"`
	AppPort      int           `mapstructure:"app_port"`
	GracefulWait time.Duration `mapstructure:"graceful_wait"`

	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Encrypt  EncryptConfig  `mapstructure:"encrypt"`
	Admin    AdminConfig    `mapstructure:"admin"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN     string `mapstructure:"dsn"`
	ShowSQL bool   `mapstructure:"show_sql"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig jwt 认证key
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
}

// EncryptConfig 加密key配置
type EncryptConfig struct {
	Key string `mapstructure:"key"`
}

// AdminConfig 管理员配置
type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Load 加载配置
func Load() *Config {
	conf := setting.New(setting.WithConfigFile("./app.yaml"))
	err := conf.Load()
	if err != nil {
		log.Fatalln("load config err: ", err)
	}

	c := &Config{}
	err = conf.ReadSection("app_conf", c)
	if err != nil {
		log.Fatalln("read config err: ", err)
	}

	return c
}
