package goniplus

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Sender", func() {
	BeforeEach(func() {
		InitSDK("APIKEY", 5)
	})
	Describe("send metric to worker", func() {
		It("should return false when url is not exists", func() {
			url := "testUrl"
			client.SetMetricURL(url)
			result := sendMetricData([]byte{})
			Expect(result).To(Equal(false))
		})
		It("should return true when worker ends connection", func() {
			url := "goni.goniapm.io:9900"
			client.SetMetricURL(url)
			result := sendMetricData([]byte{})
			Expect(result).To(Equal(true))
		})
	})
	Describe("startSender", func() {
		It("should not panic when sending metric", func() {
			url := "testUrl"
			client.SetMetricURL(url)
			time.Sleep(30 * time.Second)
		})
	})
})
