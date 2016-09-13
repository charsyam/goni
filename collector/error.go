package goniplus

import (
	pb "github.com/goniapm/goniplus-worker/metric"
	"sync"
)

var errMapLock = &sync.Mutex{}

// initErrMap initialize error map.
func initErrMap() {
	errMapLock.Lock()
	client.tMetric.errMap = make([]*pb.ApplicationMetric_Error, 0)
	errMapLock.Unlock()
}

// GetErrorMetric returns error map, and initialize error map.
func GetErrorMetric() []*pb.ApplicationMetric_Error {
	errMapLock.Lock()
	errMetric := make([]*pb.ApplicationMetric_Error, len(client.tMetric.errMap))
	copy(errMetric, client.tMetric.errMap)
	errMapLock.Unlock()
	initErrMap()
	return errMetric
}

// Error add passed tag, error (explanation) to error map.
func Error(tag, err string) {
	errMapLock.Lock()
	client.tMetric.errMap = append(client.tMetric.errMap, &pb.ApplicationMetric_Error{
		Tag:       tag,
		Log:       err,
		Timestamp: GetUnixTimestamp(),
	})
	errMapLock.Unlock()
}
