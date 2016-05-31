package goniplus

var nonInitialData bool

func getResourceData() map[string]interface{} {
	m := make(map[string]interface{})
	cpu, err := getCPUUsage()
	if err == nil {
		if !nonInitialData {
			nonInitialData = true
		} else {
			m["cpu"] = cpu
		}
	}
	return m
}
