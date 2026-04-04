package config

import (
	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化redis
func InitRedis(conf RedisConfig) redis.UniversalClient {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{
			conf.Addr,
		},
		Password: conf.Password,
		DB:       conf.DB,
	})

	return rdb
}
