package mobile

import (
	"testing"
)

func TestIntegrationPsiphonTunnel(t *testing.T) {
	config := PsiphonTunnelConfig{
		ASNDBPath     : "../../../asn.mmdb",
		BouncerBaseURL: "https://events.proteus.test.ooni.io",
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	nettest := NewPsiphonTunnelNettest(&config)
	ok := nettest.Run()
	if !ok {
		t.Fatal("PsiphonTunnel failed")
	}
}
