package goniplus

import (
	"net"
	"os"
	"strconv"
	"time"
)

// Client for server communication
type Client struct {
	apikey       string
	id           string
	sendInterval time.Duration
}

// Global GoniPlus client
var client *Client
var metricURL = "goniplus.layer123.io:9900"

// getIP() returns system's IP address in string
func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// getPid() returns application's Pid in string
func getPid() string {
	return strconv.Itoa(os.Getpid())
}

// getInstance() returns instance id
func getInstanceID() string {
	return getIP() + "-" + getPid()
}

// InitSDK initialize goniplus sdk client
func InitSDK(apikey string, interval int) {
	client = &Client{
		apikey,
		getInstanceID(),
		time.Duration(interval) * time.Second,
	}
	go client.startSender()
}
