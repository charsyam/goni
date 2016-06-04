package goniplus

import "sync"

var errMapLock = &sync.Mutex{}

// initErrMap initialize error map.
func initErrMap() {
	errMapLock.Lock()
	client.tMetric.errMap = make(map[string][]string)
	errMapLock.Unlock()
}

// GetErrorMetric returns error map, and initialize error map.
func GetErrorMetric() map[string][]string {
	errMapLock.Lock()
	data := make(map[string][]string)
	for k, v := range client.tMetric.errMap {
		data[k] = v
	}
	errMapLock.Unlock()
	initErrMap()
	return data
}

// Error add passed tag, error (explanation) to error map.
func Error(tag, err string) {
	errMapLock.Lock()
	defer errMapLock.Unlock()
	client.tMetric.errMap[tag] = append(client.tMetric.errMap[tag], err)
}
