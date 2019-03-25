package task

import (
	"errors"
	"net/url"

	"github.com/measurement-kit/measurement-kit/bouncer"
	"github.com/measurement-kit/measurement-kit/nettest"
	"github.com/measurement-kit/measurement-kit/nettest/psiphontunnel"
)

// psiphontunnelNew creates a new psiphontunnel nettest.
func psiphontunnelNew(task *Task, settings *settings) (*nettest.Nettest, error) {
	if len(settings.Inputs) > 0 {
		return nil, errors.New("psiphontunnel does not require any input")
	}
	settings.Inputs = append(settings.Inputs, "") // run once
	config := psiphontunnel.Config{
		ConfigFilePath: settings.Options.ConfigFilePath,
		WorkDirPath:    settings.Options.WorkDirPath,
	}
	nettest := psiphontunnel.NewNettest(task.ctx, config)
	nettest.ASNDatabasePath = settings.Options.GeoIPASNPath
	nettest.SoftwareName = settings.Options.SoftwareName
	nettest.SoftwareVersion = settings.Options.SoftwareVersion
	if settings.Options.BouncerBaseURL != "" {
		URL, err := url.Parse(settings.Options.BouncerBaseURL)
		if err != nil {
			return nil, err // TODO(bassosimone): better wrapping of this error
		}
		nettest.SelectedBouncer = &bouncer.Entry {
			Type: "https",
			Address: URL.Host,
		}
	}
	return nettest, nil
}
