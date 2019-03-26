package task

// options is a field of settings
type options struct {
	// BouncerBaseURL is the bouncer base URL
	BouncerBaseURL string `json:"bouncer_base_url"`

	// ConfigFilePath is the path to a config file required by a nettest
	ConfigFilePath string `json:"config_file_path"`

	// GeoIPASNPath is the path to the MaxMind ASN database
	GeoIPASNPath string `json:"geoip_asn_path"`

	// IgnoreBouncerError ignores an error when querying the bouncer
	IgnoreBouncerError bool `json:"ignore_bouncer_error"`

	// NoBouncer indicates that we don't want to use the OONI bouncer
	NoBouncer bool `json:"no_bouncer"`

	// SoftwareName is the app name
	SoftwareName string `json:"software_name"`

	// SoftwareVersion is the app version
	SoftwareVersion string `json:"software_version"`

	// WorkDirPath is the path of the dir in which the task should run
	WorkDirPath string `json:"work_dir_path"`
}

// settings contains the settings
type settings struct {
	// Annotations contains annotations for the report
	Annotations map[string]string `json:"annotations"`

	// DisabledEvents lists disables events
	DisabledEvents []string `json:"disabled_events"`

	// Inputs contains the nettest inputs
	Inputs []string `json:"inputs"`

	// InputFilepaths contains the nettest input file paths
	InputFilepaths []string `json:"input_filepaths"`

	// LogFilepath is the path of the file where to write logs
	LogFilepath string `json:"log_filepath"`

	// LogLevel is the desired level of logging
	LogLevel string `json:"log_level"`

	// Name is the task name
	Name string `json:"name"`

	// Options contains the task options
	Options options `json:"options"`

	// OutputFilepath is the path of the file where to write results
	OutputFilepath options `json:"output_filepath"`
}
