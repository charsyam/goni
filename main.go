package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	pb "github.com/goniapm/goniplus-worker/metric"
	"io/ioutil"
	"log"
	"net"
)

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
	metric := &pb.Metric{}
	if err = proto.Unmarshal(b, metric); err != nil {
		log.Println("Failed to parse metric:", err)
		return
	}
	data := metric
	dbQueue <- data
	slackQueue <- data
}

func main() {
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
