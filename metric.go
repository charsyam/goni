package goniplus

import (
	"encoding/json"
	"strconv"
	"time"
)

// tempMetric saves temporary data for calculating / collecting data
type tempMetric struct {
	// prevCPUMetric saves calculated total, idle value for next calculation
	prevCPUMetric localCPUMetric
	errMap        map[string][]string
	// Path > Method > Status > Browser > RequestData
	reqMap          map[string]map[string]map[string]map[string][]RequestData
	reqTrackMap     map[string][]string
	reqTrackTimeMap map[string][]time.Time
	reqUserMap      map[string]bool
}

// ApplicationMetric contains http data
type ApplicationMetric struct {
	Error map[string][]string                                       `json:"err"`
	HTTP  map[string]map[string]map[string]map[string][]RequestData `json:"http"`
	User  []string                                                  `json:"user"`
}

// SystemMetric contains expvar data and runtime data
// Expvar : Alloc / Sys / HeapAlloc / HeapInuse / PauseTotalNs / NumGC
// Runtime : cgo / goroutine
type SystemMetric struct {
	Expvar   map[string]interface{} `json:"expvar"`
	Resource map[string]interface{} `json:"resource"`
	Runtime  map[string]interface{} `json:"runtime"`
}

// Metric contains SystemMetric and timestamp
type Metric struct {
	APIKey      string            `json:"apikey"`
	Application ApplicationMetric `json:"app"`
	Instance    string            `json:"instance"`
	System      SystemMetric      `json:"sys"`
	Timestamp   string            `json:"time"`
}

// getTimestamp() returns RFC3339 Timestamp string
func getTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

func getUnixTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func getApplicationMetric() ApplicationMetric {
	http, user := getHTTPResponseMetric()
	metric := ApplicationMetric{
		Error: GetErrorMetric(),
		HTTP:  http,
		User:  user,
	}
	return metric
}

func getSystemMetric() SystemMetric {
	metric := SystemMetric{
		Expvar:   GetExpvar(),
		Resource: getResourceData(),
		Runtime:  getRuntimeData(),
	}
	return metric
}

func (c *Client) getMetric(update bool) ([]byte, error) {
	if update {
		c.id = getInstanceID()
	}
	metric := Metric{
		APIKey:      c.apikey,
		Application: getApplicationMetric(),
		Instance:    c.id,
		System:      getSystemMetric(),
		Timestamp:   getTimestamp(),
	}
	data, err := json.Marshal(metric)
	return data, err
}
