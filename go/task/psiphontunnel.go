package task

import (
	"github.com/measurement-kit/measurement-kit/go/nettest/nettest"
	"github.com/measurement-kit/measurement-kit/go/nettest/psiphontunnel"
)

// psiphontunnelNew creates a new psiphontunnel nettest.
func psiphontunnelNew(task *State, settings *settings) *nettest.Nettest {
	settings.Inputs = []string{""} // XXX
	config := psiphontunnel.Config{
		NettestConfig: nettest.Config{
			ASNDBPath:       settings.Options.GeoIPASNPath,
			BouncerBaseURL:  settings.Options.BouncerBaseURL,
			SoftwareName:    settings.Options.SoftwareName,
			SoftwareVersion: settings.Options.SoftwareVersion,
		},
		ConfigFilePath: settings.Options.ConfigFilePath,
		WorkDirPath:    settings.Options.WorkDirPath,
	}
	return psiphontunnel.NewNettest(task.ctx, config)
}
