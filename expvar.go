package goniplus

import (
	"encoding/json"
	"expvar"
	"strings"
)

var excludeExpvarKey = map[string]bool{
	"BySize":   true,
	"PauseEnd": true,
	"PauseNs":  true,
}

func getExpvar() map[string]interface{} {
	m := make(map[string]interface{})
	expvar.Do(func(kv expvar.KeyValue) {
		if kv.Key == "memstats" {
			data := make(map[string]interface{})
			decoder := json.NewDecoder(strings.NewReader(kv.Value.String()))
			if err := decoder.Decode(&data); err != nil {
				return
			}
			excluded := make(map[string]interface{})
			for key, value := range data {
				if _, ok := excludeExpvarKey[key]; !ok {
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
