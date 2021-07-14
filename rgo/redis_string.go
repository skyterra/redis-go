package rgo

import "github.com/gomodule/redigo/redis"

/*
 * Redis String
 */

// SET 命令额外参数
type SetOptions struct {
	ExpireTime       uint64 // EX 		过期时间 (秒）
	ExpireTimeInMs   uint64 // PX 		过期时间（毫秒）
	ExpireTimeAT     uint64 // EXAT 	过期时间戳（精确到秒）
	ExpireTImeATInMs uint64 // PXAT 	过期时间戳（精确到毫秒）
	KeepTTL          bool   // KEEPTTL	保留设置前指定键的生存时间
	IfNotExist       bool   // NX 		仅当不存在key时，设置key
	IfExist          bool   // XX 		仅当存在key时，设置key
	GetOld           bool   // GET 		返回key设置前存储的值
}

// SET命令
func (c *RedisClient) Set(key string, value interface{}, tc ...TransConn) bool {
	conn := c.getConn(tc)
	reply, _ := redis.String(conn.Do("set", key, value))
	return reply == "OK"
}

// SET命令 指定额外参数
// usage:
//		options := rgo.NewSetOptions()
//		options.ExpireTime = 60
//		options.KeepTTL = true
//		options.IfExist = true
//		options.GetOld = true
//		redisClient.SetWithOptions("hello", "world", options)
//
func (c *RedisClient) SetWithOptions(key string, value string, options *SetOptions, tc ...TransConn) bool {
	conn := c.getConn(tc)
	args := []interface{}{key, value}
	args = append(args, options.ToArgs()...)
	reply, _ := redis.String(conn.Do("set", args...))
	return reply == "OK"
}

// GET命令
func (c *RedisClient) Get(key string, tc ...TransConn) (string, bool) {
	conn := c.getConn(tc)
	reply, err := redis.String(conn.Do("get", key))
	if err != nil {
		return "", false
	}

	return reply, true
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

func NewSetOptions() *SetOptions {
	return &SetOptions{}
}
