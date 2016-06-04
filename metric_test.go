package goniplus_test

import (
	. "github.com/goniapm/goniplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Metric", func() {
	BeforeEach(func() {
		InitSDK("APIKEY", 60)
	})
	Describe("GetTimestamp", func() {
		It("should return different string after time passed", func() {
			unixStamp := GetUnixTimestamp()
			time.Sleep(2 * time.Second)
			afterUnixStamp := GetUnixTimestamp()
			Expect(unixStamp).NotTo(Equal(afterUnixStamp))
			stamp := GetTimestamp()
			time.Sleep(2 * time.Second)
			afterStamp := GetTimestamp()
			Expect(stamp).NotTo(Equal(afterStamp))
		})
	})
	Describe("GetApplicationMetric", func() {
		It("should clear application metric map", func() {
			GetApplicationMetric()
			Expect(len(GetErrorMetric())).To(Equal(0))
			httpMetric, userMetric := GetHTTPResponseMetric()
			Expect(len(httpMetric)).To(Equal(0))
			Expect(len(userMetric)).To(Equal(0))
		})
	})
	Describe("GetSystemCollect", func() {
		Context("If initial collect", func() {
			It("expvar metric should be collected", func() {
				systemMetric := GetSystemMetric()
				Expect(len(systemMetric.Expvar) > 0).To(Equal(true))
			})
			It("resource/cpu metric should not be collected", func() {
				systemMetric := GetSystemMetric()
				_, ok := systemMetric.Resource["cpu"]
				Expect(ok).To(Equal(false))
			})
			It("runtime metric should be collected", func() {
				systemMetric := GetSystemMetric()
				Expect(len(systemMetric.Runtime) > 0).To(Equal(true))
			})
		})
		Context("If non-initial collect", func() {
			It("resource/cpu metric should be collected", func() {
				systemMetric := GetSystemMetric()
				time.Sleep(time.Second)
				systemMetric = GetSystemMetric()
				_, ok := systemMetric.Resource["cpu"]
				Expect(ok).To(Equal(true))
			})
		})
	})
})
