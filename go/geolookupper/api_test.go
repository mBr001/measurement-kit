package geolookupper

import (
	"context"
	"testing"
)

func TestIntegrationLookup(t *testing.T) {
	config := Config{
		ASNDBPath: "../../asn.mmdb",
	}
	result, err := Lookup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", result)
}
