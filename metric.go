package goniplus

import (
	"encoding/json"
	"strconv"
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
	APIKey    string       `json:"apikey"`
	Instance  string       `json:"instance"`
	System    SystemMetric `json:"sys"`
	Timestamp string       `json:"time"`
}

// getTimestamp() returns RFC3339 Timestamp string
func getTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

func getUnixTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func getSystemMetric() SystemMetric {
	metric := SystemMetric{
		Expvar:  getExpvar(),
		Runtime: getRuntimeData(),
	}
	return metric
}

func (c *Client) getMetric(update bool) ([]byte, error) {
	if update {
		c.id = getInstanceID()
	}
	metric := Metric{
		APIKey:    c.apikey,
		Instance:  c.id,
		System:    getSystemMetric(),
		Timestamp: getTimestamp(),
	}
	data, err := json.Marshal(metric)
	return data, err
}
