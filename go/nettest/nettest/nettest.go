// Package nettest contains the generic code to run a nettest.
package nettest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/measurement-kit/measurement-kit/go/bouncer"
	"github.com/measurement-kit/measurement-kit/go/collector"
	"github.com/measurement-kit/measurement-kit/go/geolookupper"
)

// Config contains the generic nettest configuration set by the
// application that wants to run the nettest.
type Config struct {
	// ASNDBPath contains the ASN DB path.
	ASNDBPath string

	// BouncerBaseURL contains the bouncer base URL.
	BouncerBaseURL string

	// SoftwareName contains the software name.
	SoftwareName string

	// SoftwareVersion contains the software version.
	SoftwareVersion string
}

// dateformat is the format used by OONI for dates inside reports.
const dateformat = "2006-01-02 15:04:05"

// Nettest is a generic nettest. You should create a specific nettest
// by using the specific-nettest-package New method. In turn such
// nettest-specific packages should get a partially initialized Nettest
// by calling the NewPartial function on this package.
type Nettest struct {
	// Ctx is the context for running the nettest. This must be set
	// by the constructor of the specific nettest package.
	Ctx context.Context

	// Config is the user supplied configuration. Also this field must
	// be set by the constructor of the specific nettest package.
	Config Config

	// TestName is the name of the test. Also this field must
	// be set by the constructor of the specific nettest package.
	TestName string

	// TestVersion is the version of the test. Also this field must
	// be set by the constructor of the specific nettest package.
	TestVersion string

	// RunFunc is the function that actually implements the test. Also this field
	// must be set by the constructor of the specific nettest package.
	RunFunc func(string)interface{}

	// TestStartTime is the time when the test started. This is set
	// by the NewPartial function in this package.
	TestStartTime string

	// AvailableCollectors contains the available collectors. This field
	// is filled by DiscoverAvailableCollectors.
	AvailableCollectors []bouncer.Entry

	// SelectedCollector is the selected collector. This field is filled
	// by the SelectCollector function.
	SelectedCollector bouncer.Entry

	// GeoLookupInfo contains geolookup info. This field is filled by
	// the GeoLookup function.
	GeoLookupInfo geolookupper.Result

	// Report is the report bound to this nettest. This field is initialized
	// by the OpenReport function and closed by the Close function.
	Report collector.Report
}

// NewPartial returns a partially initialized nettest. More initialization
// is required to actually used it, as already mentioned.
func NewPartial() *Nettest {
	return &Nettest{
		TestStartTime: time.Now().UTC().Format(dateformat),
	}
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
	} else {
		return "https://a.collector.ooni.io/"
	}
}

// GeoLookup performs the geolookup (probe_ip, probe_asn, etc.)
func (nettest *Nettest) GeoLookup() error {
	info, err := geolookupper.Lookup(nettest.Ctx, geolookupper.Config{
		ASNDBPath: nettest.Config.ASNDBPath,
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
	} else {
		return "AS0"
	}
}

// probeCC is like probeASN but for the country code (CC).
func (nettest *Nettest) probeCC() string {
	if nettest.GeoLookupInfo.ProbeCC != "" {
		return nettest.GeoLookupInfo.ProbeCC
	} else {
		return "ZZ"
	}
}

// OpenReport opens a new report for the nettest.
func (nettest *Nettest) OpenReport() error {
	if nettest.Report.ID != "" {
		return nil // idempotent semantics is nice
	}
	report, err := collector.Open(nettest.Ctx, collector.Config{
		BaseURL: nettest.collectorBaseURL(),
	}, collector.Template{
		ProbeASN: nettest.probeASN(),
		ProbeCC: nettest.probeCC(),
		SoftwareName: nettest.Config.SoftwareName,
		SoftwareVersion: nettest.Config.SoftwareVersion,
		TestName: nettest.TestName,
		TestVersion: nettest.TestVersion,
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
func (nettest *Nettest) Measure(input string) Measurement {
	measurementstarttime := time.Now().UTC().Format(dateformat)
	t0 := time.Now()
	testkeys := nettest.RunFunc(input)
	return Measurement{
		DataFormatVersion: "0.2.0",
		MeasurementStartTime: measurementstarttime,
		ProbeASN: nettest.probeASN(),
		ProbeCC: nettest.probeCC(),
		ReportID: nettest.Report.ID,
		SoftwareName: nettest.Config.SoftwareName,
		SoftwareVersion: nettest.Config.SoftwareVersion,
		TestKeys: testkeys,
		TestName: nettest.TestName,
		TestRuntime: float64(time.Now().Sub(t0)) / float64(time.Second),
		TestStartTime: nettest.TestStartTime,
		TestVersion: nettest.TestVersion,
	}
}

// Submit submits a measurement. Returns the measurementID on success and
// an error on failure. Also this method is concurrency safe.
func (nettest *Nettest) Submit(measurement Measurement) (string, error) {
	measurementID, err := nettest.Report.Update(nettest.Ctx, measurement)
	if err != nil {
		return "", err
	}
	return measurementID, nil
}

// Close closes the possibly open report.
func (nettest *Nettest) Close() error {
	if nettest.Report.ID != "" {
		return nettest.Report.Close(nettest.Ctx)
	} else {
		return nil
	}
}
