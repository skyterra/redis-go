package rgo

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

/*
 * rgo: Redis-Go
 */

// Redis客户端
type RedisClient struct {
	pool *redis.Pool
}

// 通用命令接口
func (c *RedisClient) Do(command string, args ...interface{}) (interface{}, error) {
	conn := c.pool.Get()
	return conn.Do(command, args...)
}

func (c *RedisClient) getConn(tc []TransConn) redis.Conn {
	if len(tc) > 0 {
		return tc[0]
	}

	return c.pool.Get()
}

// 连接redis，返回连接池
func DialRedis(host, password string, port, db int) (*RedisClient, error) {
	if port == 0 {
		port = 6379
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, redis.DialDatabase(db), redis.DialPassword(password))
		},
	}

	conn := pool.Get()
	defer conn.Close()

	if conn.Err() != nil {
		return nil, conn.Err()
	}

	if r, _ := redis.String(conn.Do("PING")); r != "PONG" {
		return nil, errors.New("connect redis failed")
	}

	return &RedisClient{pool: pool}, nil
}

// 关闭连接池
func CloseRedis(c *RedisClient) {
	if c.pool != nil {
		c.pool.Close()
	}
}
