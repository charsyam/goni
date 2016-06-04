package goniplus_test

import (
	. "github.com/goniapm/goniplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CPU", func() {
	Describe("GetCPUUsage", func() {
		Context("With system not supports /proc/stat", func() {
			It("should return 0.0", func() {
				u, err := GetCPUUsage()
				if err != nil {
					if err.Error() == "Cannot read CPU data" {
						Expect(u).To(Equal(0.0))
					}
				}
			})
		})
	})
})
