package rediskit

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rexue2019/micro/cfg"
	"github.com/rexue2019/util/errorkit"
	"github.com/rexue2019/util/logging"
)

type RedisClient struct {
	cfgRedis *cfg.RedisConfig
	*redis.Client
}

func NewRedisClient(config *cfg.RedisConfig) *RedisClient {
	o := &RedisClient{}
	o.cfgRedis = config
	o.initRedisClient()
	return o
}

func (u *RedisClient) Close() {

}

func (u *RedisClient) initRedisClient() {
	if u.Client == nil {
		if u.cfgRedis.Host == "" || u.cfgRedis.Port == 0 {
			logging.Log.FatalError(errorkit.WrapError(nil).SetMessage("invalid redis config"))
		}

		addr := fmt.Sprintf("%s:%d", u.cfgRedis.Host, u.cfgRedis.Port)
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: u.cfgRedis.Password,
			DB:       0,
		})

		//通过withContext，可以自带tracing的能力
		u.Client = client

		logging.Log.Infow(logging.KeyScope, "redis", logging.KeyEvent, "connectSuccess")
		err := u.Client.Ping().Err()
		if err != nil {
			logging.Log.FatalError(errorkit.WrapError(err).AddOp("redis.Ping"))
		}
	}
	return
}
