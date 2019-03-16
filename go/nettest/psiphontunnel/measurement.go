package psiphontunnel

import (
	"context"

	"github.com/measurement-kit/measurement-kit/go/measurement"
)

type measurementResult struct {
	Result result
}

// Failure returns the measurement failure
func (r *measurementResult) Failure() string {
	return r.Result.Failure
}

// NewMeasurement runs a psiphontunnel measurement in a background goroutine
// using the provided context and config. Returns a channel where the
// measurement is posted. The channel will be closed after the measurement
// has been posted.
func NewMeasurement(ctx context.Context, config Config) <-chan measurement.Result {
	out := make(chan measurement.Result)
	go func() {
		defer close(out)
		var r measurementResult
		r.Result = run(ctx, config)
		out <- &r
	}()
	return out
}
