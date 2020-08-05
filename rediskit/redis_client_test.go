package rediskit

import (
	"github.com/rexue2019/micro/cfg"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	cli := NewRedisClient(&cfg.RedisConfig{
		Host:     "r-wz9wz98e3hd0bx1542.redis.rds.aliyuncs.com",
		Port:     6379,
		Password: "",
	})
	t.Log(cli)
}
