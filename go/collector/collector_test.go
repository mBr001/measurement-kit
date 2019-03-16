package collector

import (
	"context"
	"testing"

	"github.com/measurement-kit/measurement-kit/go/model"
)

type fakeTestKeys struct {
	ClientResolver string `json:"client_resolver"`
}

func makeMeasurement(rt Template, ID string) model.Measurement {
	return model.Measurement{
		DataFormatVersion:    "0.2.0",
		ID:                   "bdd20d7a-bba5-40dd-a111-9863d7908572",
		MeasurementStartTime: "2018-11-01 15:33:20",
		ProbeASN:             rt.ProbeASN,
		ProbeCC:              rt.ProbeCC,
		ReportID:             ID,
		SoftwareName:         rt.SoftwareName,
		SoftwareVersion:      rt.SoftwareVersion,
		TestKeys: fakeTestKeys{
			ClientResolver: "91.80.37.104",
		},
		TestName:      rt.TestName,
		TestRuntime:   5.0565230846405,
		TestStartTime: "2018-11-01 15:33:17",
		TestVersion:   rt.TestVersion,
	}
}

func TestIntegration(t *testing.T) {
	config := Config{
		BaseURL: "https://collector-sandbox.ooni.io",
	}
	template := Template{
		ProbeASN:        "AS0",
		ProbeCC:         "ZZ",
		SoftwareName:    "measurement-kit",
		SoftwareVersion: "0.0.1",
		TestName:        "dummy",
		TestVersion:     "0.0.1",
	}
	ctx := context.Background()
	report, err := Open(ctx, config, template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Report %s: OPEN", report.ID)
	defer func() {
		err := report.Close(ctx)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Report %s: CLOSED", report.ID)
	}()
	ID, err := report.Update(ctx, makeMeasurement(template, report.ID))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Report %s: UPDATED (Measurement ID: %s)", report.ID, ID)
}
