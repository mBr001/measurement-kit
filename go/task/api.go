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

// New starts a task with the specified settings.
func New(settings string) *State {
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
func (task *State) WaitForNextEvent() string {
	const terminated = `{"key": "status.terminated", "value": {}}`
	event, ok := <-task.ch
	if !ok {
		return terminated
	}
	return event
}

// IsDone returns true if the task is done, false otherwise.
func (task *State) IsDone() bool {
	return atomic.LoadInt64(&task.done) != 0
}

// Interrupt interrupts a running task.
func (task *State) Interrupt() {
	task.cancel()
}
