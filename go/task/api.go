// Package task implements the Measurement Kit task API
package task

import (
	"context"
	"sync/atomic"
)

// State is a task run by Measurement Kit
type State struct {
	cancel context.CancelFunc
	ch     chan string
	ctx    context.Context
	done   int64
}

// Start starts a task with the specified settings. Returns a valid
// pointer if settings are good, a null pointer otherwise.
func Start(settings string) *State {
	ctx, cancel := context.WithCancel(context.Background())
	state := &State{
		cancel: cancel,
		ch:     make(chan string),
		ctx:    ctx,
		done:   0,
	}
	go runWithSerializedSettings(state, settings)
	return state
}

// WaitForNextEvent blocks until task generates the next event. Returns a
// valid pointer on success, a null pointer on failure.
func WaitForNextEvent(task *State) string {
	const terminated = `{"key": "status.terminated", "value": {}}`
	if task == nil {
		return terminated
	}
	event, ok := <-task.ch
	if !ok {
		return terminated
	}
	return event
}

// IsDone returns true if the task is done, false otherwise.
func IsDone(task *State) bool {
	if task == nil {
		return true
	}
	return atomic.LoadInt64(&task.done) != 0
}

// Interrupt interrupts a running task.
func Interrupt(task *State) {
	if task != nil {
		task.cancel()
	}
}
