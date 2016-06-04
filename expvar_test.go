package goniplus

import (
	"fmt"
	"testing"
)

func TestGetExpvar(t *testing.T) {
	fmt.Println("TestGetExpvar")
	expvar := getExpvar()
	if !(len(expvar) > 0) {
		t.Error("getExpvar() should return more than 1 metric")
	}
}
