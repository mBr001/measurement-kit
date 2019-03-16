package geolookupper

import (
	"context"
)

// Config contains the settings
type Config struct {
	// ASNDBPath is the ASN DB path
	ASNDBPath string
}

// Result contains the result of the geolookup
type Result struct {
	// ProbeIP is the probe IP
	ProbeIP string

	// ProbeASN is the probe ASN
	ProbeASN string

	// ProbeCC is the probe country code
	ProbeCC string

	// ProbeNetworkName is the probe network name
	ProbeNetworkName string
}

// Lookup performs a geolookup using context and config. It returns
// the results on success and an error on failure.
func Lookup(ctx context.Context, config Config) (Result, error) {
	var result Result
	err := lookupIPAndCC(ctx, &result)
	if err != nil {
		return result, err
	}
	err = lookupASNAndOrg(config, result.ProbeIP, &result)
	return result, err
}
