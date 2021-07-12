package rgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"redis-go/rgo"
)

var _ = Describe("RedisString", func() {
	Context("SetWithOptions", func() {
		It("should be succeed", func() {
			redisClient, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
			Expect(err).Should(Succeed())

			options := rgo.NewSetOptions()
			options.ExpireTime = 60
			options.KeepTTL = true
			options.IfExist = true
			options.GetOld = true
			redisClient.SetWithOptions("hello", "world", options)
		})
	})

	Context("Get", func() {
		It("should be succeed", func() {
			redisClient, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
			Expect(err).Should(Succeed())

			Expect(redisClient.Set("test", "this is test")).Should(BeTrue())
			val, exist := redisClient.Get("test")
			Expect(exist).Should(BeTrue())
			Expect(val == "this is test").Should(BeTrue())

			val, exist = redisClient.Get("test1")
			Expect(exist).Should(BeFalse())
			Expect(val == "").Should(BeTrue())
		})
	})
})
