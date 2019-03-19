// Package nettest contains the generic code to run a nettest.
package nettest

// Config contains the generic nettest conf.
type Config struct {
	// ASNDBPath contains the ASN DB path
	ASNDBPath string

	// BouncerBaseURL contains the bouncer base URL
	BouncerBaseURL string

	// Inputs contains the nettest inputs. Remember to provide here
	// an empty string when the nettest takes no input.
	Inputs []string

	// SoftwareName contains the software name
	SoftwareName string

	// SoftwareVersion contains the software version
	SoftwareVersion string

	// TestName contains the test name
	TestName string

	// TestVersion contains the test version
	TestVersion string
}
