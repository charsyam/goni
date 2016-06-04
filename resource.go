package goniplus

var nonInitialData bool

func getResourceData() map[string]interface{} {
	m := make(map[string]interface{})
	cpu, err := GetCPUUsage()
	if err == nil {
		if !nonInitialData {
			nonInitialData = true
		} else {
			m["cpu"] = cpu
		}
	}
	return m
}
