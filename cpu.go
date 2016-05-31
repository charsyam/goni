package goniplus

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

type localCPUMetric struct {
	total, idle uint64
}

var prevCPUMetric localCPUMetric

func calcCPUUsage(fields []string) float64 {
	var idle, total uint64
	for i := 1; i < len(fields); i++ {
		u, _ := strconv.ParseUint(fields[i], 10, 64)
		if i == 4 || i == 5 {
			idle += u
		}
		total += u
	}
	v := float64((total-prevCPUMetric.total)-(idle-prevCPUMetric.idle)) / float64(total-prevCPUMetric.total)
	prevCPUMetric.idle = idle
	prevCPUMetric.total = total
	return v
}

func getCPUUsage() (float64, error) {
	d, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0.0, errors.New("Cannot read CPU data")
	}
	data := string(d)
	lines := strings.Split(data, "\n")
	fields := strings.Fields(lines[0])
	if fields[0][:3] == "cpu" {
		usage := calcCPUUsage(fields)
		return usage, nil
	}
	return 0.0, errors.New("Cannot parse CPU data")
}
