package rgo_test

import (
	"redis-go/rgo"
	"strconv"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisLock", func() {
	Context("Acquire Lock", func() {
		It("should be succeed", func() {
			redisClient, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
			Expect(err).Should(Succeed())

			// 设置文档默认读取数量为1
			docID := "doc:xxx:123:xxx"
			redisClient.Set(docID, 1)

			total := 10
			wg := sync.WaitGroup{}
			wg.Add(total)

			for i := 0; i < total; i++ {
				go func() {
					id, ok := redisClient.AcquireLock(docID, 100, 10*1000)
					Expect(ok).Should(BeTrue())

					count, ok := redisClient.Get(docID)
					Expect(ok).Should(BeTrue())

					number, _ := strconv.Atoi(count)
					redisClient.Set(docID, number+1)

					redisClient.ReleaseLock(docID, id)
					wg.Done()
				}()
			}

			wg.Wait()

			reply, _ := redisClient.Get(docID)
			number, _ := strconv.Atoi(reply)
			Expect(number == 11).Should(BeTrue())
		})
	})

})
