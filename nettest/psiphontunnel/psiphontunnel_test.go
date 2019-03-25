package psiphontunnel

import (
	"context"
	"testing"

	"github.com/measurement-kit/measurement-kit/bouncer"
)

func TestRunIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	testkeys := run(context.Background(), config)
	if testkeys.Failure != "" {
		t.Fatal("Failure is not empty")
	}
	if testkeys.BootstrapTime <= 0.0 {
		t.Fatal("BootstrapTime is not positive")
	}
}

func TestNewNettestIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	nettest := NewNettest(context.Background(), config)
	nettest.ASNDatabasePath = "../../asn.mmdb"
	nettest.SoftwareName = "measurement-kit"
	nettest.SoftwareVersion = "0.11.0-alpha"
	nettest.SelectedBouncer = &bouncer.Entry{
		Type:    "https",
		Address: "events.proteus.test.ooni.io",
	}
	err := nettest.AutomaticallySelectCollector()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("SelectedCollector: %+v", nettest.SelectedCollector)
	err = nettest.GeoLookup()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GeoInfo: %+v", nettest.GeoInfo)
	err = nettest.OpenReport()
	if err != nil {
		t.Fatal(err)
	}
	defer nettest.CloseReport()
	t.Logf("Report: %+v", nettest.Report)
	measurement := nettest.NewMeasurement()
	nettest.Measure("", &measurement)
	t.Logf("measurement: %+v", measurement)
	err = nettest.SubmitMeasurement(&measurement)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("measurementID: %+v", measurement.OOID)
}
