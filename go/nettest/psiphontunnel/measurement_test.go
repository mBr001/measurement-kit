package psiphontunnel

import (
	"context"
	"testing"

	"github.com/measurement-kit/measurement-kit/go/nettest/measurement"
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
		t.Error("Failure is not empty")
	}
}
