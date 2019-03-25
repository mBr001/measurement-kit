package geolookup

import (
	"context"
	"testing"
)

func TestIntegrationLookupIPAndCC(t *testing.T) {
	var result Result
	err := lookupIPAndCC(context.Background(), &result)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("IP: %s", result.ProbeIP)
	t.Logf("CC: %s", result.ProbeCC)
}
