package goniplus

import "runtime"

func getRuntimeData() map[string]interface{} {
	m := make(map[string]interface{})
	m["cgo"] = runtime.NumCgoCall()
	m["goroutine"] = runtime.NumGoroutine()
	return m
}
