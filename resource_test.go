package goniplus_test

import (
	. "github.com/goniapm/goniplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource", func() {
	BeforeEach(func() {
		InitSDK("APIKEY", 60)
	})
	Describe("GetResource", func() {
		Context("If initial collect", func() {
			It("cpu data should not be included", func() {
				resourceMetric := GetResource()
				_, ok := resourceMetric["cpu"]
				Expect(ok).To(Equal(false))
			})
		})
		Context("If non-initial collect", func() {
			It("cpu data should be included", func() {
				resourceMetric := GetResource()
				_, ok := resourceMetric["cpu"]
				Expect(ok).To(Equal(false))
				resourceMetric = GetResource()
				_, ok = resourceMetric["cpu"]
				Expect(ok).To(Equal(true))
			})
		})
	})
})
