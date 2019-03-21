package task

import (
	"encoding/json"
)

// eventValue is the value field of an event
type eventValue struct {
	Failure    string  `json:"failure,omitempty"`
	Idx        int64   `json:"idx,omitempty"`
	JSONStr    string  `json:"json_str,omitempty"`
	Message    string  `json:"message,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}

// event is an event emitted by a task.
type event struct {
	Key   string     `json:"key"`
	Value eventValue `json:"value"`
}

// emit emits an event with the given task
func emit(task *Task, event event) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	if task == nil {
		return
	}
	task.ch <- string(data)
}
