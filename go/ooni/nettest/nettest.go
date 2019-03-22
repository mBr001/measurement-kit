// Package nettest contains generic code for running nettests.
//
// You typically create a Nettest instance by invoking the factory
// function exposed by a specific nettest package, and the use this
// generic code to manage the test lifecycle. For example
//
//     nettest := psiphontunnel.NewNettest(context.Background(), config)
//     defer nettest.Close()
//     err := nettest.DiscoverAvailableCollectors()
//     if err != nil {
//       log.Fatal(err)
//     }
//     err = nettest.SelectCollector()
//     if err != nil {
//       log.Fatal(err)
//     }
//     err = nettest.GeoLookup()
//     if err != nil {
//       log.Fatal(err)
//     }
//     err = nettest.OpenReport()
//     if err != nil {
//       log.Fatal(err)
//     }
//     measurement := nettest.Measure("")
//     measurementID, err := nettest.Submit(measurement)
//     if err != nil {
//       log.Fatal(err)
//     }
//     log.Printf("measurementID: %+v\n", measurementID)
//
// This API is such that every small operation that a test must perform
// is a separate operation. This allows you to handle errors and results
// of each separate operation in the way you find most convenient.
package nettest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/measurement-kit/measurement-kit/go/ooni/bouncer"
	"github.com/measurement-kit/measurement-kit/go/ooni/collector"
	"github.com/measurement-kit/measurement-kit/go/ooni/geolookup"
	"github.com/measurement-kit/measurement-kit/go/ooni/measurement"
)

// DateFormat is the format used by OONI for dates inside reports.
const DateFormat = "2006-01-02 15:04:05"

// Config contains the generic nettest configuration set by the
// application that wants to run the nettest.
type Config struct {
	// ASNDatabasePath contains the ASN DB path.
	ASNDatabasePath string

	// BouncerBaseURL contains the bouncer base URL.
	BouncerBaseURL string

	// SoftwareName contains the software name.
	SoftwareName string

	// SoftwareVersion contains the software version.
	SoftwareVersion string
}

// Func is the function that implements a nettest.
type Func = func(string) interface{}

// Nettest is a nettest.
type Nettest struct {
	// Ctx is the context for running the nettest.
	Ctx context.Context

	// Config is the user supplied configuration.
	Config Config

	// TestName is the name of the test.
	TestName string

	// TestVersion is the version of the test.
	TestVersion string

	// Func is the function that actually implements the test.
	Func Func

	// TestStartTime is the time when the test started. This is set
	// by the New function in this package.
	TestStartTime string

	// AvailableCollectors contains the available collectors. This field
	// is filled by DiscoverAvailableCollectors.
	AvailableCollectors []bouncer.Entry

	// SelectedCollector is the selected collector. This field is filled
	// by the SelectCollector function.
	SelectedCollector bouncer.Entry

	// GeoLookupInfo contains geolookup info. This field is filled by
	// the GeoLookup function.
	GeoLookupInfo geolookup.Result

	// Report is the report bound to this nettest. This field is initialized
	// by the OpenReport function and closed by the Close function.
	Report collector.Report
}

// DiscoverAvailableCollectors discovers the available collectors.
func (nettest *Nettest) DiscoverAvailableCollectors() error {
	collectors, err := bouncer.GetCollectors(nettest.Ctx, bouncer.Config{
		BaseURL: nettest.Config.BouncerBaseURL,
	})
	if err != nil {
		return err
	}
	nettest.AvailableCollectors = collectors
	return nil
}

// SelectCollector selects a collector from the available collectors.
func (nettest *Nettest) SelectCollector() error {
	for _, collector := range nettest.AvailableCollectors {
		if collector.Type == "https" {
			nettest.SelectedCollector = collector
			return nil
		}
	}
	return errors.New("No suitable collectors found")
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
		ASNDatabasePath: nettest.Config.ASNDatabasePath,
	})
	if err != nil {
		return err
	}
	nettest.GeoLookupInfo = info
	return nil
}

// probeASN is a convenience method for getting an always valid probe ASN.
func (nettest *Nettest) probeASN() string {
	if nettest.GeoLookupInfo.ProbeASN != "" {
		return nettest.GeoLookupInfo.ProbeASN
	}
	return "AS0"
}

// probeCC is like probeASN but for the country code (CC).
func (nettest *Nettest) probeCC() string {
	if nettest.GeoLookupInfo.ProbeCC != "" {
		return nettest.GeoLookupInfo.ProbeCC
	}
	return "ZZ"
}

// OpenReport opens a new report for the nettest.
func (nettest *Nettest) OpenReport() error {
	if nettest.Report.ID != "" {
		return nil // idempotent semantics is nice
	}
	report, err := collector.Open(nettest.Ctx, collector.Config{
		BaseURL: nettest.collectorBaseURL(),
	}, collector.ReportTemplate{
		ProbeASN:        nettest.probeASN(),
		ProbeCC:         nettest.probeCC(),
		SoftwareName:    nettest.Config.SoftwareName,
		SoftwareVersion: nettest.Config.SoftwareVersion,
		TestName:        nettest.TestName,
		TestVersion:     nettest.TestVersion,
	})
	if err != nil {
		return err
	}
	nettest.Report = report
	return nil
}

// Measure runs a nettest measurement with the provided input and returns the
// measurement object. Pass an empty string for input-less nettests. It is
// safe to call this method from different goroutines concurrently.
func (nettest *Nettest) Measure(input string) measurement.Measurement {
	measurement := nettest.NewMeasurement()
	measurement.Input = input
	t0 := time.Now()
	measurement.TestKeys = nettest.Func(input)
	elapsed := float64(time.Now().Sub(t0)) / float64(time.Second)
	measurement.TestRuntime = elapsed
	return measurement
}

// NewMeasurement returns a new measurement. The fields that the user should
// initialize are Inputs, TestKeys, and TestRuntime. All the other fields are
// already initialized by NewMeasurement.
func (nettest *Nettest) NewMeasurement() measurement.Measurement {
	return measurement.Measurement{
		DataFormatVersion:    "0.2.0",
		MeasurementStartTime: time.Now().UTC().Format(DateFormat),
		ProbeASN:             nettest.probeASN(),
		ProbeCC:              nettest.probeCC(),
		ReportID:             nettest.Report.ID,
		SoftwareName:         nettest.Config.SoftwareName,
		SoftwareVersion:      nettest.Config.SoftwareVersion,
		TestName:             nettest.TestName,
		TestStartTime:        nettest.TestStartTime,
		TestVersion:          nettest.TestVersion,
	}
}

// Submit submits a measurement to the discovered collector. Returns the
// measurementID on success and an error on failure. It is safe to call this
// method from different goroutines concurrently.
func (nettest *Nettest) Submit(measurement measurement.Measurement) (string, error) {
	measurementID, err := nettest.Report.Update(nettest.Ctx, measurement)
	if err != nil {
		return "", err
	}
	return measurementID, nil
}

// Close makes sure that, if we opened a report, we close it.
func (nettest *Nettest) Close() error {
	if nettest.Report.ID != "" {
		return nettest.Report.Close(nettest.Ctx)
	}
	return nil
}
