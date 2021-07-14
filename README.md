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