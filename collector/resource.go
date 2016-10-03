package goniplus

import pb "github.com/monitflux/goniplus-worker/metric"

// GetResource returns resource metric
func GetResource() *pb.SystemMetric_Resource {
	resourceMetric := pb.SystemMetric_Resource{}
	cpu, _ := GetCPUUsage()
	if !client.tMetric.isResourceInitialCollect {
		client.tMetric.isResourceInitialCollect = true
	} else {
		resourceMetric.Cpu = cpu
	}
	return &resourceMetric
}
