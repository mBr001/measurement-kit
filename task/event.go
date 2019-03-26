package task

import (
	"encoding/json"
)

// eventValue is the value field of an event.
//
// Implementation note: because of the way in which JSON serialization
// works in golang, we cannot use omitempty here because there are cases
// like e.g. `idx` where a zero value has meaning. So, use the opposite
// approach where we always emit all keys and the client, which should
// know what keys have meaning for specific events, will correctly only
// consider the fields that make sense in that context.
type eventValue struct {
	// Failure is set when there is a failure
	Failure string `json:"failure"`

	// Idx is the index of the nettest
	Idx int64 `json:"idx"`

	// JSONStr is the serialization of a measurement
	JSONStr string `json:"json_str"`

	// Message is a message bound to the event
	Message string `json:"message"`

	// Percentage indicates the progress
	Percentage float64 `json:"percentage"`
}

// event is an event emitted by a task.
type event struct {
	// Key is the event key
	Key string `json:"key"`

	// Value is the event value
	Value interface{} `json:"value"`
}

// emit emits an event with the given task
func emit(task *Task, event event) {
	data, err := json.Marshal(event)
	if err != nil {
		// TODO(bassosimone): the following error could be more informative, yet
		// it's better to have some non-super informative error rather than
		// completely suppressing the event emission.
		data = []byte(`{"key":"bug.json_dump","value":{"failure": "internal_error"}}`)
	}
	task.ch <- string(data)
}

// failureStartup is the failure.startup event body
type failureStartup struct {
	// Failure is the error that occurred
	Failure string `json:"failure"`
}

// emitFailureStartup emits the failure.startup event
func emitFailureStartup(task *Task, err error) {
	emit(task, event{Key: "failure.startup", Value: failureStartup{
		Failure: err.Error(),
	}})
}

// statusEnd is the status.end event body
type statusEnd struct {
	// DownloadedKB informs you about the downloaded KBs
	DownloadedKB float64 `json:"downloaded_kb"`

	// Failure is set when there is a failure
	Failure string `json:"failure"`

	// UploadedKB informs you about the uploaded KBs
	UploadedKB float64 `json:"uploaded_kb"`
}

// emitStatusEnd emits the status.end event
func emitStatusEnd(task *Task, statusEnd statusEnd) {
	emit(task, event{Key: "status.end", Value: statusEnd})
}

// statusProgress is the status.progress event body
type statusProgress struct {
	// Message is a message bound to the event
	Message string `json:"message"`

	// Percentage indicates the progress
	Percentage float64 `json:"percentage"`
}

// emitStatusProgress emits a status.progress event
func emitStatusProgress(task *Task, percentage float64, message string) {
	emit(task, event{Key: "status.progress", Value: statusProgress{
		Message: message, Percentage: percentage}})
}

// emitStatusQueued emits a status.queued event
func emitStatusQueued(task *Task) {
	task.ch <- `{"key":"status.queued","value":{}}`
}

// emitStatusStarted emits a status.started event
func emitStatusStarted(task *Task) {
	task.ch <- `{"key":"status.started","value":{}}`
}
