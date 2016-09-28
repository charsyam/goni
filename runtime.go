package goniplus

import (
	pb "github.com/monitflux/goniplus-worker/metric"
	"runtime"
)

// GetRuntime returns runtime metric map.
func GetRuntime() *pb.SystemMetric_Runtime {
	runtimeMetric := &pb.SystemMetric_Runtime{}
	runtimeMetric.Cgo = runtime.NumCgoCall()
	runtimeMetric.Goroutine = int32(runtime.NumGoroutine())
	return runtimeMetric
}
