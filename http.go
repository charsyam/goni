package goniplus

import (
	"fmt"
	"github.com/mssola/user_agent"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Request contains Random task id, tag and time information
type Request struct {
	userAgent    string
	id           string
	method       string
	path         string
	start        time.Time
	finishedAt   string
	response     string
	responseTime int64
}

// RequestData contains server response data
type RequestData struct {
	Breadcrumb   []string `json:"crumb,omitempty"`
	Panic        bool     `json:"panic,omitempty"`
	ResponseTime int64    `json:"res"`
	Timestamp    string   `json:"time"`
}

// Path > Method > Status > Browser > RequestData
var reqMap = make(map[string]map[string]map[string]map[string][]RequestData)
var reqMapLock = &sync.Mutex{}
var reqTrackMap = make(map[string][]string)
var reqTrackMapLock = &sync.Mutex{}

func initHTTPMap() {
	reqMap = make(map[string]map[string]map[string]map[string][]RequestData)
}

func getHTTPResponseMetric() map[string]map[string]map[string]map[string][]RequestData {
	reqMapLock.Lock()
	defer reqMapLock.Unlock()
	respData := make(map[string]map[string]map[string]map[string][]RequestData, len(reqMap))
	for k, v := range reqMap {
		respData[k] = v
	}
	initHTTPMap()
	return respData
}

// Return request id (string) for request tracking
func createRequestID(m, p string) string {
	id := fmt.Sprintf("%v@%s_%s", time.Now().UnixNano(), m, p)
	return id
}

// StartRequestTrack starts request tracking
func StartRequestTrack(r *http.Request) *Request {
	// Create Request for tracking request
	req := &Request{
		id:        createRequestID(r.Method, r.URL.String()),
		method:    r.Method,
		path:      r.URL.String(),
		start:     time.Now(),
		userAgent: r.Header.Get("User-Agent"),
	}
	// Add id to request header for tracking breadcrumb inside request
	r.Header.Add("Goni-tracking-id", req.id)
	// TODO : Add id to task queue for current request graph
	return req
}

// LeaveBreadcrumb add specified tag to array
func LeaveBreadcrumb(r *http.Request, tag string) {
	// Get tracking id from header
	id := r.Header.Get("Goni-tracking-id")
	reqTrackMapLock.Lock()
	defer reqTrackMapLock.Unlock()
	if arr, ok := reqTrackMap[id]; !ok {
		arr = make([]string, 0)
		reqTrackMap[id] = arr
	}
	reqTrackMap[id] = append(reqTrackMap[id], tag)
}

// FinishRequestTrack finishes request tracking
func (r *Request) FinishRequestTrack(status int, panic bool) {
	t := time.Now()
	r.responseTime = int64(t.Sub(r.start) / time.Millisecond)
	r.finishedAt = strconv.FormatInt(t.Unix(), 10)
	r.response = strconv.Itoa(status)
	r.addRequestData(panic)
}

func (r *Request) addRequestData(panic bool) {
	reqMapLock.Lock()
	defer reqMapLock.Unlock()
	// Map check
	if mP, ok := reqMap[r.path]; !ok {
		mP = make(map[string]map[string]map[string][]RequestData)
		reqMap[r.path] = mP
	}
	if mM, ok := reqMap[r.path][r.method]; !ok {
		mM = make(map[string]map[string][]RequestData)
		reqMap[r.path][r.method] = mM
	}
	if mR, ok := reqMap[r.path][r.method][r.response]; !ok {
		mR = make(map[string][]RequestData)
		reqMap[r.path][r.method][r.response] = mR
	}
	// UserAgent
	ua := user_agent.New(r.userAgent)
	browserName, browserVersion := ua.Browser()
	browser := fmt.Sprintf("%s_%s", browserName, browserVersion)
	reqTrackMapLock.Lock()
	defer reqTrackMapLock.Unlock()
	reqMap[r.path][r.method][r.response][browser] = append(reqMap[r.path][r.method][r.response][browser], RequestData{
		Breadcrumb:   reqTrackMap[r.id],
		Panic:        panic,
		ResponseTime: r.responseTime,
		Timestamp:    getUnixTimestamp(),
	})
	delete(reqTrackMap, r.id)
}
