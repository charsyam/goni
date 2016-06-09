package goniplus

import (
	"encoding/json"
	"expvar"
	"strings"
)

// ExcludeMemstatKey is a string not to be included
// in expvar memstats metric map. Key is separated by space.
var ExcludeMemstatKey = "BySize PauseEnd PauseNs"

// GetExpvar returns expvar metric map.
// For memstats, key specified in ExcludeMemstatKey will be excluded.
func GetExpvar() map[string]string {
	m := make(map[string]string)
	expvar.Do(func(kv expvar.KeyValue) {
		if kv.Key == "memstats" {
			data := make(map[string]interface{})
			decoder := json.NewDecoder(strings.NewReader(kv.Value.String()))
			if err := decoder.Decode(&data); err != nil {
				return
			}
			excluded := make(map[string]interface{})
			for key, value := range data {
				if !strings.Contains(ExcludeMemstatKey, key) {
					excluded[key] = value
				}
			}
			str, err := json.Marshal(excluded)
			if err != nil {
				return
			}
			m[kv.Key] = string(str)
			return
		}
		m[kv.Key] = kv.Value.String()
	})
	return m
}
