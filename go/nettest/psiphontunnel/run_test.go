package psiphontunnel

import (
	"context"
	"testing"

	"github.com/measurement-kit/measurement-kit/go/measurement"
)

func TestRunIntegration(t *testing.T) {
	config := Config{
		MeasurementConfig: measurement.Config{
			Timeout: 10,
		},
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	result := run(context.Background(), config)
	if result.Failure != "" {
		t.Fatal("Failure is not empty")
	}
	if result.BootstrapTime <= 0.0 {
		t.Fatal("BootstrapTime is not positive")
	}
}
