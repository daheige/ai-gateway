package config

import "time"

// Config 项目配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Encrypt  EncryptConfig
	Admin    AdminConfig
}

// ServerConfig 服务端口
type ServerConfig struct {
	Port string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN string
}

// RedisConfig redis配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// JWTConfig jwt 认证key
type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

// EncryptConfig 加密key配置
type EncryptConfig struct {
	Key string
}

// AdminConfig 管理员配置
type AdminConfig struct {
	Username string
	Password string
}

// Load 加载配置
// todo 这里应该哦组配置文件读取
func Load() *Config {
	var jwtKey = "2134bd00292111f1a5d3ff57c47c1333"
	var encKey = "12345678901234567890123456789012"
	return &Config{
		Server:   ServerConfig{Port: ":8080"},
		Database: DatabaseConfig{DSN: "root:root123456@tcp(127.0.0.1:3306)/ai_gateway?charset=utf8mb4&parseTime=True"},
		Redis:    RedisConfig{Addr: "localhost:6379", Password: "", DB: 0},
		JWT:      JWTConfig{Secret: jwtKey, ExpireTime: 24 * time.Hour},
		Encrypt:  EncryptConfig{Key: encKey},
		Admin:    AdminConfig{Username: "admin", Password: "admin123456"},
	}
}
