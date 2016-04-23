package goniplus

import "expvar"

func getExpvar() map[string]interface{} {
	m := make(map[string]interface{})
	expvar.Do(func(kv expvar.KeyValue) {
		m[kv.Key] = kv.Value
	})
	return m
}
