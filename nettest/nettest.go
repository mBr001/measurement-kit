// Package nettest contains code for running nettests.
//
// This API is such that every small operation that a test must perform
// is a separate operation. This allows you to handle errors and results
// of each separate operation in the way you find most convenient.
//
// An nettest is a nettest that is not implemented as part of
// this library but is implemented elsewhere.
//
// Creating a nettest
//
// To create a nettest just instantiate it:
//
//     var nettest nettest.Nettest
//
// You must fill the following fields:
//
// - nettest.Ctx with a context for the nettest
//
// - nettest.TestName with the name of the nettest
//
// - nettest.TestVersion with the nettest version
//
// - nettest.SoftwareName with the app name
//
// - nettest.SoftwareVersion with the app version
//
// - nettest.TestStartTime with the UTC test start time formatted according
// the format expected by OONI (you can use nettest.FormatTimeNowUTC
// to initialize this field with the current UTC time, or nettest.DateFormat
// to format another time according to the proper format -- just remember
// that you must use the UTC time here)
//
// - nettest.Measure if the nettest is written in Go, otherwise, if the
// nettest is still in C++, you'll need to call C++ code to get back
// the test keys and properly finish initializing a measurement.
//
// For example
//
//     nettest.Ctx = context.Background()
//     nettest.TestName = "nettest"
//     nettest.TestVersion = "0.0.1"
//     nettest.SoftwareName = "example"
//     nettest.SoftwareVersion = "0.0.1"
//     nettest.TestStartTime = nettest.FormatTimeNowUTC()
//     nettest.Measure = func(input string, m *measurement.Measurement) {
//       // perform measurement and initialize m
//     }
//
// You may also want to fill other fields or use nettest specific
// functionality for automatically filling them.
//
// Selecting a bouncer
//
// The bouncer is used to discover collectors and test helpers. If you
// don't have a specific bouncer in mind, just skip this step and we'll
// use the default bouncer. Otherwise, do something like:
//
//     nettest.SelectedBouncer = &bouncer.Entry{
//       Type: "https",
//       Address: "bouncer.example.com",
//     }
//
// Where Address must be the bouncer FQDN with optional port.
//
// Selecting a collector
//
// You must select a collector to submit nettest measurements to
// the OONI collector. We recommend you to automatically discover
// a collector. Otherwise just initialize nettest.Collector.
//
// To automatically discover a collector do the following:
//
//     err := nettest.DiscoverAvailableCollectors()
//     if err != nil {
//       // handle error
//     }
//
// This will populate the nettest.AvailableCollector field. At this
// point you can either manually select a collector or do:
//
//     err = nettest.AutomaticallySelectCollector()
//     if err != nil {
//       // handle error
//     }
//
// This will populate the nettest.Collector field. If you just run
// nettest.AutomaticallySelectCollector without running also
// nettest.DiscoverAvailableCollectors, the available collectors
// will be automatically discovered for you.
//
// Collectors will be automatically discovered using the OONI
// bouncer. You can set the nettest.SelectedBouncer field to
// force the code using a specific bouncer.
//
// Selecting test helpers
//
// If your test needs test helpers, you should discover the available
// test helpers using:
//
//     err = nettest.DiscoverAvailableTestHelpers()
//     if err != nil {
//       // handle error
//     }
//
// This will fill the nettest.AvailableTestHelpers field. If you know
// about test helpers with other means, otherwise, you can just skip
// nettest.DiscoverAvailableTestHelpers and just initialize the
// nettest.AvailableTestHelpers field.
//
// Test helpers will be automatically discovered using the OONI
// bouncer. You can set the nettest.SelectedBouncer field to
// force the code using a specific bouncer.
//
// Geolocate the probe
//
// Geolocating a probe means discover its IP, CC (country code),
// ASN (autonomous system number), and network name (i.e. the
// commercial name bound to the ASN).
//
// If you already know this values, just initialize them; e.g.:
//
//     nettest.GeoInfo.ProbeIP = "93.147.252.33"
//     nettest.GeoInfo.ProbeCC = "IT"
//     nettest.GeoInfo.ProbeASN = "AS30722"
//     nettest.GeoInfo.ProbeIP = "Vodafone Italia"
//
// Otherwise, you need to initialize the CountryDatabasePath and
// the ASNDatabasePath fields to point to valid and current MaxMind
// MMDB databases; e.g.,
//
//     nettest.CountryDatabasePath = "country.mmdb"
//     nettest.ASNDatabasePath = "asn.mmdb"
//
// Then run:
//
//     err = nettest.GeoLookup()
//
// This will fill the nettest.GeoInfo field.
//
// Opening a report
//
// You need a nettest.SelectedCollector to do that. Then run:
//
//     err = nettest.OpenReport()
//     if err != nil {
//       // handle error
//     }
//     defer nettest.CloseReport()
//
// This will initialize the nettest.Report.ID field. If this field
// is already initialized, this step will fail. This means, among
// other things, that you can only open a report once.
//
// Creating a new measurement
//
// You are now ready to perform measurements. Ask the nettest to
// create for you a measurement with:
//
//     measurement := nettest.NewMeasurement()
//
// This will initialize all measurement fields except:
//
// - measurement.TestKeys, which should contains a JSON serializable
// interface{} containing the nettest specific results
//
// - measurement.MeasurementRuntime, which should contain the measurement
// runtime in seconds as a floating point
//
// - measurement.Input, which should only be initialized if your
// nettest requires input
//
// If nettest.Measure is initialized, it will do that for you:
//
//     nettest.Measure(input, &measurement)
//
// where input is an empty string is the nettest does not take any
// input. Otherwise, you'll need to call C++ code to get the test
// keys and initialize Runtime and Input yourself. Either way, when
// you're done, you can submit the measurement.
//
// Note that, by default, the ProbeIP in the measurement will be set
// to "127.0.0.1". If you want to submit the real probe IP, you'll
// need to override measurement.ProbeIP with nettest.GeoInfo.ProbeIP.
//
// Submitting a measurement
//
// To submit a measurement, you need a nettest.SelectedCollector and
// the report should have been openned. Then run:
//
//     err := nettest.SubmitMeasurement(&measurement)
//     if err != nil {
//       // handle error
//     }
//
// If successful, this will set the measurement.OOID field, which
// may be empty if the collector does not support if. If this field
// isn't empty, later you can use this OOID to get the measurement
// from the OONI API.
package nettest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/measurement-kit/measurement-kit/bouncer"
	"github.com/measurement-kit/measurement-kit/collector"
	"github.com/measurement-kit/measurement-kit/geolookup"
	"github.com/measurement-kit/measurement-kit/measurement"
)

// DateFormat is the format used by OONI for dates inside reports.
const DateFormat = "2006-01-02 15:04:05"

// FormatTimeNowUTC formats the current time in UTC using the OONI format.
func FormatTimeNowUTC() string {
	return time.Now().UTC().Format(DateFormat)
}

// MeasureFunc is the function running a measurement. Pass an empty string
// if the nettest does not take input. Remember to initialize the fields
// of measurement that are not initialized by NewMeasurement (see above for
// a complete list of such fields).
type MeasureFunc = func(input string, measurement *measurement.Measurement)

// Nettest is a nettest.
type Nettest struct {
	// Ctx is the context for the nettest.
	Ctx context.Context

	// TestName is the test name.
	TestName string

	// TestVersion is the test version.
	TestVersion string

	// SoftwareName contains the software name.
	SoftwareName string

	// SoftwareVersion contains the software version.
	SoftwareVersion string

	// TestStartTime is the UTC time when the test started.
	TestStartTime string

	// Measure runs the measurement.
	Measure MeasureFunc

	// SelectedBouncer is the selected bouncer.
	SelectedBouncer *bouncer.Entry

	// SelectedCollector is the selected collector.
	SelectedCollector *bouncer.Entry

	// AvailableCollectors contains all the available collectors.
	AvailableCollectors []bouncer.Entry

	// AvailableCollectors contains all the available test helpers.
	AvailableTestHelpers []bouncer.Entry

	// CountryDatabasePath contains the country MMDB database path.
	CountryDatabasePath string

	// ASNDatabasePath contains the ASN MMDB database path.
	ASNDatabasePath string

	// GeoInfo contains the geolookup result.
	GeoInfo geolookup.Result

	// Report is the report bound to this nettest.
	Report collector.Report
}

// bouncerBaseURL returns the bouncer base URL.
func (nettest *Nettest) bouncerBaseURL() string {
	if nettest.SelectedBouncer != nil {
		return fmt.Sprintf("https://%s/", nettest.SelectedBouncer.Address)
	}
	return "https://bouncer.ooni.io"
}

// DiscoverAvailableCollectors discovers the available collectors.
func (nettest *Nettest) DiscoverAvailableCollectors() error {
	collectors, err := bouncer.GetCollectors(nettest.Ctx, bouncer.Config{
		BaseURL: nettest.bouncerBaseURL(),
	})
	if err != nil {
		return err
	}
	nettest.AvailableCollectors = collectors
	return nil
}

// AutomaticallySelectCollector automatically selects a collector.
func (nettest *Nettest) AutomaticallySelectCollector() error {
	if len(nettest.AvailableCollectors) <= 0 {
		err := nettest.DiscoverAvailableCollectors()
		if err != nil {
			return err
		}
	}
	for _, collector := range nettest.AvailableCollectors {
		if collector.Type == "https" {
			nettest.SelectedCollector = &collector
			return nil
		}
	}
	return errors.New("No suitable collectors found")
}

// DiscoverAvailableTestHelpers discovers the available test helpers.
func (nettest *Nettest) DiscoverAvailableTestHelpers() error {
	testHelpers, err := bouncer.GetTestHelpers(nettest.Ctx, bouncer.Config{
		BaseURL: nettest.bouncerBaseURL(),
	})
	if err != nil {
		return err
	}
	nettest.AvailableTestHelpers = testHelpers
	return nil
}

// collectorBaseURL is an internal convenience method to compute
// the collector's base URL from the selected collector.
func (nettest *Nettest) collectorBaseURL() string {
	if nettest.SelectedCollector.Address != "" {
		return fmt.Sprintf("https://%s/", nettest.SelectedCollector.Address)
	}
	return "https://a.collector.ooni.io/"
}

// GeoLookup performs the geolookup (probe_ip, probe_asn, etc.)
func (nettest *Nettest) GeoLookup() error {
	info, err := geolookup.Perform(nettest.Ctx, geolookup.Config{
		ASNDatabasePath: nettest.ASNDatabasePath,
	})
	if err != nil {
		return err
	}
	nettest.GeoInfo = info
	return nil
}

// probeASN is a convenience method for getting an always valid probe ASN.
func (nettest *Nettest) probeASN() string {
	if nettest.GeoInfo.ProbeASN != "" {
		return nettest.GeoInfo.ProbeASN
	}
	return "AS0"
}

// probeCC is like probeASN but for the country code (CC).
func (nettest *Nettest) probeCC() string {
	if nettest.GeoInfo.ProbeCC != "" {
		return nettest.GeoInfo.ProbeCC
	}
	return "ZZ"
}

// OpenReport opens a new report for the nettest.
func (nettest *Nettest) OpenReport() error {
	if nettest.Report.ID != "" {
		return errors.New("Report is already open")
	}
	report, err := collector.Open(nettest.Ctx, collector.Config{
		BaseURL: nettest.collectorBaseURL(),
	}, collector.ReportTemplate{
		ProbeASN:        nettest.probeASN(),
		ProbeCC:         nettest.probeCC(),
		SoftwareName:    nettest.SoftwareName,
		SoftwareVersion: nettest.SoftwareVersion,
		TestName:        nettest.TestName,
		TestVersion:     nettest.TestVersion,
	})
	if err != nil {
		return err
	}
	nettest.Report = report
	return nil
}

// NewMeasurement returns a new measurement for this nettest. You should
// fill fields that are not initialized; see above for a description
// of what fields WILL NOT be initialized.
func (nettest *Nettest) NewMeasurement() measurement.Measurement {
	return measurement.Measurement{
		DataFormatVersion:    "0.2.0",
		MeasurementStartTime: time.Now().UTC().Format(DateFormat),
		ProbeIP:              "127.0.0.1",
		ProbeASN:             nettest.probeASN(),
		ProbeCC:              nettest.probeCC(),
		ReportID:             nettest.Report.ID,
		SoftwareName:         nettest.SoftwareName,
		SoftwareVersion:      nettest.SoftwareVersion,
		TestName:             nettest.TestName,
		TestStartTime:        nettest.TestStartTime,
		TestVersion:          nettest.TestVersion,
	}
}

// SubmitMeasurement submits a measurement to the selected collector. It is
// safe to call this function from different goroutines concurrently as long
// as the measurement is not shared by the goroutines.
func (nettest *Nettest) SubmitMeasurement(measurement *measurement.Measurement) error {
	measurementID, err := nettest.Report.Update(nettest.Ctx, *measurement)
	if err != nil {
		return err
	}
	measurement.OOID = measurementID
	return nil
}

// CloseReport closes an open report.
func (nettest *Nettest) CloseReport() error {
	return nettest.Report.Close(nettest.Ctx)
}
