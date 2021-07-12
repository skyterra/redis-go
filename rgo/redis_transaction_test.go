package rgo_test

import (
	"redis-go/rgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisTransaction", func() {
	Context("Transaction", func() {
		It("should be succeed", func() {
			redisClient, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
			Expect(err).Should(Succeed())

			reply := redisClient.Transaction(func(tc rgo.TransConn) {
				redisClient.Set("trans1", "this is trans1", tc)
				redisClient.Set("trans2", "this is trans2", tc)
				redisClient.Set("trans3", "this is trans3", tc)
			})

			Expect(len(reply) == 3).Should(BeTrue())
			for i := 0; i < len(reply); i++ {
				Expect(reply[i].(string) == "OK").Should(BeTrue())
			}
		})
	})
})
