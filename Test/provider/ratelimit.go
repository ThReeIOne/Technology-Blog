package provider

import (
	"github.com/gomodule/redigo/redis"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/redigostore"
)

type RateLimiter struct {
	Second *throttled.GCRARateLimiter
	Minute *throttled.GCRARateLimiter
}

func (r *RateLimiter) New(redisClient *redis.Pool) *RateLimiter {
	return &RateLimiter{
		Second: newLimiter(redisClient, 1, throttled.PerSec(1)),
		Minute: newLimiter(redisClient, 1, throttled.PerMin(1)),
	}
}

func newLimiter(redisClient *redis.Pool, maxBurst int, maxRate throttled.Rate) *throttled.GCRARateLimiter {
	store, _ := redigostore.New(redisClient, "limiter:", 0)
	quota := throttled.RateQuota{
		MaxRate:  maxRate,
		MaxBurst: maxBurst,
	}
	rateLimiter, _ := throttled.NewGCRARateLimiter(store, quota)
	return rateLimiter
}

func (r *RateLimiter) Start() {
}

func (r *RateLimiter) Close() {
}
