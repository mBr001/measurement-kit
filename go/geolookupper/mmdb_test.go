package geolookupper

import (
	"testing"
)

func TestIntegrationLookupASNAndOrg(t *testing.T) {
	config := Config{
		ASNDBPath: "../../asn.mmdb",
	}
	var result Result
	err := lookupASNAndOrg(config, "8.8.8.8", &result)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ASN: %s", result.ProbeASN)
	t.Logf("NetworkName: %s", result.ProbeNetworkName)
}
