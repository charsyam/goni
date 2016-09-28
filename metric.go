package goniplus

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/monitflux/goniplus-worker/metric"
	"strconv"
	"time"
)

// tempMetric contains temporary data for calculating / collecting data.
type tempMetric struct {
	errMap                   []*pb.ApplicationMetric_Error
	isResourceInitialCollect bool
	// prevCPUMetric saves calculated total, idle value for next calculation
	prevCPUMetric   localCPUMetric
	reqBrowserMap   map[string]map[string]uint32
	reqIDMap        map[string]*Request
	reqTrackMap     map[string][]string
	reqTrackTimeMap map[string][]time.Time
	reqUserMap      map[string]bool
	transactionMap  map[string]*pb.ApplicationMetric_TransactionDetail
}

// GetUnixTimestamp returns UnixTimestamp in string.
func GetUnixTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// GetApplicationMetric returns ApplicationMetric.
func GetApplicationMetric() *pb.ApplicationMetric {
	transaction, realtime, user := GetTransactionMetric()
	appMetric := &pb.ApplicationMetric{
		Error:       GetErrorMetric(),
		Realtime:    realtime,
		Transaction: transaction,
		User:        user,
	}
	return appMetric
}

// GetSystemMetric returns SystemMetric.
func GetSystemMetric() *pb.SystemMetric {
	systemMetric := &pb.SystemMetric{
		Expvar:   GetExpvar(),
		Resource: GetResource(),
		Runtime:  GetRuntime(),
	}
	return systemMetric
}

// GetMetric returns marshalled metric.
func (c *Client) getMetric(update bool) ([]byte, error) {
	if update {
		c.id = getInstanceID()
	}
	metric := &pb.Metric{
		Apikey:      c.apikey,
		Instance:    c.id,
		Timestamp:   GetUnixTimestamp(),
		Application: GetApplicationMetric(),
		System:      GetSystemMetric(),
	}
	data, err := proto.Marshal(metric)
	return data, err
}
