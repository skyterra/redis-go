package rgo

import (
	"github.com/gomodule/redigo/redis"
)

/*
 * Redis 事务
 */

type TransConn = redis.Conn

// Redis 事务函数
// usage:
//		redisClient, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
//		reply, err := redisClient.Transaction(func(tc rgo.TransConn) {
//			// 以下命令要么全部执行，要么全部不执行
//			redisClient.Set("trans1", "this is trans1", tc)
//			redisClient.Set("trans2", "this is trans2", tc)
//			redisClient.Set("trans3", "this is trans3", tc)
//		})
func (c *RedisClient) Transaction(doCommands func(conn TransConn)) ([]interface{}, error) {
	conn := c.pool.Get()
	conn.Do("multi")
	doCommands(conn)
	return redis.Values(conn.Do("exec"))
}

// Redis 事务函数
// usage:
//		redisClient, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
//		reply := redisClient.Transaction("WatchKey", func(tc rgo.TransConn) {
//			// 如果"WatchKey"在事务执行过程中，有变化，则以下事务执行失败
//			redisClient.Set("trans1", "this is trans1", tc)
//			redisClient.Set("trans2", "this is trans2", tc)
//			redisClient.Set("trans3", "this is trans3", tc)
//		})
func (c *RedisClient) TransactionWithWatch(watchKey string, doCommands func(conn TransConn)) ([]interface{}, error) {
	conn := c.pool.Get()

	conn.Do("watch", watchKey)
	defer conn.Do("unwatch", watchKey)

	conn.Do("multi")
	doCommands(conn)
	return redis.Values(conn.Do("exec"))
}
