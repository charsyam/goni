package goniplus

import "net/http"

// startHTTPTrack create new Task for tracking http request
func startHTTPTrack(r *http.Request) *Task {
	tag := r.Method + "\\" + r.URL.String() // set Method + Path for tag
	return startTask(tag)
}

// finishHTTPTrack finish tracking for http request
func finishHTTPTrack(t *Task, respCode int) {
	t.finishTask()
}
