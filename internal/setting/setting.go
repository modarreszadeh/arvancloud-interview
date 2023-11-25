package setting

import "github.com/modarreszadeh/arvancloud-interview/pkg/redis"

var (
	Port  = ":5000"
	Redis = redis.Config{Addr: "redis:6379", Password: "redispass", DB: 0}
)
