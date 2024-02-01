package provider

import (
	"Technology-Blog/Test/config"
	"Technology-Blog/Test/log"
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	Pool *redis.Pool
}

func (r *Redis) New() *Redis {
	r.Pool = &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				config.Get(config.RedisAddr),
				redis.DialPassword(config.Get(config.RedisPassword)),
			)
		},
	}
	return r
}

func (r *Redis) Start() {

}

func (r *Redis) Close() {
	if err := r.Pool.Close(); err != nil {
		_ = log.Error(err)
	}
}
