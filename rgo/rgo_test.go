package rgo_test

import (
	"redis-go/rgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stredis", func() {
	Context("SetWithOptions", func() {
		It("should be succeed", func() {
			conn, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
			Expect(err).Should(Succeed())

			options := rgo.NewSetOptions()
			options.ExpireTime = 60
			options.KeepTTL = true
			options.IfExist = true
			options.GetOld = true
			conn.SetWithOptions("hello", "world", options)
		})
	})

	Context("Get", func() {
		It("should be succeed", func() {
			conn, err := rgo.DialRedis("127.0.0.1", "", 0, 0)
			Expect(err).Should(Succeed())

			Expect(conn.Set("test", "this is test")).Should(BeTrue())
			val, exist := conn.Get("test")
			Expect(exist).Should(BeTrue())
			Expect(val == "this is test").Should(BeTrue())

			val, exist = conn.Get("test1")
			Expect(exist).Should(BeFalse())
			Expect(val == "").Should(BeTrue())
		})
	})
})
