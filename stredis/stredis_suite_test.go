package stredis_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStredis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stredis Suite")
}
