package nettest

// Measurement contains a measurement result
type Measurement struct {
	// Annotations contains results annotations
	Annotations map[string]string `json:"annotations"`

	// DataFormatVersion is the version of the data format
	DataFormatVersion string `json:"data_format_version"`

	// ID is the locally generated measurement ID
	ID string `json:"id"`

	// Input is the measurement input
	Input string `json:"input"`

	// InputHashes contains input hashes
	InputHashes []string `json:"input_hashes"`

	// MeasurementStartTime is the time when the measurement started
	MeasurementStartTime string `json:"measurement_start_time"`

	// Options contains command line options
	Options []string `json:"options"`

	// ProbeASN contains the probe autonomous system number
	ProbeASN string `json:"probe_asn"`

	// ProbeCC contains the probe country code
	ProbeCC string `json:"probe_cc"`

	// ProbeCity contains the probe city
	ProbeCity string `json:"probe_city"`

	// ProbeIP contains the probe IP
	ProbeIP string `json:"probe_ip"`

	// ReportID contains the report ID
	ReportID string `json:"report_id"`

	// SoftwareName contains the software name
	SoftwareName string `json:"software_name"`

	// SoftwareVersion contains the software version
	SoftwareVersion string `json:"software_version"`

	// TestHelpers contains the test helpers
	TestHelpers map[string]string `json:"test_helpers"`

	// TestKeys contains the real test result
	TestKeys interface{} `json:"test_keys"`

	// TestName contains the test name
	TestName string `json:"test_name"`

	// TestRuntime contains the test runtime
	TestRuntime float64 `json:"test_runtime"`

	// TestStartTime contains the test start time
	TestStartTime string `json:"test_start_time"`

	// TestVersion contains the test version
	TestVersion string `json:"test_version"`
}
