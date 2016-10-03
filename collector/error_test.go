package goniplus_test

import (
	. "github.com/monitflux/goniplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	BeforeEach(func() {
		InitSDK("APIKEY", 60)
	})
	Describe("Creating error", func() {
		It("should added to map", func() {
			Error("tag", "explanation")
			errorMetric := GetErrorMetric()
			Expect(len(errorMetric)).To(Equal(1))
		})
	})
	Describe("Collecting error", func() {
		It("should return empty map after get", func() {
			Error("tag", "explanation")
			errorMetric := GetErrorMetric()
			Expect(len(errorMetric)).To(Equal(1))
			errorMetric = GetErrorMetric()
			Expect(len(errorMetric)).To(Equal(0))
		})
	})
})
