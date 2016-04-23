package goniplus

import "testing"

func TestGetExpvar(t *testing.T) {
	q := getExpvar()
	if !(len(q) > 0) {
		t.Fail()
	}
}
