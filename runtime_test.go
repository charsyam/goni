package goniplus

import "testing"

func TestGetRuntimeData(t *testing.T) {
	q := getRuntimeData()
	if !(len(q) > 0) {
		t.Fail()
	}
}
