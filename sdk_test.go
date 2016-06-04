package goniplus

import (
	"fmt"
	"testing"
	"time"
)

func TestGetIP(t *testing.T) {
	fmt.Println("TestGetIP")
	ip := getIP()
	if !(len(ip) > 0) {
		t.Error("Cannot get IP")
	}
}

func TestGetPID(t *testing.T) {
	fmt.Println("TestGetPID")
	pid := getPid()
	if !(len(pid) > 0) {
		t.Error("Cannot get PID")
	}
}

func TestGetInstanceID(t *testing.T) {
	fmt.Println("TestGetInstanceID")
	iid := getInstanceID()
	if !(len(iid) > 2) {
		t.Error("Instance ID length should be more than 2")
	}
}

func TestInitSDK(t *testing.T) {
	fmt.Println("TestInitSDK")
	InitSDK("APIKEY", 60)
	if client.apikey != "APIKEY" {
		t.Error("APIKey not configured to client")
	}
	t.Log("Send interval unit should be second")
	if client.sendInterval/time.Second != 60 {
		t.Error("Send interval unit should be second")
	}
}
