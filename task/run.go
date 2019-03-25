package task

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/measurement-kit/measurement-kit/nettest"
)

// discovercollectors discovers the available collectors
func discovercollectors(task *Task, nettest *nettest.Nettest) bool {
	err := nettest.DiscoverAvailableCollectors()
	if err != nil {
		emit(task, event{Key: "failure.startup", Value: eventValue{
			Failure: err.Error(),
		}})
		return false
	}
	emit(task, event{Key: "status.progress", Value: eventValue{
		Percentage: 0.05,
		Message: fmt.Sprintf(
			"AvailableCollectors: %+v", nettest.AvailableCollectors),
	}})
	return true
}

// selectcollector selects a collector
func selectcollector(task *Task, nettest *nettest.Nettest) bool {
	err := nettest.AutomaticallySelectCollector()
	if err != nil {
		emit(task, event{Key: "failure.startup", Value: eventValue{
			Failure: err.Error(),
		}})
		return false
	}
	emit(task, event{Key: "status.progress", Value: eventValue{
		Percentage: 0.1,
		Message:    fmt.Sprintf("SelectedCollector: %+v", nettest.SelectedCollector),
	}})
	return true
}

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

// makenettest creates a new nettest or returns nil
func makenettest(task *Task, settings *settings) *nettest.Nettest {
	if settings.Name == "psiphontunnel" {
		nettest, err := psiphontunnelNew(task, settings)
		if err != nil {
			emit(task, event{Key: "failure.startup", Value: eventValue{
				Failure: err.Error(),
			}})
			return nil
		}
		return nettest
	}
	emit(task, event{Key: "failure.startup", Value: eventValue{
		Failure: "Unknown nettest name",
	}})
	return nil
}

// runWithSettings runs a nettest with settings.
func runWithSettings(task *Task, settings settings) {
	nettest := makenettest(task, &settings)
	if nettest == nil {
		return
	}
	if !discovercollectors(task, nettest) ||
		!selectcollector(task, nettest) ||
		!geolookup(task, nettest) ||
		!openreport(task, nettest) {
		return
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
			return
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
}

// runWithSerializedSettings runs a task with serialized settings
func runWithSerializedSettings(task *Task, serializedsettings string) {
	defer close(task.ch)
	defer func() {
		atomic.StoreInt64(&task.done, 1)
	}()
	var settings settings
	err := json.Unmarshal([]byte(serializedsettings), &settings)
	if err != nil {
		emit(task, event{Key: "failure.startup", Value: eventValue{
			Failure: err.Error(),
		}})
		return
	}
	runWithSettings(task, settings)
}
