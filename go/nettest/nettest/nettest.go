// Package nettest contains the generic code to run a nettest.
package nettest

import (
	"context"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/measurement-kit/measurement-kit/go/bouncer"
	"github.com/measurement-kit/measurement-kit/go/collector"
	"github.com/measurement-kit/measurement-kit/go/geolookupper"
	"github.com/measurement-kit/measurement-kit/go/model"
)

// Config contains the generic nettest conf.
type Config struct {
	// ASNDBPath contains the ASN DB path
	ASNDBPath string

	// BouncerBaseURL contains the bouncer base URL
	BouncerBaseURL string

	// Inputs contains the nettest inputs. Remember to provide here
	// an empty string when the nettest takes no input.
	Inputs []string

	// SoftwareName contains the software name
	SoftwareName string

	// SoftwareVersion contains the software version
	SoftwareVersion string

	// TestName contains the test name
	TestName string

	// TestVersion contains the test version
	TestVersion string
}

// TODO(bassosimone): we still need a way to emit events.

// Run runs a test with the specified context and configuration, invoking
// a runner function for each input. Returns whether a fatal error preventing
// from starting up the nettest itself has occurred. A nil return value does
// not indicate that some operations like submitting a report succeeded.
func Run(ctx context.Context, conf Config, fn func(string) interface{}) error {
	collectors, err := bouncer.GetCollectors(ctx, bouncer.Config{
		BaseURL: conf.BouncerBaseURL,
	})
	if err != nil {
		return err
	}
	var collectorBaseURL string
	for _, collector := range collectors {
		if collector.Type == "https" {
			collectorBaseURL = fmt.Sprintf("https://%s/", collector.Address)
			break
		}
	}
	log.Infof("Collectors: %+v", collectors)
	if collectorBaseURL == "" {
		return err
	}
	log.Infof("collectorBaseURL: %+v", collectorBaseURL)
	geolookup, err := geolookupper.Lookup(ctx, geolookupper.Config{
		ASNDBPath: conf.ASNDBPath,
	})
	if err != nil {
		return err
	}
	log.Infof("Geolookup: %+v", geolookup)
	report, err := collector.Open(ctx, collector.Config{
		BaseURL: collectorBaseURL,
	}, collector.Template{
		ProbeASN: geolookup.ProbeASN,
		ProbeCC: geolookup.ProbeCC,
		SoftwareName: conf.SoftwareName,
		SoftwareVersion: conf.SoftwareVersion,
		TestName: conf.TestName,
		TestVersion: conf.TestVersion,
	})
	if err != nil {
		return err
	}
	defer report.Close(ctx)
	log.Infof("Report.ID: %+v", report.ID)
	const dateformat = "2006-01-02 15:04:05"
	teststarttime := time.Now().UTC().Format(dateformat)
	for _, input := range(conf.Inputs) {
		measurementstarttime := time.Now().UTC().Format(dateformat)
		t0 := time.Now()
		testkeys := fn(input)
		log.Infof("TestKeys: %+v", testkeys)
		measurement := model.Measurement{
			DataFormatVersion: "0.2.0",
			MeasurementStartTime: measurementstarttime,
			ProbeASN: geolookup.ProbeASN,
			ProbeCC: geolookup.ProbeCC,
			ReportID: report.ID,
			SoftwareName: conf.SoftwareName,
			SoftwareVersion: conf.SoftwareVersion,
			TestKeys: testkeys,
			TestName: conf.TestName,
			TestRuntime: float64(time.Now().Sub(t0)) / float64(time.Second),
			TestStartTime: teststarttime,
			TestVersion: conf.TestVersion,
		}
		log.Infof("Measurement: %+v", measurement)
		measurementID, err := report.Update(ctx, measurement)
		// TODO(bassosimone): we should report this error to the app
		if err != nil {
			continue
		}
		log.Infof("measurementID: %s", measurementID)
	}
	return nil
}
