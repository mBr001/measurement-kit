// Package mobile contains the mobile library
package mobile

// Nettest is the interface of all tests
type Nettest interface {
	// Run runs a nettest until it completes
	Run() bool
}
