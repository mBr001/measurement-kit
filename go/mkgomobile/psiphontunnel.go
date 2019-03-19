package mkgomobile

import (
	"context"

	"github.com/apex/log"
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

// psiphonTunnelNettest is the private psiphontunnel mobile nettest
type psiphonTunnelNettest struct {
	// config is the configuration provided by the mobile user
	config *PsiphonTunnelConfig

	// ctx is the context for this nettest
	ctx context.Context
}

// NewPsiphonTunnelNettest creates a new PsiphonTunnel nettest
func NewPsiphonTunnelNettest(config *PsiphonTunnelConfig) Nettest {
	return &psiphonTunnelNettest{
		config: config,
		ctx:    context.Background(),
	}
}

func (psiphon *psiphonTunnelNettest) Run() bool {
	log.Infof("psiphon.config: %+v", psiphon.config)
	config := psiphontunnel.Config{
		NettestConfig: nettest.Config{
			ASNDBPath:       psiphon.config.ASNDBPath,
			BouncerBaseURL:  psiphon.config.BouncerBaseURL,
			SoftwareName:    psiphon.config.SoftwareName,
			SoftwareVersion: psiphon.config.SoftwareVersion,
		},
		ConfigFilePath: psiphon.config.ConfigFilePath,
		WorkDirPath:    psiphon.config.WorkDirPath,
	}
	log.Infof("config: %+v", config)
	nettest := psiphontunnel.NewNettest(psiphon.ctx, config)
	defer nettest.Close()
	err := nettest.DiscoverAvailableCollectors()
	if err != nil {
		log.WithError(err).Warn("nettest.DiscoverAvailableCollectors failed")
		return false
	}
	log.Infof("AvailableCollectors: %+v", nettest.AvailableCollectors)
	err = nettest.SelectCollector()
	if err != nil {
		log.WithError(err).Warn("nettest.SelectCollector failed")
		return false
	}
	log.Infof("SelectedCollector: %+v", nettest.SelectedCollector)
	err = nettest.GeoLookup()
	if err != nil {
		log.WithError(err).Warn("nettest.GeoLookup failed")
		return false
	}
	log.Infof("GeoLookupInfo: %+v", nettest.GeoLookupInfo)
	err = nettest.OpenReport()
	if err != nil {
		log.WithError(err).Warn("nettest.OpenReport failed")
		return false
	}
	log.Infof("Report: %+v", nettest.Report)
	measurement := nettest.Measure("")
	log.Infof("measurement: %+v", measurement)
	measurementID, err := nettest.Submit(measurement)
	if err != nil {
		log.WithError(err).Warn("nettest.Submit failed")
		return false
	}
	log.Infof("measurementID: %+v", measurementID)
	return true
}
