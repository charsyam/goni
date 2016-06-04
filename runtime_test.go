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
		It("should return more than 1 metric", func() {
			runtimeMetric := GetRuntime()
			Expect(len(runtimeMetric) > 0).To(Equal(true))
		})
	})
})
