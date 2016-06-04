package goniplus

import "sync"

var errMapLock = &sync.Mutex{}

func initErrMap() {
	errMapLock.Lock()
	client.tMetric.errMap = make(map[string][]string)
	errMapLock.Unlock()
}

func getErrorMetric() map[string][]string {
	errMapLock.Lock()
	data := make(map[string][]string)
	for k, v := range client.tMetric.errMap {
		data[k] = v
	}
	errMapLock.Unlock()
	initErrMap()
	return data
}

// Error add application error to metric array
func Error(tag, err string) {
	errMapLock.Lock()
	defer errMapLock.Unlock()
	client.tMetric.errMap[tag] = append(client.tMetric.errMap[tag], err)
}
