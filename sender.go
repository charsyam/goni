package goniplus

import (
	"log"
	"net"
	"time"
)

// sendMetricData(data []byte) returns true if data successfully sent
func sendMetricData(data []byte) bool {
	conn, err := net.Dial("tcp", metricURL)
	if err != nil {
		log.Println(err)
		return false
	}
	defer conn.Close()
	if _, err = conn.Write(data); err != nil {
		log.Println(err)
		return false
	}
	return true
}

// startSender starts timer for sending metrics to server
// every `sendInterval` second
func (c *Client) startSender() {
	cnt := 0
	for range time.Tick(c.sendInterval) {
		updateInstanceID := false
		if cnt == 5 {
			updateInstanceID = true
			cnt = 0
		} else {
			cnt++
		}
		data, err := c.GetMetric(updateInstanceID)
		if err != nil {
			log.Println(err)
			continue
		}
		go sendMetricData(data)
	}
}
