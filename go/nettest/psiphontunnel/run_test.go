package psiphontunnel

import (
	"context"
	"testing"
)

func TestRunIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "/tmp/psiphon.json",
		Timeout: 10,
		WorkDirPath: "/tmp/",
	}
	result := run(context.Background(), config)
	if result.Failure != "" {
		t.Error("Failure is not empty")
	}
	if result.BootstrapTime <= 0.0 {
		t.Error("BootstrapTime is not positive")
	}
}
