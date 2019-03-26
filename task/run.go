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

// urlToAddress returns the Host of a given URL or an error.
func urlToAddress(URL string) (string, error) {
	exploded, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	return exploded.Host, nil
}

// jsonMarshalString is like json.Marshal but returns a string
func jsonMarshalString(input interface{}) (string, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// loop loops over the available input
func loop(task *Task, settings settings, nettest *nettest.Nettest) statusEnd {
	// TODO(bassosimone): read input files and randomize
	// TODO(bassosimone): parallelism
	// TODO(bassosimone): max_runtime
	for idx, input := range settings.Inputs {
		emitStatusMeasurementStart(task, idx, input)
		// TODO(bassosimone): fill annotations etc
		measurement := nettest.NewMeasurement()
		nettest.Measure(input, &measurement)
		jsonStr, err := jsonMarshalString(measurement)
		if err != nil {
			emitBugJsonDump(task)
			continue
		}
		emitMeasurement(task, idx, input, jsonStr)
		if !settings.Options.NoCollector {
			err = nettest.SubmitMeasurement(&measurement)
			if err != nil {
				emitFailureMeasurementSubmission(task, err, idx, input, jsonStr)
			} else {
				emitStatusMeasurementSubmission(task, idx, input)
			}
		}
		emitStatusMeasurementDone(task, idx, input)
	}
	emitStatusProgress(task, 0.9, "ending the test")
	return statusEnd{}
}

// openReport opens the report
func openReport(task *Task, settings settings, nettest *nettest.Nettest) statusEnd {
	// TODO(bassosimone): add support for writing report to file
	nettest.SoftwareName = settings.Options.SoftwareName
	nettest.SoftwareVersion = settings.Options.SoftwareVersion
	if !settings.Options.NoCollector {
		if settings.Options.CollectorBaseURL != "" {
			address, err := urlToAddress(settings.Options.CollectorBaseURL)
			if err != nil {
				emitFailureStartup(task, err)
				return statusEnd{Failure: "generic_error"}
			}
			nettest.SelectedCollector = &bouncer.Entry{
				Type: "https",
				Address: address,
			}
		}
		err := nettest.OpenReport()
		if err != nil {
			emitFailureReportCreate(task, err)
			if !settings.Options.IgnoreOpenReportError {
				return statusEnd{Failure: "generic_error"}
			}
		} else {
			emitStatusReportCreate(task, nettest.Report.ID)
			defer nettest.CloseReport()
		}
	}
	emitStatusProgress(task, 0.4, "open report")
	return loop(task, settings, nettest)
}

// geoLookup performs a geolookup of the probe
func geoLookup(task *Task, settings settings, nettest *nettest.Nettest) statusEnd {
	nettest.GeoInfo.ProbeIP = "127.0.0.1"
	nettest.GeoInfo.ProbeASN = "AS0"
	nettest.GeoInfo.ProbeCC = "ZZ"
	nettest.GeoInfo.ProbeNetworkName = ""
	nettest.ASNDatabasePath = settings.Options.GeoIPASNPath
	nettest.CountryDatabasePath = settings.Options.GeoIPCountryPath
	err := nettest.GeoLookup()
	if err != nil {
		emitLogWarning(task, fmt.Sprintf("nettest.GeoLookup: %s", err.Error()))
		// FALLTHROUGH
	}
	emitStatusGeoIPLookup(task, nettest.GeoInfo)
	emitStatusProgress(task, 0.2, "geoip lookup")
	return openReport(task, settings, nettest)
}

// queryBouncer queries the OONI bouncer
func queryBouncer(task *Task, settings settings, nettest *nettest.Nettest) statusEnd {
	emitStatusStarted(task)
	if !settings.Options.NoBouncer {
		baseURL := "https://bouncer.ooni.io/"
		if settings.Options.BouncerBaseURL != "" {
			baseURL = settings.Options.BouncerBaseURL
		}
		address, err := urlToAddress(baseURL)
		if err != nil {
			emitFailureStartup(task, err)
			return statusEnd{Failure: "value_error"}
		}
		nettest.SelectedBouncer = &bouncer.Entry{
			Type:    "https",
			Address: address,
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
	return geoLookup(task, settings, nettest)
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
