package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"
	"sync/atomic"

	"github.com/measurement-kit/measurement-kit/bouncer"
	"github.com/measurement-kit/measurement-kit/nettest"
)

// geolookup performs a geoip lookup
func geolookup(task *Task, nettest *nettest.Nettest) bool {
	err := nettest.GeoLookup()
	if err != nil {
		emit(task, event{Key: "failure.startup", Value: eventValue{
			Failure: err.Error(),
		}})
		return false
	}
	emit(task, event{Key: "status.progress", Value: eventValue{
		Percentage: 0.15,
		Message:    fmt.Sprintf("GeoInfo: %+v", nettest.GeoInfo),
	}})
	return true
}

// openreport opens a report
func openreport(task *Task, nettest *nettest.Nettest) bool {
	err := nettest.OpenReport()
	if err != nil {
		emit(task, event{Key: "failure.report_create", Value: eventValue{
			Failure: err.Error(),
		}})
		return false
	}
	emit(task, event{Key: "status.progress", Value: eventValue{
		Percentage: 0.2,
		Message:    fmt.Sprintf("Report: %+v", nettest.Report),
	}})
	return true
}

// runWithSettings runs a nettest with specific settings.
func runWithSettings(task *Task, settings settings, nettest *nettest.Nettest) statusEnd {
	if !geolookup(task, nettest) ||
		!openreport(task, nettest) {
		return statusEnd{}
	}
	defer nettest.CloseReport()
	for idx, input := range settings.Inputs {
		measurement := nettest.NewMeasurement()
		nettest.Measure(input, &measurement)
		data, err := json.Marshal(measurement)
		if err != nil {
			emit(task, event{Key: "bug.json_dump", Value: eventValue{
				Failure: err.Error(),
			}})
			return statusEnd{}
		}
		emit(task, event{Key: "measurement", Value: eventValue{
			Idx: int64(idx), JSONStr: string(data)}})
		err = nettest.SubmitMeasurement(&measurement)
		if err != nil {
			emit(task, event{Key: "failure.measurement_submission", Value: eventValue{
				Failure: err.Error(),
				JSONStr: string(data),
			}})
			continue
		}
	}
	emit(task, event{Key: "status.progress", Value: eventValue{
		Percentage: 1.0, Message: "Nettest complete"}})
	return statusEnd{}
}

// queryBouncer queries the OONI bouncer
func queryBouncer(task *Task, settings settings, nettest *nettest.Nettest) statusEnd {
	emitStatusStarted(task)
	if !settings.Options.NoBouncer {
		baseURL := "https://bouncer.ooni.io/"
		if settings.Options.BouncerBaseURL != "" {
			baseURL = settings.Options.BouncerBaseURL
		}
		URL, err := url.Parse(baseURL)
		if err != nil {
			emitFailureStartup(task, err)
			return statusEnd{Failure: "value_error"}
		}
		nettest.SelectedBouncer = &bouncer.Entry{
			Type:    "https",
			Address: URL.Host,
		}
		err = nettest.DiscoverAvailableCollectors()
		if err != nil {
			emitFailureStartup(task, err)
			if !settings.Options.IgnoreBouncerError {
				return statusEnd{Failure: "generic_error"}
			}
		}
		err = nettest.DiscoverAvailableTestHelpers()
		if err != nil {
			emitFailureStartup(task, err)
			if !settings.Options.IgnoreBouncerError {
				return statusEnd{Failure: "generic_error"}
			}
		}
		err = nettest.AutomaticallySelectCollector()
		if err != nil {
			emitFailureStartup(task, err)
			if !settings.Options.IgnoreBouncerError {
				return statusEnd{Failure: "generic_error"}
			}
		}
		emitStatusProgress(task, 0.1, "contacted bouncer")
	}
	// TODO(bassosimone): this function and the ones that it calls are
	// just basic stubs that need a second round of review
	return runWithSettings(task, settings, nettest)
}

// makeNettestAndRun makes a nettest and runs the task.
func makeNettestAndRun(task *Task, settings settings) statusEnd {
	if settings.Name == "PsiphonTunnel" {
		nettest, err := psiphontunnelNew(task, &settings)
		if err != nil {
			emitFailureStartup(task, err)
			return statusEnd{Failure: "value_error"}
		}
		return queryBouncer(task, settings, nettest)
	}
	emitFailureStartup(task, errors.New("Unknown task name"))
	return statusEnd{Failure: "value_error"}
}

// openLogFileAndRun opens the log file and runs the task.
func openLogFileAndRun(task *Task, settings settings) statusEnd {
	if settings.LogFilepath != "" {
		filep, err := os.OpenFile(
			settings.LogFilepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		// Note that the specification says that we're supposed to ignore
		// any error resulting in opening or writing the logfile
		if err != nil {
			task.logFile = filep
			defer filep.Close()
		}
	}
	return makeNettestAndRun(task, settings)
}

// parseAndRun parses serializedsettings and runs the task.
func parseAndRun(task *Task, serializedsettings string) statusEnd {
	var settings settings
	err := json.Unmarshal([]byte(serializedsettings), &settings)
	if err != nil {
		emitFailureStartup(task, err)
		return statusEnd{Failure: "value_error"}
	}
	return openLogFileAndRun(task, settings)
}

// semaphore prevents tests from running in parallel
var semaphore sync.Mutex

// taskMain runs a task with serialized settings
func taskMain(task *Task, settings string) {
	emitStatusQueued(task)
	semaphore.Lock() // blocked until my turn
	// TODO(bassosimone): measure the amount of data consumed
	emitStatusEnd(task, parseAndRun(task, settings))
	close(task.ch)
	atomic.StoreInt64(&task.done, 1)
	semaphore.Unlock() // allow another test to run
}
