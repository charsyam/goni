package goniplus_test

import (
	"encoding/json"
	. "github.com/goniapm/goniplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Expvar", func() {
	BeforeEach(func() {
		InitSDK("APIKEY", 60)
	})
	Describe("Collecting metric", func() {
		It("should return more than 1 metric", func() {
			expvarMetric := GetExpvar()
			Expect(len(expvarMetric) > 0).To(Equal(true))
		})
		It("memstats should not contain key in ExcludeExpvarKey", func() {
			expvarMetric := GetExpvar()
			memstats := expvarMetric["memstats"]
			decoder := json.NewDecoder(strings.NewReader(memstats))
			data := make(map[string]interface{})
			if err := decoder.Decode(&data); err != nil {
				Fail("Cannot parse memstats")
			}
			excludeKey := strings.Fields(ExcludeMemstatKey)
			for memstatKey := range data {
				for _, exKey := range excludeKey {
					Expect(memstatKey).NotTo(Equal(exKey))
				}
			}
		})
	})
})
