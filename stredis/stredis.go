package stredis

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// redis连接池
type ConnPool struct {
	pool *redis.Pool
}

// redis set命令，指定附加项
type SetOptions struct {
	ExpireTime       uint64 // EX，过期时间 (秒）
	ExpireTimeInMs   uint64 // PX，过期时间（毫秒）
	ExpireTimeAT     uint64 // EXAT，过期时间戳（精确到秒）
	ExpireTImeATInMs uint64 // PXAT，过期时间戳（精确到毫秒）
	KeepTTL          bool   // KEEPTTL，保留设置前指定键的生存时间

	IfNotExist bool // NX，仅当不存在key时，设置key
	IfExist    bool // XX，仅当存在key时，设置key

	GetOld bool // GET，返回key设置前存储的值
}

func (c *ConnPool) Do(command string, args ...interface{}) (interface{}, error) {
	conn := c.pool.Get()
	return conn.Do(command, args...)
}

// redis set command
func (c *ConnPool) Set(key string, value string) bool {
	conn := c.pool.Get()
	reply, _ := redis.String(conn.Do("set", key, value))
	return reply == "OK"
}

// redis set command with options
func (c *ConnPool) SetWithOptions(key string, value string, options *SetOptions) bool {
	conn := c.pool.Get()
	args := []interface{}{key, value}
	args = append(args, options.ToArgs()...)
	reply, _ := redis.String(conn.Do("set", args...))
	return reply == "OK"
}

// 关闭连接池
func (c *ConnPool) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

func (opt *SetOptions) ToArgs() []interface{} {
	var args []interface{}

	switch {
	case opt.ExpireTime > 0:
		args = append(args, "ex")
		args = append(args, opt.ExpireTime)
	case opt.ExpireTimeInMs > 0:
		args = append(args, "px")
		args = append(args, opt.ExpireTimeInMs)
	case opt.ExpireTimeAT > 0:
		args = append(args, "exat")
		args = append(args, opt.ExpireTimeAT)
	case opt.ExpireTImeATInMs > 0:
		args = append(args, "pxat")
		args = append(args, opt.ExpireTImeATInMs)
	case opt.KeepTTL:
		args = append(args, "keepttl")
	}

	switch {
	case opt.IfNotExist:
		args = append(args, "nx")
	case opt.IfExist:
		args = append(args, "xx")
	}

	if opt.GetOld {
		args = append(args, "get")
	}

	return args
}

// 连接redis，返回连接池
func DialRedis(host, password string, port, db int) (*ConnPool, error) {
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

	return &ConnPool{pool: pool}, nil
}

func NewSetOptions() *SetOptions {
	return &SetOptions{}
}
