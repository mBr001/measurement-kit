package task

import (
	"fmt"
	"testing"
)

func TestPsiphonTunnelIntegration(t *testing.T) {
	config := `{
		"name": "psiphontunnel",
		"options": {
			"bouncer_base_url": "https://events.proteus.test.ooni.io",
			"config_file_path": "/tmp/psiphon.json",
			"geoip_asn_path": "../../asn.mmdb",
			"software_name": "measurement-kit",
			"software_version": "0.11.0-alpha",
			"work_dir_path": "/tmp"
		}
	}`
	task := Start(config)
	for !IsDone(task) {
		event := WaitForNextEvent(task)
		fmt.Println(event)
	}
}
