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

var transactionMapLock = &sync.Mutex{}
var reqIDMapLock = &sync.Mutex{}
var reqTrackMapLock = &sync.Mutex{}
var reqTrackTimeMapLock = &sync.Mutex{}
var reqUserMapLock = &sync.Mutex{}

func initHTTPMap(forceInit bool) {
	if forceInit {
		// reqMap
		transactionMapLock.Lock()
		client.tMetric.transactionMap = make(map[string]*pb.ApplicationMetric_TransactionDetail)
		transactionMapLock.Unlock()
		// reqIDMap
		reqIDMapLock.Lock()
		client.tMetric.reqIDMap = make(map[string]*Request)
		reqIDMapLock.Unlock()
		// reqTrackMap
		reqTrackMapLock.Lock()
		client.tMetric.reqTrackMap = make(map[string][]string)
		reqTrackMapLock.Unlock()
		// reqTrackTimeMap
		reqTrackTimeMapLock.Lock()
		client.tMetric.reqTrackTimeMap = make(map[string][]time.Time)
		reqTrackTimeMapLock.Unlock()
	}
	// reqUserMap
	reqUserMapLock.Lock()
	client.tMetric.reqUserMap = make(map[string]bool)
	reqUserMapLock.Unlock()
}

// GetTransactionMetric returns http metric map.
func GetTransactionMetric() (*pb.ApplicationMetric_Transaction, []*pb.ApplicationMetric_Realtime, []*pb.ApplicationMetric_User) {
	transactionMapLock.Lock()
	reqUserMapLock.Lock()
	reqIDMapLock.Lock()
	var transactionMetric []*pb.ApplicationMetric_TransactionDetail
	for k, v := range client.tMetric.transactionMap {
		transactionMetric = append(transactionMetric, v)
		delete(client.tMetric.transactionMap, k)
	}
	realtimeMap := make(map[int64]int64)
	var realtimeMetric []*pb.ApplicationMetric_Realtime
	for _, v := range client.tMetric.reqIDMap {
		k := int64(time.Since(v.start)/time.Millisecond) / 250
		if k > 60 {
			k = 60
		}
		realtimeMap[k]++
	}
	for k, v := range realtimeMap {
		realtimeMetric = append(realtimeMetric, &pb.ApplicationMetric_Realtime{
			Timegroup: k,
			Count:     v,
		})
	}
	var userMetric []*pb.ApplicationMetric_User
	for user := range client.tMetric.reqUserMap {
		userMetric = append(userMetric, &pb.ApplicationMetric_User{
			Ip: user,
		})
	}
	reqIDMapLock.Unlock()
	reqUserMapLock.Unlock()
	transactionMapLock.Unlock()
	initHTTPMap(false)
	return &pb.ApplicationMetric_Transaction{Detail: transactionMetric}, realtimeMetric, userMetric
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
	reqIDMapLock.Lock()
	client.tMetric.reqIDMap[req.id] = req
	reqIDMapLock.Unlock()
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
	reqIDMapLock.Lock()
	delete(client.tMetric.reqIDMap, r.id)
	reqIDMapLock.Unlock()
	reqKey := fmt.Sprintf("%s (%s)", r.path, r.method)
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
	transactionMapLock.Lock()
	if _, ok := client.tMetric.transactionMap[reqKey]; !ok {
		client.tMetric.transactionMap[reqKey] = &pb.ApplicationMetric_TransactionDetail{
			Path:     reqKey,
			Data:     make([]*pb.ApplicationMetric_TransactionData, 0),
			Method:   r.method,
			Realpath: r.path,
		}
	}
	detail := client.tMetric.transactionMap[reqKey]
	detail.Data = append(detail.Data, &pb.ApplicationMetric_TransactionData{
		Status: &pb.ApplicationMetric_TransactionStatus{
			Status:    r.response,
			Duration:  r.responseTime,
			Panic:     panic,
			Timestamp: strconv.FormatInt(r.start.Unix(), 10),
		},
		Browser: browser,
		Breadcrumb: &pb.ApplicationMetric_Breadcrumb{
			Tag:  client.tMetric.reqTrackMap[r.id],
			TagT: crumbT,
		},
	})
	transactionMapLock.Unlock()
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
