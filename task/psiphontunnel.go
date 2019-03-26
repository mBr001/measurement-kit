package task

import (
	"errors"

	"github.com/measurement-kit/measurement-kit/nettest"
	"github.com/measurement-kit/measurement-kit/nettest/psiphontunnel"
)

// psiphontunnelNew creates a new psiphontunnel nettest.
func psiphontunnelNew(task *Task, settings *settings) (*nettest.Nettest, error) {
	if len(settings.Inputs) != 0 {
		return nil, errors.New("PsiphonTunnel does not take any input")
	}
	settings.Inputs = append(settings.Inputs, "") // run once
	config := psiphontunnel.Config{
		ConfigFilePath: settings.Options.ConfigFilePath,
		WorkDirPath:    settings.Options.WorkDirPath,
	}
	return psiphontunnel.NewNettest(task.ctx, config), nil
}
