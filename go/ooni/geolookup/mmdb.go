package geolookup

import (
	"fmt"
	"net"

	// oschwald is a maxmind developer, therefore I expect this package
	// to be reasonably support even though it's not official
	"github.com/oschwald/maxminddb-golang"
)

// lookupASNAndOrg lookups the probe ASN and organization, and stores them
// inside of result, on success; returns an error, on failure.
func lookupASNAndOrg(config Config, IP string, result *Result) error {
	db, err := maxminddb.Open(config.ASNDatabasePath)
	if err != nil {
		return err
	}
	defer db.Close()
	dataIP := net.ParseIP(IP)
	var record struct {
		ASN int    `maxminddb:"autonomous_system_number"`
		Org string `maxminddb:"autonomous_system_organization"`
	}
	err = db.Lookup(dataIP, &record)
	if err != nil {
		return err
	}
	result.ProbeASN = fmt.Sprintf("AS%d", record.ASN)
	result.ProbeNetworkName = record.Org
	return nil
}
