package psiphontunnel

import (
	"context"
	"testing"

	"github.com/measurement-kit/measurement-kit/go/measurement"
)

func TestNewMeasurementIntegration(t *testing.T) {
	config := Config{
		MeasurementConfig: measurement.Config{
			Timeout: 10,
		},
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath: "/tmp/",
	}
	measurement := <-NewMeasurement(context.Background(), config)
	if measurement.Failure() != "" {
		t.Fatal("Failure is not empty")
	}
}
