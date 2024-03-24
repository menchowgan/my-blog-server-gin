package redis

import (
	"gmc-blog-server/config"
	"time"

	"github.com/go-redis/redis"
)

var (
	RedisDb *redis.Client
)

func Init() error {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		Password: "",
		DB:       0,
	})

	_, err := RedisDb.Ping().Result()

	if err != nil {
		return err
	}

	return nil
}

func GetKey(typeS string, value string) string {
	switch typeS {
	case "user":
		return UserInfo + value
	case "photo":
		return PhotoInfo + value
	case "plan":
		return PlanInfo + value
	default:
		return ""
	}
}

func Test() (string, error) {
	s, err := RedisDb.Set("test-key", "redis-learning", 1000*time.Second).Result()
	if err == redis.Nil || err != nil {
		return "", err
	}
	return s, nil
}
