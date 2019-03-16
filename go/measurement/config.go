// Package measurement contains a generic measurement.
package measurement

// Config contains generic measurement settings.
type Config struct {
	// ASNDBPath contains the ASN database path.
	ASNDBPath string

	// BouncerBaseURL is the base URL to communicate with the bouncer.
	BouncerBaseURL string

	// CABundlePath contains the CA bundle path.
	CABundlePath string

	// CountryDBPath contains the country DB path.
	CountryDBPath string

	// Timeout is the number of seconds after which the
	// specific measurement will fail.
	Timeout int
}
