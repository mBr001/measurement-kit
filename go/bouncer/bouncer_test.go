package bouncer

import (
	"context"
	"testing"
)

func integration(t *testing.T, f func(context.Context, Config) ([]Entry, error)) {
	entries, err := f(context.Background(), Config{
		BaseURL: "https://events.proteus.test.ooni.io",
	})
	if err != nil {
		t.Error(err)
	}
	for _, entry := range entries {
		t.Logf("%+v", entry)
	}
}

func TestIntegrationGetCollectors(t *testing.T) {
	integration(t, GetCollectors)
}

func TestIntegrationGetTestHelpers(t *testing.T) {
	integration(t, GetTestHelpers)
}
