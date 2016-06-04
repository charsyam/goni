package goniplus

// GetResource returns resource metric map.
func GetResource() map[string]interface{} {
	m := make(map[string]interface{})
	cpu, _ := GetCPUUsage()
	if !client.tMetric.isResourceInitialCollect {
		client.tMetric.isResourceInitialCollect = true
	} else {
		m["cpu"] = cpu
	}
	return m
}
