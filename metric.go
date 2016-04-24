package goniplus

import (
	"encoding/json"
	"time"
)

// SystemMetric contains expvar data and runtime data
//
// Expvar : Alloc / Sys / HeapAlloc / HeapInuse / PauseTotalNs / NumGC
// Runtime : cgo / goroutine
type SystemMetric struct {
	Expvar  map[string]interface{} `json:"expvar"`
	Runtime map[string]interface{} `json:"runtime"`
}

// Metric contains SystemMetric and timestamp
type Metric struct {
	System    SystemMetric `json:"sys"`
	Timestamp string       `json:"time"`
}

// getTime() returns RFC3339 Timestamp string
func getTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

func getSystemMetric() SystemMetric {
	metric := SystemMetric{
		Expvar:  getExpvar(),
		Runtime: getRuntimeData(),
	}
	return metric
}

func getMetric() ([]byte, error) {
	metric := Metric{
		System:    getSystemMetric(),
		Timestamp: getTimestamp(),
	}
	data, err := json.Marshal(metric)
	return data, err
}
