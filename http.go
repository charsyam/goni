package goniplus

import (
	"fmt"
	pb "github.com/goniapm/goniplus-worker/metric"
	"github.com/mssola/user_agent"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Request is a type that wraps request information
type Request struct {
	userAgent    string
	id           string
	ip           string
	method       string
	path         string
	start        time.Time
	finishedAt   string
	response     string
	responseTime int64
}

// RequestData is a type for keeping request tracking data
type RequestData struct {
	Breadcrumb     []string `json:"crumb,omitempty"`
	BreadcrumbTime []int64  `json:"crumbT,omitempty"`
	Panic          bool     `json:"panic,omitempty"`
	ResponseTime   int64    `json:"res"`
	Timestamp      string   `json:"time"`
}

var reqMapLock = &sync.Mutex{}
var reqBrowserMapLock = &sync.Mutex{}
var reqTrackMapLock = &sync.Mutex{}
var reqTrackTimeMapLock = &sync.Mutex{}
var reqUserMapLock = &sync.Mutex{}

func initHTTPMap() {
	// reqMap
	reqMapLock.Lock()
	client.tMetric.reqMap = make(map[string]*pb.ApplicationMetric_HTTPDetail)
	reqMapLock.Unlock()
	// reqBrowserMap
	reqBrowserMapLock.Lock()
	client.tMetric.reqBrowserMap = make(map[string]map[string]uint32)
	reqBrowserMapLock.Unlock()
	// reqTrackMap
	reqTrackMapLock.Lock()
	client.tMetric.reqTrackMap = make(map[string][]string)
	reqTrackMapLock.Unlock()
	// reqTrackTimeMap
	reqTrackTimeMapLock.Lock()
	client.tMetric.reqTrackTimeMap = make(map[string][]time.Time)
	reqTrackTimeMapLock.Unlock()
	// reqUserMap
	reqUserMapLock.Lock()
	client.tMetric.reqUserMap = make(map[string]bool)
	reqUserMapLock.Unlock()
}

// GetHTTPResponseMetric returns http metric map.
func GetHTTPResponseMetric() (*pb.ApplicationMetric_HTTP, []*pb.ApplicationMetric_User) {
	reqMapLock.Lock()
	reqUserMapLock.Lock()
	reqBrowserMapLock.Lock()
	var reqMetric []*pb.ApplicationMetric_HTTPDetail
	for _, v := range client.tMetric.reqMap {
		browserMap := client.tMetric.reqBrowserMap[v.Path]
		var browserMetric []*pb.ApplicationMetric_Browser
		for browser, count := range browserMap {
			browserMetric = append(browserMetric, &pb.ApplicationMetric_Browser{
				Browser: browser,
				Count:   count,
			})
		}
		v.Browser = browserMetric
		reqMetric = append(reqMetric, v)
	}
	reqBrowserMapLock.Unlock()
	var userMetric []*pb.ApplicationMetric_User
	for user := range client.tMetric.reqUserMap {
		userMetric = append(userMetric, &pb.ApplicationMetric_User{
			Ip: user,
		})
	}
	reqUserMapLock.Unlock()
	reqMapLock.Unlock()
	initHTTPMap()
	return &pb.ApplicationMetric_HTTP{Detail: reqMetric}, userMetric
}

// CreateRequestID returns request id (string) for request tracking
func CreateRequestID(m, p string) string {
	id := fmt.Sprintf("%v@%s_%s", time.Now().UnixNano(), m, p)
	return id
}

// StartRequestTrack starts request tracking
func StartRequestTrack(r *http.Request) *Request {
	// Create Request for tracking request
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	req := &Request{
		id:        CreateRequestID(r.Method, r.URL.String()),
		ip:        ip,
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
	if arr, ok := client.tMetric.reqTrackMap[id]; !ok {
		arr = make([]string, 0)
		client.tMetric.reqTrackMap[id] = arr
	}
	client.tMetric.reqTrackMap[id] = append(client.tMetric.reqTrackMap[id], tag)
	reqTrackMapLock.Unlock()
	reqTrackTimeMapLock.Lock()
	client.tMetric.reqTrackTimeMap[id] = append(client.tMetric.reqTrackTimeMap[id], time.Now())
	reqTrackTimeMapLock.Unlock()
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
	reqKey := r.method + " " + r.path
	// UserAgent
	ua := user_agent.New(r.userAgent)
	browserName, browserVersion := ua.Browser()
	browser := fmt.Sprintf("%s_%s", browserName, browserVersion)
	// IP
	reqUserMapLock.Lock()
	if _, ok := client.tMetric.reqUserMap[r.ip]; !ok {
		client.tMetric.reqUserMap[r.ip] = true
	}
	reqUserMapLock.Unlock()
	reqTrackTimeMapLock.Lock()
	var crumbT []int64
	var totalD int64
	crumbLen := len(client.tMetric.reqTrackTimeMap[r.id])
	for i := 0; i < crumbLen; i++ {
		if i == 0 {
			d := int64(client.tMetric.reqTrackTimeMap[r.id][i].Sub(r.start) / time.Millisecond)
			totalD += d
			crumbT = append(crumbT, d)
			if i == crumbLen-1 {
				crumbT = append(crumbT, r.responseTime-totalD)
			}
			continue
		}
		d := int64(client.tMetric.reqTrackTimeMap[r.id][i].Sub(
			client.tMetric.reqTrackTimeMap[r.id][i-1]) / time.Millisecond)
		crumbT = append(crumbT, d)
		totalD += d
		if i == crumbLen-1 {
			crumbT = append(crumbT, r.responseTime-totalD)
		}
	}
	delete(client.tMetric.reqTrackTimeMap, r.id)
	reqTrackTimeMapLock.Unlock()
	reqTrackMapLock.Lock()
	reqMapLock.Lock()
	if _, ok := client.tMetric.reqMap[reqKey]; !ok {
		client.tMetric.reqMap[reqKey] = &pb.ApplicationMetric_HTTPDetail{
			Path:    reqKey,
			Status:  make([]*pb.ApplicationMetric_HTTPStatus, 0),
			Browser: make([]*pb.ApplicationMetric_Browser, 0),
			Breadcrumb: &pb.ApplicationMetric_Breadcrumb{
				Crumb: make([]*pb.ApplicationMetric_BreadcrumbDetail, 0),
			},
		}
	}
	detail := client.tMetric.reqMap[reqKey]
	detail.Status = append(detail.Status, &pb.ApplicationMetric_HTTPStatus{
		Status:    r.response,
		Duration:  r.responseTime,
		Panic:     panic,
		Timestamp: r.start.Format(time.RFC3339),
	})
	detail.Breadcrumb.Crumb =
		append(detail.Breadcrumb.Crumb, &pb.ApplicationMetric_BreadcrumbDetail{
			Tag:  client.tMetric.reqTrackMap[r.id],
			TagT: crumbT,
		})
	reqMapLock.Unlock()
	reqBrowserMapLock.Lock()
	if _, ok := client.tMetric.reqBrowserMap[reqKey]; !ok {
		browserMap := make(map[string]uint32)
		client.tMetric.reqBrowserMap[reqKey] = browserMap
	}
	client.tMetric.reqBrowserMap[reqKey][browser]++
	reqBrowserMapLock.Unlock()
	delete(client.tMetric.reqTrackMap, r.id)
	reqTrackMapLock.Unlock()
}

/*
 * net/http middleware
 */

// ResponseWriterWrapper wrap ResponseWriter for net/http middleware support
type ResponseWriterWrapper struct {
	w      http.ResponseWriter
	status int
	size   int
}

// Header is implementaion of net/http ResponseWriter Header()
func (rww *ResponseWriterWrapper) Header() http.Header {
	return rww.w.Header()
}

// Write is implementaion of snet/http ResponseWriter Write()
func (rww *ResponseWriterWrapper) Write(b []byte) (int, error) {
	if rww.status == 0 {
		rww.status = http.StatusOK
	}
	size, err := rww.w.Write(b)
	rww.size += size
	return size, err
}

// WriteHeader sets the header with status
func (rww *ResponseWriterWrapper) WriteHeader(s int) {
	rww.w.WriteHeader(s)
	rww.status = s
}

// Status returns the HTTP status code
func (rww *ResponseWriterWrapper) Status() int {
	return rww.status
}

// Middleware returns net/http middleware
func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := StartRequestTrack(r)
		defer func() {
			if err := recover(); err != nil {
				t.FinishRequestTrack(500, true)
			}
		}()
		rww := &ResponseWriterWrapper{w: w}
		h(rww, r)
		t.FinishRequestTrack(rww.status, false)
	})
}
