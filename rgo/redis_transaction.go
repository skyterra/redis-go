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
//		reply := redisClient.Transaction(func(tc rgo.TransConn) {
//			redisClient.Set("trans1", "this is trans1", tc)
//			redisClient.Set("trans2", "this is trans2", tc)
//			redisClient.Set("trans3", "this is trans3", tc)
//		})
func (c *RedisClient) Transaction(operation func(conn TransConn)) []interface{} {
	conn := c.pool.Get()
	conn.Do("multi")
	operation(conn)
	reply, _ := redis.Values(conn.Do("exec"))
	return reply
}
