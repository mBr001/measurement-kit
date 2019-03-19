// Package mkgomobile contains the mobile library
package mkgomobile

// Nettest is the interface of all tests
type Nettest interface {
	// Run runs a nettest until it completes
	Run() bool
}
