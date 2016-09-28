package goniplus_test

import (
	. "github.com/monitflux/goniplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource", func() {
	BeforeEach(func() {
		InitSDK("APIKEY", 60)
	})
	Describe("GetResource", func() {
		Context("If initial collect", func() {
			It("cpu data should be 0.0", func() {
				resourceMetric := GetResource()
				Expect(resourceMetric.Cpu).To(Equal(0.0))
			})
		})
		Context("If non-initial collect", func() {
			It("cpu data should be included", func() {
				resourceMetric := GetResource()
				Expect(resourceMetric.Cpu).To(Equal(0.0))
				resourceMetric = GetResource()
				Expect(resourceMetric.Cpu >= 0).To(Equal(true))
			})
		})
	})
})
