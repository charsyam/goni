package goniplus_test

import (
	. "github.com/goniapm/goniplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runtime", func() {
	BeforeEach(func() {
		InitSDK("APIKEY", 60)
	})
	Describe("Collecting metric", func() {
		It("should contain Cgo", func() {
			runtimeMetric := GetRuntime()
			Expect(runtimeMetric.Cgo).ToNot(BeNil())
		})
		It("should contain Goroutine", func() {
			runtimeMetric := GetRuntime()
			Expect(runtimeMetric.Goroutine).ToNot(BeNil())
		})
	})
})
