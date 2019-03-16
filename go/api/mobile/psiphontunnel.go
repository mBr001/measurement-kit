package mobile

import (
	"context"

	"github.com/measurement-kit/measurement-kit/go/nettest/nettest"
	"github.com/measurement-kit/measurement-kit/go/nettest/psiphontunnel"

)

// PsiphonTunnelConfig contains the psiphontunnel nettest config
type PsiphonTunnelConfig struct {
	// ASNDBPath contains the ASN DB path
	ASNDBPath string

	// BouncerBaseURL contains the bouncer base URL
	BouncerBaseURL string

	// ConfigFilePath contains the psiphon config file path
	ConfigFilePath string

	// SoftwareName contains the software name
	SoftwareName string

	// SoftwareVersion contains the software version
	SoftwareVersion string

	// WorkDirPath contains the psiphon workdir path
	WorkDirPath string
}

type psiphonTunnelNettest struct {
	config *PsiphonTunnelConfig
	ctx context.Context
}

// NewPsiphonTunnelNettest creates a new PsiphonTunnel nettest
func NewPsiphonTunnelNettest(config *PsiphonTunnelConfig) Nettest {
	return &psiphonTunnelNettest{
		config: config,
		ctx: context.Background(),
	}
}

func (nt *psiphonTunnelNettest) Run() bool {
	config := nettest.Config{
		ASNDBPath: nt.config.ASNDBPath,
		BouncerBaseURL: nt.config.BouncerBaseURL,
		Inputs: []string{""},
		SoftwareName: nt.config.SoftwareName,
		SoftwareVersion: nt.config.SoftwareVersion,
		TestName: "psiphontunnel",
		TestVersion: "0.0.1",
	}
	err := nettest.Run(nt.ctx, config, func(input string) interface{} {
		return psiphontunnel.Run(nt.ctx, psiphontunnel.Config{
			ConfigFilePath: nt.config.ConfigFilePath,
			WorkDirPath: nt.config.WorkDirPath,
		})
	})
	return err == nil
}
