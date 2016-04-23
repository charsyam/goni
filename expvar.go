package goniplus

import "expvar"

func getExpvar() map[string]string {
	m := make(map[string]string)
	expvar.Do(func(kv expvar.KeyValue) {
		m[kv.Key] = kv.Value.String()
	})
	return m
}
