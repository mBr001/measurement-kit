// Package internal contains internal nettest bits
package internal

import (
	"context"
	"time"

	"github.com/measurement-kit/measurement-kit/go/ooni/nettest"
)

// NewNettest returns a new nettest instance.
func NewNettest(ctx context.Context, config nettest.Config, name, version string, fn nettest.Func) *nettest.Nettest {
	return &nettest.Nettest{
		Config:        config,
		Ctx:           ctx,
		Func:          fn,
		TestName:      name,
		TestStartTime: time.Now().UTC().Format(nettest.DateFormat),
		TestVersion:   version,
	}
}
