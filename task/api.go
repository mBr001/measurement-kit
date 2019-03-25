// Package task implements the Measurement Kit task API
package task

import (
	"context"
	"sync/atomic"
)

// Task is a task run by Measurement Kit
type Task struct {
	cancel context.CancelFunc
	ch     chan string
	ctx    context.Context
	done   int64
}

// Start starts a task with the specified settings.
func Start(settings string) *Task {
	ctx, cancel := context.WithCancel(context.Background())
	state := &Task{
		cancel: cancel,
		ch:     make(chan string),
		ctx:    ctx,
		done:   0,
	}
	go runWithSerializedSettings(state, settings)
	return state
}

// Terminated is the event returned when a task is done
const Terminated = `{"key": "status.terminated", "value": {}}`

// WaitForNextEvent blocks until task generates the next event. Returns a
// valid pointer on success, a null pointer on failure.
func (task *Task) WaitForNextEvent() string {
	event, ok := <-task.ch
	if !ok {
		return Terminated
	}
	return event
}

// IsDone returns true if the task is done, false otherwise.
func (task *Task) IsDone() bool {
	return atomic.LoadInt64(&task.done) != 0
}

// Interrupt interrupts a running task.
func (task *Task) Interrupt() {
	task.cancel()
}
