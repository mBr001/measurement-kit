// Package mkgomobile contains the Measurement Kit API exposed on mobile.
package mkgomobile

import (
	"github.com/measurement-kit/measurement-kit/go/task"
)

// Task is a task running some operations. The usage is as follows:
//
//     task := mkgomobile.StartTask(settings)
//     for !task.IsDone() {
//      event := task.WaitForNextEvent()
//      // process the event
//     }
//
// See https://github.com/measurement-kit/measurement-kit/tree/master/include/measurement_kit
// for the documentation regarding settings and events.
type Task struct {
	task *task.Task
}

// StartTask starts a new task.
func StartTask(settings string) *Task {
	return &Task{task: task.Start(settings)}
}

// IsDone returns whether a task is done.
func (t *Task) IsDone() bool {
	return t.task.IsDone()
}

// WaitForNextEvent blocks until the task emits the next event.
func (t *Task) WaitForNextEvent() string {
	return t.task.WaitForNextEvent()
}

// Interrupt interrupts a running task.
func (t *Task) Interrupt() {
	t.task.Interrupt()
}
