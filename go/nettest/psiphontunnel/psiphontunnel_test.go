package psiphontunnel

import (
	"context"
	"testing"
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
