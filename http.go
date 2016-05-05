package goniplus

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Request contains Random task id, tag and time information
type Request struct {
	id           string
	method       string
	path         string
	start        time.Time
	finishedAt   string
	response     string
	responseTime int64
}

// RequestData contains response time with timestamp
type RequestData struct {
	ResponseTime int64  `json:"res"`
	Timestamp    string `json:"time"`
}

// Path > Method > RequestData
var reqMap = make(map[string]map[string]map[string][]RequestData)
var reqMapLock = &sync.Mutex{}

func initHTTPMap() {
	reqMap = make(map[string]map[string]map[string][]RequestData)
}

func getHTTPResponseMetric() map[string]map[string]map[string][]RequestData {
	reqMapLock.Lock()
	data := make(map[string]map[string]map[string][]RequestData, len(reqMap))
	for k, v := range reqMap {
		data[k] = v
	}
	initHTTPMap()
	reqMapLock.Unlock()
	return data
}

// Return request id (string) for request tracking
func createRequestID(m, p string) string {
	id := fmt.Sprintf("%v@%s_%s", time.Now().UnixNano(), m, p)
	return id
}

// startRequestTrack starts request tracking
func startRequestTrack(r *http.Request) *Request {
	// Create Request for tracking request
	req := &Request{
		id:     createRequestID(r.Method, r.URL.String()),
		method: r.Method,
		path:   r.URL.String(),
		start:  time.Now(),
	}
	// Add id to request header for tracking breadcrumb inside request
	r.Header.Add("Goni-tracking-id", req.id)
	// TODO : Add id to task queue for current request graph
	return req
}

// finishRequestTrack finishes request tracking
func (r *Request) finishRequestTrack(status int) {
	t := time.Now()
	r.responseTime = int64(t.Sub(r.start) / time.Millisecond)
	r.finishedAt = strconv.FormatInt(t.Unix(), 10)
	r.response = strconv.Itoa(status)
	r.addRequestData()
}

func (r *Request) addRequestData() {
	reqMapLock.Lock()
	if mP, ok := reqMap[r.path]; !ok {
		mP = make(map[string]map[string][]RequestData)
		reqMap[r.path] = mP
	}
	if mM, ok := reqMap[r.path][r.method]; !ok {
		mM = make(map[string][]RequestData)
		reqMap[r.path][r.method] = mM
	}
	reqMap[r.path][r.method][r.response] = append(reqMap[r.path][r.method][r.response], RequestData{
		ResponseTime: r.responseTime,
		Timestamp:    getUnixTimestamp(),
	})
	reqMapLock.Unlock()
}
