package mobile

import (
	"context"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/measurement-kit/measurement-kit/go/bouncer"
	"github.com/measurement-kit/measurement-kit/go/collector"
	"github.com/measurement-kit/measurement-kit/go/geolookupper"
	"github.com/measurement-kit/measurement-kit/go/model"
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
	collectors, err := bouncer.GetCollectors(nt.ctx, bouncer.Config{
		BaseURL: nt.config.BouncerBaseURL,
	})
	if err != nil {
		log.WithError(err).Warn("bouncer.GetCollectors failed")
		return false
	}
	var collectorBaseURL string
	for _, c := range collectors {
		if c.Type == "https" {
			collectorBaseURL = fmt.Sprintf("https://%s/", c.Address)
			break
		}
	}
	log.Infof("Collectors: %+v", collectors)
	if collectorBaseURL == "" {
		log.WithError(err).Warn(
			"bouncer.GetCollectors returned no suitable collectors")
		return false
	}
	log.Infof("collectorBaseURL: %+v", collectorBaseURL)
	geolookup, err := geolookupper.Lookup(nt.ctx, geolookupper.Config{
		ASNDBPath: nt.config.ASNDBPath,
	})
	if err != nil {
		log.WithError(err).Warn("geolookupper.Lookup failed")
		return false
	}
	log.Infof("Geolookup: %+v", geolookup)
	report, err := collector.Open(nt.ctx, collector.Config{
		BaseURL: collectorBaseURL,
	}, collector.Template{
		ProbeASN: geolookup.ProbeASN,
		ProbeCC: geolookup.ProbeCC,
		SoftwareName: nt.config.SoftwareName,
		SoftwareVersion: nt.config.SoftwareVersion,
		TestName: "psiphontunnel",
		TestVersion: "0.0.1",
	})
	if err != nil {
		log.WithError(err).Warn("collector.Open failed")
		return false
	}
	defer report.Close(nt.ctx)
	log.Infof("Report.ID: %+v", report.ID)
	const dateformat = "2006-01-02 15:04:05"
	measurementstarttime := time.Now().UTC().Format(dateformat)
	t0 := time.Now()
	testkeys := psiphontunnel.Run(nt.ctx, psiphontunnel.Config{
		ConfigFilePath: nt.config.ConfigFilePath,
		WorkDirPath: nt.config.WorkDirPath,
	})
	log.Infof("TestKeys: %+v", testkeys)
	measurement := model.Measurement{
		DataFormatVersion: "0.2.0",
		MeasurementStartTime: measurementstarttime,
		ProbeASN: geolookup.ProbeASN,
		ProbeCC: geolookup.ProbeCC,
		ReportID: report.ID,
		SoftwareName: nt.config.SoftwareName,
		SoftwareVersion: nt.config.SoftwareVersion,
		TestKeys: testkeys,
		TestName: "psiphontunnel",
		TestRuntime: float64(time.Now().Sub(t0)) / float64(time.Second),
		TestStartTime: measurementstarttime,
		TestVersion: "0.0.1",
	}
	log.Infof("Measurement: %+v", measurement)
	measurementID, err := report.Update(nt.ctx, measurement)
	if err != nil {
		log.WithError(err).Warn("report.Update failed")
		return false
	}
	log.Infof("measurementID: %s", measurementID)
	return true
}
