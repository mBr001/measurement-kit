package psiphontunnel

import (
	"context"
	"testing"

	"github.com/measurement-kit/measurement-kit/go/nettest/nettest"
)

func TestRunIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	testkeys := Run(context.Background(), config)
	if testkeys.Failure != "" {
		t.Fatal("Failure is not empty")
	}
	if testkeys.BootstrapTime <= 0.0 {
		t.Fatal("BootstrapTime is not positive")
	}
}

func TestNewNettestIntegration(t *testing.T) {
	config := Config{
		NettestConfig: nettest.Config{
			ASNDBPath:       "../../../asn.mmdb",
			BouncerBaseURL:  "https://events.proteus.test.ooni.io",
			SoftwareName:    "measurement-kit",
			SoftwareVersion: "0.11.0-alpha",
		},
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	nettest := NewNettest(context.Background(), config)
	defer nettest.Close()
	err := nettest.DiscoverAvailableCollectors()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("AvailableCollectors: %+v", nettest.AvailableCollectors)
	err = nettest.SelectCollector()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("SelectedCollector: %+v", nettest.SelectedCollector)
	err = nettest.GeoLookup()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GeoLookupInfo: %+v", nettest.GeoLookupInfo)
	err = nettest.OpenReport()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Report: %+v", nettest.Report)
	measurement := nettest.Measure("")
	t.Logf("measurement: %+v", measurement)
	measurementID, err := nettest.Submit(measurement)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("measurementID: %+v", measurementID)
}
