package rgo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultAcquireTimeout uint64 = 100       // 获取锁的timeout时间，默认100ms
	DefaultLockTimeout    uint64 = 10 * 1000 // 锁过期时间，默认10s

	RetryInterval time.Duration = 1        // 获取锁失败后，重试间隔，默认1ms
	LockPrefix                  = "stlock" // 锁前缀
)

// 获取分布式锁
// usage:
//		// 并行更新文档读取数量，初始时，文档读取数量为1
// 		docID := "doc:xxx:123:xxx"
//		redisClient.Set(docID, 1)
//
//		// 请求对当前文档添加分布式锁
//		lockID, ok := redisClient.AcquireLock(docID, 100, 10*1000)
//		if !ok {
//			//获取锁失败
//			return
//		}
//
//		// 变更文档读取数量
//		count, ok := redisClient.Get(docID)
//		number, _ := strconv.Atoi(count)
//		redisClient.Set(docID, number+1)
//
//		// 释放当前文档的锁
//		redisClient.ReleaseLock(docID, lockID)
func (c *RedisClient) AcquireLock(resName string, acquireTimeout, lockTimeout uint64) (string, bool) {
	if acquireTimeout <= 0 {
		acquireTimeout = DefaultAcquireTimeout
	}

	if lockTimeout <= 0 {
		lockTimeout = DefaultLockTimeout
	}

	identifier := uuid.NewString()
	endTime := NowMs() + acquireTimeout
	resName = fmt.Sprintf("%s:%s", LockPrefix, resName)

	for NowMs() < endTime {
		options := NewSetOptions()
		options.IfNotExist = true
		options.ExpireTimeInMs = lockTimeout

		ok := c.SetWithOptions(resName, identifier, options)
		if ok {
			return identifier, true
		}

		// 获取锁失败，1ms后重试
		time.Sleep(RetryInterval * time.Millisecond)
	}

	return "", false
}

// 释放分布式锁
func (c *RedisClient) ReleaseLock(resName, identifier string) error {
	resName = fmt.Sprintf("%s:%s", LockPrefix, resName)

	// 检查是否是当前进程添加的锁
	reply, _ := c.Get(resName)
	if reply != identifier {
		return nil
	}

	// 删除锁
	_, err := c.TransactionWithWatch(resName, func(tc TransConn) {
		tc.Do("del", resName)
	})

	return err
}

// 添加分布式锁执行redis命令
// usage:
// redisClient.WithLock(docID, func() {
// 		count, _ := redisClient.Get(docID)
//		number, _ := strconv.Atoi(count)
//		redisClient.Set(docID, number+1)
// })
func (c *RedisClient) WithLock(resName string, doCommand func()) bool {
	lockID, ok := c.AcquireLock(resName, 0, 0)
	if !ok {
		return false
	}

	doCommand()

	c.ReleaseLock(resName, lockID)
	return true
}
