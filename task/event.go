package task

import (
	"encoding/json"

	"github.com/measurement-kit/measurement-kit/geolookup"
)

// event is an event emitted by a task.
type event struct {
	// Key is the event key
	Key string `json:"key"`

	// Value is the event value
	Value interface{} `json:"value"`
}

// bugJsonDump is emitted when we cannot serialize a JSON.
var bugJsonDump = `{"key":"bug.json_dump","value":{"failure":"internal_error"}}`

// emitBugJsonDump emits a bug.json_dump event
func emitBugJsonDump(task *Task) {
	task.ch <- bugJsonDump
}

// emit emits an event with the given task
func emit(task *Task, event event) {
	data, err := json.Marshal(event)
	if err != nil {
		emitBugJsonDump(task)
		return
	}
	task.ch <- string(data)
}

// failureMeasurementSubmission is the failure.measurement_submission event body
type failureMeasurementSubmission struct {
	// Failure is the failure that occurred
	Failure string `json:"failure"`

	// Idx is the measurement index
	Idx int `json:"idx"`

	// Input is the measurement input
	Input string `json:"input"`

	// JSONStr is the serialize measurement
	JSONStr string `json:"json_str"`
}

// emitFailureMeasurementSubmission emits the failure.measurement_submission event
func emitFailureMeasurementSubmission(task *Task, err error, idx int, input, jsonStr string) {
	emit(task, event{Key: "failure.measurement_submission", Value: failureMeasurementSubmission{
		Failure: err.Error(),
		Idx: idx,
		Input: input,
		JSONStr: jsonStr,
	}});
}

// failureReportCreate is the failure.report_create event body
type failureReportCreate struct {
	// Failure is the error that occurred
	Failure string `json:"failure"`
}

// emitFailureReportCreate emits the failure.report_create event
func emitFailureReportCreate(task *Task, err error) {
	emit(task, event{Key: "failure.report_create", Value: failureStartup{
		Failure: err.Error(),
	}})
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

// logRecord is the log event body
type logRecord struct {
	// LogLevel is the log level
	LogLevel string `json:"log_level"`

	// Message is the log message
	Message string `json:"message"`
}

// emitLogWarning emits a warning message
func emitLogWarning(task *Task, message string) {
	emit(task, event{Key: "log", Value: logRecord{
		LogLevel: "WARNING",
		Message: message,
	}})
}

// measurement is the measurement event body
type measurement struct {
	// Idx is the measurement index
	Idx int `json:"idx"`

	// Input is the measurement input
	Input string `json:"input"`

	// JSONStr is the serialize measurement
	JSONStr string `json:"json_str"`
}

// emitMeasurement emits the measurement event
func emitMeasurement(task *Task, idx int, input, jsonStr string) {
	emit(task, event{Key: "measurement", Value: measurement{
		Idx: idx,
		Input: input,
		JSONStr: jsonStr,
	}});
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

// emitStatusGeoIPLookup emits the status.geoip_lookup event
func emitStatusGeoIPLookup(task *Task, result geolookup.Result) {
	emit(task, event{Key: "status.geoip_lookup", Value: result})
}

// statusMeasurementDone is the status.measurement_done event body
type statusMeasurementDone struct {
	// Idx is the measurement index
	Idx int `json:"idx"`

	// Input is the measurement input
	Input string `json:"input"`
}

// emitStatusMeasurementDone emits the status.measurement_done event
func emitStatusMeasurementDone(task *Task, idx int, input string) {
	emit(task, event{Key: "status.measurement_done", Value: statusMeasurementDone{
		Idx: idx,
		Input: input,
	}});
}

// statusMeasurementStart is the status.measurement_start event body
type statusMeasurementStart struct {
	// Idx is the measurement index
	Idx int `json:"idx"`

	// Input is the measurement input
	Input string `json:"input"`
}

// emitStatusMeasurementStart emits the status.measurement_start event
func emitStatusMeasurementStart(task *Task, idx int, input string) {
	emit(task, event{Key: "status.measurement_start", Value: statusMeasurementStart{
		Idx: idx,
		Input: input,
	}});
}

// statusMeasurementSubmission is the status.measurement_submission event body
type statusMeasurementSubmission struct {
	// Idx is the measurement index
	Idx int `json:"idx"`

	// Input is the measurement input
	Input string `json:"input"`
}

// emitStatusMeasurementSubmission emits the status.measurement_submission event
func emitStatusMeasurementSubmission(task *Task, idx int, input string) {
	emit(task, event{Key: "status.measurement_submission", Value: statusMeasurementSubmission{
		Idx: idx,
		Input: input,
	}});
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

// statusReportCreate is the status.report_create event body
type statusReportCreate struct {
	// ReportID is the report ID
	ReportID string `json:"report_id"`
}

// emitStatusReportCreate emits a status.report_create event
func emitStatusReportCreate(task *Task, ID string) {
	emit(task, event{Key: "status.report_create", Value: statusReportCreate{
		ReportID: ID}})
}

// emitStatusQueued emits a status.queued event
func emitStatusQueued(task *Task) {
	task.ch <- `{"key":"status.queued","value":{}}`
}

// emitStatusStarted emits a status.started event
func emitStatusStarted(task *Task) {
	task.ch <- `{"key":"status.started","value":{}}`
}
