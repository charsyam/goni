package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	pb "github.com/monitflux/goniplus-worker/metric"
	influxlib "github.com/influxdata/influxdb/client/v2"
	"io/ioutil"
	"log"
	"net"
)

var influx influxlib.Client
var mySQL *sql.DB

func handleData(conn net.Conn, dbQueue, slackQueue chan *pb.Metric) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered at handleData: ", r)
		}
	}()
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		return
	}
	db := &pb.Metric{}
	if err = proto.Unmarshal(b, db); err != nil {
		log.Println("Failed to parse metric:", err)
		return
	}
	slack := &pb.Metric{}
	if err = proto.Unmarshal(b, slack); err != nil {
		log.Println("Failed to parse metric:", err)
		return
	}
	dbQueue <- db
	slackQueue <- slack
}

func main() {
	influxConn, err := getInflux()
	if err != nil {
		log.Fatalln(err)
		return
	}
	influx = influxConn
	mySQLConn, err := getMySQL()
	if err != nil {
		log.Fatalln(err)
		return
	}
	mySQL = mySQLConn
	defer mySQL.Close()
	maxWorkers := 2
	queueSize := 256
	dbQueue := make(chan *pb.Metric, queueSize)
	slackQueue := make(chan *pb.Metric, queueSize)
	dispatcher := newDispatcher(dbQueue, slackQueue, maxWorkers)
	dispatcher.run()
	ln, err := net.Listen("tcp", ":9900")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleData(conn, dbQueue, slackQueue)
	}
}
