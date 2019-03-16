// Package measurement contains a generic measurement.
package measurement

// Result is a measurement result.
type Result interface {
	// Failure is an empty string on success and an error describing
	// the measurement error on failure.
	Failure() string
}
