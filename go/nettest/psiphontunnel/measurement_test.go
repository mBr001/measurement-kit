package psiphontunnel

import (
	"context"
	"testing"
)

func TestNewMeasurementIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "/tmp/psiphon.json",
		Timeout: 10,
		WorkDirPath: "/tmp/",
	}
	measurement := <-NewMeasurement(context.Background(), config)
	if measurement.Failure() != "" {
		t.Error("Failure is not empty")
	}
}
