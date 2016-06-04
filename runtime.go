package goniplus

import "runtime"

// GetRuntime returns runtime metric map.
func GetRuntime() map[string]interface{} {
	m := make(map[string]interface{})
	m["cgo"] = runtime.NumCgoCall()
	m["goroutine"] = runtime.NumGoroutine()
	return m
}
