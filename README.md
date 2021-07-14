# redis-go
基于[Redigo](https://github.com/gomodule/redigo)库，提供更为直接的redis访问接口

## redis 事务
普通事务
```go
redisClient, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
reply, err := redisClient.Transaction(func(tc rgo.TransConn) {
   	// 以下命令要么全部执行，要么全部不执行
	redisClient.Set("trans1", "this is trans1", tc)
	redisClient.Set("trans2", "this is trans2", tc)
	redisClient.Set("trans3", "this is trans3", tc)
})
```

带有watch的事务
```go
redisClient, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
reply := redisClient.Transaction("WatchKey", func(tc rgo.TransConn) {
	// 如果"WatchKey"在事务执行过程中，有变化，则以下事务执行失败
	redisClient.Set("trans1", "this is trans1", tc)
	redisClient.Set("trans2", "this is trans2", tc)
	redisClient.Set("trans3", "this is trans3", tc)
})
```

## 分布式锁
```go
// 并行更新文档读取数量，初始时，文档读取数量为1
docID := "doc:xxx:123:xxx"
redisClient.Set(docID, 1)

// 请求对当前文档添加分布式锁
lockID, ok := redisClient.AcquireLock(docID, 100, 10*1000)
if !ok {
	//获取锁失败
	return
}

// 变更文档读取数量
count, ok := redisClient.Get(docID)
number, _ := strconv.Atoi(count)
redisClient.Set(docID, number+1)

// 释放当前文档的锁
redisClient.ReleaseLock(docID, lockID)
```