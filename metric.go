package goniplus

import (
	"encoding/json"
	"strconv"
	"time"
)

// tempMetric contains temporary data for calculating / collecting data.
type tempMetric struct {
	errMap                   map[string][]string
	isResourceInitialCollect bool
	// prevCPUMetric saves calculated total, idle value for next calculation
	prevCPUMetric localCPUMetric
	// Path > Method > Status > Browser > RequestData
	reqMap          map[string]map[string]map[string]map[string][]RequestData
	reqTrackMap     map[string][]string
	reqTrackTimeMap map[string][]time.Time
	reqUserMap      map[string]bool
}

// ApplicationMetric contains error / http / user data.
type ApplicationMetric struct {
	Error map[string][]string                                       `json:"err"`
	HTTP  map[string]map[string]map[string]map[string][]RequestData `json:"http"`
	User  []string                                                  `json:"user"`
}

// SystemMetric contains expvar / resource / runtime data.
type SystemMetric struct {
	Expvar   map[string]interface{} `json:"expvar"`
	Resource map[string]interface{} `json:"resource"`
	Runtime  map[string]interface{} `json:"runtime"`
}

// Metric contains SystemMetric and timestamp.
type Metric struct {
	APIKey      string            `json:"apikey"`
	Application ApplicationMetric `json:"app"`
	Instance    string            `json:"instance"`
	System      SystemMetric      `json:"sys"`
	Timestamp   string            `json:"time"`
}

// GetTimestamp returns RFC3339 Timestamp in string.
func GetTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// GetUnixTimestamp returns UnixTimestamp in string.
func GetUnixTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// GetApplicationMetric returns ApplicationMetric.
func GetApplicationMetric() ApplicationMetric {
	http, user := GetHTTPResponseMetric()
	metric := ApplicationMetric{
		Error: GetErrorMetric(),
		HTTP:  http,
		User:  user,
	}
	return metric
}

// GetSystemMetric returns SystemMetric.
func GetSystemMetric() SystemMetric {
	metric := SystemMetric{
		Expvar:   GetExpvar(),
		Resource: GetResource(),
		Runtime:  GetRuntime(),
	}
	return metric
}

// GetMetric returns marshalled metric.
func (c *Client) GetMetric(update bool) ([]byte, error) {
	if update {
		c.id = getInstanceID()
	}
	metric := Metric{
		APIKey:      c.apikey,
		Application: GetApplicationMetric(),
		Instance:    c.id,
		System:      GetSystemMetric(),
		Timestamp:   GetTimestamp(),
	}
	data, err := json.Marshal(metric)
	return data, err
}
