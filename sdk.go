package goniplus

import "time"

// Client for server communication
type Client struct {
	apikey       string
	sendInterval time.Duration
}

var client *Client

// InitSDK initialize goniplus sdk client
func InitSDK(apikey string, interval int) {
	client = &Client{
		apikey,
		time.Duration(interval) * time.Second,
	}
}
