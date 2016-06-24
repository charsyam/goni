package goniplus

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/goniapm/goniplus-worker/metric"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Metric", func() {
	var (
		client *Client
	)
	BeforeEach(func() {
		client = InitSDK("APIKEY", 60)
	})
	Describe("GetTimestamp", func() {
		It("should return different string after time passed", func() {
			unixStamp := GetUnixTimestamp()
			time.Sleep(2 * time.Second)
			afterUnixStamp := GetUnixTimestamp()
			Expect(unixStamp).NotTo(Equal(afterUnixStamp))
		})
	})
	Describe("GetApplicationMetric", func() {
		It("should clear application metric map", func() {
			GetApplicationMetric()
			Expect(len(GetErrorMetric())).To(Equal(0))
			transactionMetric, _, userMetric := GetTransactionMetric()
			Expect(len(transactionMetric.Detail)).To(Equal(0))
			Expect(len(userMetric)).To(Equal(0))
		})
	})
	Describe("GetSystemMetric", func() {
		Context("If initial collect", func() {
			It("expvar metric should be collected", func() {
				systemMetric := GetSystemMetric()
				Expect(len(systemMetric.Expvar) > 0).To(Equal(true))
			})
			It("resource/cpu metric should be 0.0", func() {
				systemMetric := GetSystemMetric()
				Expect(systemMetric.Resource.Cpu).To(Equal(0.0))
			})
		})
		Context("If non-initial collect", func() {
			It("resource/cpu metric should be collected", func() {
				systemMetric := GetSystemMetric()
				time.Sleep(time.Second)
				systemMetric = GetSystemMetric()
				Expect(systemMetric.Resource.Cpu >= 0).To(Equal(true))
			})
		})
	})
	Describe("GetMetric", func() {
		It("should be unmarshalled", func() {
			metric, err := client.getMetric(true)
			if err != nil {
				Fail("Failed to collect metric")
			}
			marshalled := &pb.Metric{}
			err = proto.Unmarshal(metric, marshalled)
			Expect(err).To(BeNil())
		})
	})
})
