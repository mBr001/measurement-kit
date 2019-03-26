// Package task implements the Measurement Kit task API
package task

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
)

// Task is a task run by Measurement Kit.
type Task struct {
	// cancel cancels the task context.
	cancel context.CancelFunc

	// ch is the channel where events are emitted.
	ch chan string

	// ctx is the task context.
	ctx context.Context

	// done indicates whether the task is done.
	done int64

	// downloadedKB measures the downloaded KBs.
	downloadedKB float64

	// logFile is the log file
	logFile *os.File

	// logMutex protects logFile
	logMutex sync.Mutex

	// uploadedKB measures the uploaded KBs.
	uploadedKB float64
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
	go taskMain(state, settings)
	return state
}

// Terminated is the event returned when a task is done
const Terminated = `{"key":"status.terminated","value":{}}`

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
