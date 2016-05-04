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

var client *Client
var metricURL = "goniplus.layer123.io:9900"

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

func getPid() string {
	return strconv.Itoa(os.Getpid())
}

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
