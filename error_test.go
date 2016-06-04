package goniplus_test

import (
	. "github.com/goniapm/goniplus"

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
		It("should have multiple tags", func() {
			Error("tagA", "explanation")
			Error("tagB", "explanation")
			errorMetric := GetErrorMetric()
			Expect(len(errorMetric)).To(Equal(2))
		})
		It("should have multiple values", func() {
			tag := "tag"
			Error(tag, "explanationA")
			Error(tag, "explanationB")
			errorMetric := GetErrorMetric()
			Expect(len(errorMetric)).To(Equal(1))
			Expect(errorMetric[tag][0]).ToNot(Equal(errorMetric[tag][1]))
		})
		It("should have duplicated values", func() {
			tag := "tag"
			explanation := "explanation"
			Error(tag, explanation)
			Error(tag, explanation)
			errorMetric := GetErrorMetric()
			Expect(len(errorMetric)).To(Equal(1))
			Expect(errorMetric[tag][0]).To(Equal(errorMetric[tag][1]))
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
