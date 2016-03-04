package goniplus

import (
	"math/rand"
	"time"
)

// Task contains Random task id, tag and time information
type Task struct {
	id       int64
	tag      string
	isSent   bool
	start    time.Time
	duration time.Duration
}

// Return non-negative random taskid (int64)
// taskID can return duplicated value
func getTaskID() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63()
}

// startTask returns new task
func startTask(taskTag string) *Task {
	i := getTaskID()
	// Create task
	t := &Task{
		id:    i,
		tag:   taskTag,
		start: time.Now(),
	}
	// TODO : Add task to taskMap (Server Queue)
	return t
}

// finishTask add elapsed time to task
func (t *Task) finishTask() {
	t.duration = time.Since(t.start)
	// TODO : Remove task from taskMap and add to finishedTaskMap
}
