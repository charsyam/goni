package main

import pb "github.com/goniapm/goniplus-worker/metric"

const (
	typeDB = iota
	typeSlack
)

type dispatcher struct {
	dbQueue         chan *pb.Metric
	dbWorkerPool    chan chan *pb.Metric
	slackQueue      chan *pb.Metric
	slackWorkerPool chan chan *pb.Metric
	maxWorkers      int
}

type worker struct {
	workerType  int
	metricQueue chan *pb.Metric
	workerPool  chan chan *pb.Metric
	quit        chan bool
}

func newDispatcher(dbQueue, slackQueue chan *pb.Metric, maxWorkers int) *dispatcher {
	dbPool := make(chan chan *pb.Metric, maxWorkers)
	slackPool := make(chan chan *pb.Metric, maxWorkers)
	return &dispatcher{
		dbQueue:         dbQueue,
		dbWorkerPool:    dbPool,
		slackQueue:      slackQueue,
		slackWorkerPool: slackPool,
		maxWorkers:      maxWorkers,
	}
}

func (d *dispatcher) run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := newWorker(typeDB, d.dbWorkerPool)
		worker.start()
		worker = newWorker(typeSlack, d.slackWorkerPool)
		worker.start()
	}

	go d.dispatch()
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case data := <-d.dbQueue:
			go func(data *pb.Metric) {
				dataChannel := <-d.dbWorkerPool
				dataChannel <- data
			}(data)
		case data := <-d.slackQueue:
			go func(data *pb.Metric) {
				dataChannel := <-d.slackWorkerPool
				dataChannel <- data
			}(data)
		}
	}
}

func newWorker(workerType int, workerPool chan chan *pb.Metric) worker {
	return worker{
		workerType:  workerType,
		metricQueue: make(chan *pb.Metric),
		workerPool:  workerPool,
		quit:        make(chan bool),
	}
}

func (w worker) start() {
	go func() {
		for {
			w.workerPool <- w.metricQueue
			select {
			case metric := <-w.metricQueue:
				switch w.workerType {
				case typeDB:
					insertMetric(metric)
					break
				case typeSlack:
					sendSlackNotification(metric)
					break
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w worker) stop() {
	go func() {
		w.quit <- true
	}()
}
