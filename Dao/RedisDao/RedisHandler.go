package RedisDao

import "github.com/redis/go-redis/v9"

var redisHandler *redis.Client

func GetRedisDBHandler(client *redis.Client) {
	redisHandler = client
}
