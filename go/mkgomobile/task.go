package mkgomobile

import (
	"github.com/measurement-kit/measurement-kit/go/task"
)

type Task struct {
	task *task.Task
}

func (t *Task) Start(settings string) bool {
	if t.task != nil {
		return false
	}
	t.task = task.Start(settings)
	return true
}

func (t *Task) IsDone() bool {
	if t.task == nil {
		return true
	}
	return t.task.IsDone()
}

func (t *Task) WaitForNextEvent() string {
	if t.task == nil {
		return task.Terminated
	}
	return t.task.WaitForNextEvent()
}

func (t *Task) Interrupt() {
	if t.task != nil {
		t.task.Interrupt()
	}
}
