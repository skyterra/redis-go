package stredis_test

import (
	"redis-go/stredis"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stredis", func() {
	Context("SetWithOptions", func() {
		It("should be succeed", func() {
			conn, err := stredis.DialRedis("127.0.0.1", "", 0, 0)
			Expect(err).Should(Succeed())

			options := stredis.NewSetOptions()
			options.ExpireTime = 60
			options.KeepTTL = true
			options.IfExist = true
			options.GetOld = true
			conn.SetWithOptions("hello", "world", options)
		})
	})
})
