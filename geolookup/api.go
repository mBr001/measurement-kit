// Package geolookup allows to geolookup a OONI probe.
//
// Specifically, the objective of the geolookup is to discover:
//
// 1. the OONI probe's IP address (aka probe IP);
//
// 2. the autonomous system number (ASN) associated to such IP (aka probe ASN);
//
// 3. the code of the country in which the IP is (aka probe CC);
//
// 4. the name associated to the probe' ASN (aka probe network name).
//
// To this end, we use a combination of remote services and local
// MaxMind databases. These values are returned in a format suitable
// for using them with the OONI collector.
//
// See https://github.com/ooni/spec/blob/master/backends/bk-003-collector.md.
package geolookup

import (
	"context"
)

// Config contains the settings
type Config struct {
	// ASNDatabasePath is the path to the ASN MaxMind database
	ASNDatabasePath string
}

// Result contains the result of the geolookup
type Result struct {
	// ProbeIP is the probe IP (e.g. `127.0.0.1`)
	ProbeIP string `json:"probe_ip"`

	// ProbeASN is the probe ASN (e.g. `AS123`)
	ProbeASN string `json:"probe_asn"`

	// ProbeCC is the probe country code (e.g. `IT`)
	ProbeCC string `json:"probe_cc"`

	// ProbeNetworkName is the name associated with the
	// probe ASN (e.g. `Vodafone`)
	ProbeNetworkName string `json:"probe_network_name"`
}

// Perform performs a geolookup using context and config. It returns
// the results on success and an error on failure.
func Perform(ctx context.Context, config Config) (Result, error) {
	var result Result
	err := lookupIPAndCC(ctx, &result)
	if err != nil {
		return result, err
	}
	err = lookupASNAndOrg(config, result.ProbeIP, &result)
	return result, err
}
