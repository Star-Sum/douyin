package RedisDao

import "context"

// FlushDb 在服务器结束后删除redis内存数据
func FlushDb() {
	redisHandler.FlushDB(context.Background())
}
