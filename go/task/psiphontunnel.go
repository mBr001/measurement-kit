package task

import (
	"github.com/measurement-kit/measurement-kit/go/ooni/nettest"
	"github.com/measurement-kit/measurement-kit/go/ooni/nettest/psiphontunnel"
)

// psiphontunnelNew creates a new psiphontunnel nettest.
func psiphontunnelNew(task *Task, settings *settings) *nettest.Nettest {
	settings.Inputs = []string{""} // XXX
	config := psiphontunnel.Config{
		NettestConfig: nettest.Config{
			ASNDatabasePath: settings.Options.GeoIPASNPath,
			BouncerBaseURL:  settings.Options.BouncerBaseURL,
			SoftwareName:    settings.Options.SoftwareName,
			SoftwareVersion: settings.Options.SoftwareVersion,
		},
		ConfigFilePath: settings.Options.ConfigFilePath,
		WorkDirPath:    settings.Options.WorkDirPath,
	}
	return psiphontunnel.NewNettest(task.ctx, config)
}
