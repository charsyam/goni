package goniplus

import "sync"

var errMap = make(map[string][]string)
var errMapLock = &sync.Mutex{}

func initErrMap() {
	errMap = make(map[string][]string)
}

func getErrorMetric() map[string][]string {
	errMapLock.Lock()
	defer errMapLock.Unlock()
	data := make(map[string][]string)
	for k, v := range errMap {
		data[k] = v
	}
	initErrMap()
	return data
}

// Error add application error to metric array
func Error(tag, err string) {
	errMapLock.Lock()
	defer errMapLock.Unlock()
	errMap[tag] = append(errMap[tag], err)
}
