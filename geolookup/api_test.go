package geolookup

import (
	"context"
	"testing"
)

func TestIntegrationLookup(t *testing.T) {
	config := Config{
		ASNDatabasePath: "../asn.mmdb",
	}
	result, err := Perform(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", result)
}
