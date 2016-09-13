package goniplus_test

import (
	. "github.com/goniapm/goniplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/url"
	"time"
)

var _ = Describe("Http", func() {
	var (
		req *http.Request
	)
	BeforeEach(func() {
		InitSDK("APIKEY", 60)
		req = &http.Request{
			Header:     make(http.Header),
			Method:     "GET",
			RemoteAddr: "127.0.0.1",
			URL: &url.URL{
				Path: "PATH",
			},
		}
		req.Header.Set("User-Agent", "Chrome 41.0.2228.0")
	})
	Describe("Create Transaction ID", func() {
		It("should return different Transaction id", func() {
			method := "GET"
			path := "PATH"
			id1 := CreateRequestID(method, path)
			id2 := CreateRequestID(method, path)
			Expect(id1).NotTo(Equal(id2))
		})
	})
	Describe("Transaction Track (http)", func() {
		It("http header should have a tracking id", func() {
			StartRequestTrack(req)
			Expect(len(req.Header.Get("Goni-tracking-id")) > 0).To(Equal(true))
		})
		It("LeaveBreadcrumb should leave breadcrumb", func() {
			reqTrack := StartRequestTrack(req)
			time.Sleep(time.Second)
			LeaveBreadcrumb(req, "tag")
			time.Sleep(time.Second)
			reqTrack.FinishRequestTrack(200, false)
			transactionMetric, _, _ := GetTransactionMetric()
			Expect(len(transactionMetric.Detail[0].Data[0].Breadcrumb.Tag) > 0).To(Equal(true))
			Expect(len(transactionMetric.Detail[0].Data[0].Breadcrumb.TagT) > 0).To(Equal(true))
		})
		It("FinishRequestTrack should add data to map", func() {
			reqTrack := StartRequestTrack(req)
			time.Sleep(time.Second)
			reqTrack.FinishRequestTrack(200, false)
			transactionMetric, _, _ := GetTransactionMetric()
			Expect(len(transactionMetric.Detail) > 0).To(Equal(true))
		})
	})
})
