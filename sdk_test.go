package goniplus

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Sdk", func() {
	BeforeEach(func() {
		InitSDK("APIKEY", 60)
	})
	Describe("Get identifiers", func() {
		It("should return IP", func() {
			ip := getIP()
			Expect(len(ip) > 0).To(Equal(true))
		})
		It("should return Pid", func() {
			pid := getPid()
			Expect(len(pid) > 0).To(Equal(true))
		})
	})
	Describe("Create instanceID", func() {
		It("should formatted as `$ip-$pid`", func() {
			ip := getIP()
			pid := getPid()
			iID := getInstanceID()
			Expect(iID).To(Equal(ip + "-" + pid))
		})
	})
	Describe("Init SDK", func() {
		It("APIKey must be set", func() {
			Expect(client.apikey).To(Equal("APIKEY"))
		})
		It("Send interval unit should be second", func() {
			Expect(int(client.sendInterval / time.Second)).To(Equal(60))
		})
	})
	Describe("Set option", func() {
		It("setMetricURL should change metric url", func() {
			url := "testUrl"
			client.SetMetricURL(url)
			Expect(metricURL).To(Equal(url))
		})
	})
})
