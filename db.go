package main

import (
	"encoding/json"
	pb "github.com/goniapm/goniplus-worker/metric"
	influxlib "github.com/influxdata/influxdb/client/v2"
	"log"
	"strconv"
	"time"
)

func addExpvar(data *pb.Metric, timestamp time.Time, bp influxlib.BatchPoints) {
	var memstats map[string]interface{}
	err := json.Unmarshal([]byte(data.System.Expvar["memstats"]), &memstats)
	if err != nil {
		log.Println(err)
		return
	}
	tags := map[string]string{
		"apikey": data.Apikey,
	}
	fields := map[string]interface{}{
		"alloc":        memstats["Alloc"],
		"sys":          memstats["Sys"],
		"heapalloc":    memstats["HeapAlloc"],
		"heapinuse":    memstats["HeapInuse"],
		"pausetotalns": memstats["PauseTotalNs"],
		"numgc":        memstats["NumGC"],
		"instance":     data.Instance,
	}
	pt, err := influxlib.NewPoint("expvar", tags, fields, timestamp)
	if err != nil {
		log.Println(err)
		return
	}
	bp.AddPoint(pt)
}

func addResource(data *pb.Metric, timestamp time.Time, bp influxlib.BatchPoints) {
	tags := map[string]string{
		"apikey":   data.Apikey,
		"instance": data.Instance,
	}
	fields := map[string]interface{}{
		"cpu": data.System.Resource.Cpu,
	}
	pt, err := influxlib.NewPoint("resource", tags, fields, timestamp)
	if err != nil {
		return
	}
	bp.AddPoint(pt)
}

func addRealtime(data *pb.Metric, bp influxlib.BatchPoints) {
	for _, v := range data.Application.Realtime {
		tags := map[string]string{
			"apikey":    data.Apikey,
			"timegroup": strconv.FormatInt(v.Timegroup, 10),
		}
		fields := map[string]interface{}{
			"count": v.Count,
		}
		pt, err := influxlib.NewPoint("realtime", tags, fields)
		if err != nil {
			continue
		}
		bp.AddPoint(pt)
	}
}

func addRuntime(data *pb.Metric, timestamp time.Time, bp influxlib.BatchPoints) {
	tags := map[string]string{
		"apikey": data.Apikey,
	}
	fields := map[string]interface{}{
		"cgo":       data.System.Runtime.Cgo,
		"goroutine": data.System.Runtime.Goroutine,
		"instance":  data.Instance,
	}
	pt, err := influxlib.NewPoint("runtime", tags, fields, timestamp)
	if err != nil {
		return
	}
	bp.AddPoint(pt)
}

func addTransaction(data *pb.Metric, bp influxlib.BatchPoints) {
	detail := data.Application.Transaction.Detail
	for _, transaction := range detail {
		for _, transactionData := range transaction.Data {
			tags := map[string]string{
				"apikey":   data.Apikey,
				"browser":  transactionData.Browser,
				"method":   transaction.Method,
				"path":     transaction.Path,
				"realpath": transaction.Realpath,
				"status":   transactionData.Status.Status,
			}
			fields := map[string]interface{}{
				"res":      transactionData.Status.Duration,
				"instance": data.Instance,
			}
			if len(transactionData.Breadcrumb.Tag) != 0 {
				breadcrumb, err := json.Marshal(transactionData.Breadcrumb.Tag)
				if err != nil {
					continue
				}
				breadcrumbT, err := json.Marshal(transactionData.Breadcrumb.TagT)
				if err != nil {
					continue
				}
				tags["breadcrumb"] = string(breadcrumb)
				fields["breadcrumbT"] = string(breadcrumbT)
			}
			if transactionData.Status.Panic {
				fields["panic"] = true
			}
			t, err := strconv.ParseInt(transactionData.Status.Timestamp, 10, 64)
			if err != nil {
				continue
			}
			timestamp := time.Unix(t, 0)
			pt, err := influxlib.NewPoint("http", tags, fields, timestamp)
			if err != nil {
				continue
			}
			bp.AddPoint(pt)
		}
	}
}

func addUser(data *pb.Metric, timestamp time.Time, bp influxlib.BatchPoints) {
	user := data.Application.User
	if len(user) == 0 {
		return
	}
	tags := map[string]string{
		"apikey": data.Apikey,
	}
	fields := map[string]interface{}{
		"count":    len(user),
		"instance": data.Instance,
	}
	pt, err := influxlib.NewPoint("httpUser", tags, fields, timestamp)
	if err != nil {
		return
	}
	bp.AddPoint(pt)
}

func insertMetric(data *pb.Metric) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered at insertMetric: ", r)
		}
	}()
	bp, err := influxlib.NewBatchPoints(influxlib.BatchPointsConfig{
		Database:  "goniplus",
		Precision: "s",
	})
	if err != nil {
		log.Println(err)
		return
	}
	t, err := strconv.ParseInt(data.Timestamp, 10, 64)
	if err != nil {
		return
	}
	timestamp := time.Unix(t, 0)
	addExpvar(data, timestamp, bp)
	addRealtime(data, bp)
	addResource(data, timestamp, bp)
	addRuntime(data, timestamp, bp)
	addUser(data, timestamp, bp)
	addTransaction(data, bp)
	err = influx.Write(bp)
	if err != nil {
		log.Println(err)
	}
}
