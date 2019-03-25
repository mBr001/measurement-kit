package task

// options is a field of settings
type options struct {
	BouncerBaseURL  string `json:"bouncer_base_url"`
	ConfigFilePath  string `json:"config_file_path"`
	GeoIPASNPath    string `json:"geoip_asn_path"`
	SoftwareName    string `json:"software_name"`
	SoftwareVersion string `json:"software_version"`
	WorkDirPath     string `json:"work_dir_path"`
}

// settings contains the settings
type settings struct {
	Annotations    map[string]string `json:"annotations"`
	DisabledEvents []string          `json:"disabled_events"`
	Inputs         []string          `json:"inputs"`
	InputFilepaths []string          `json:"input_filepaths"`
	LogFilepath    string            `json:"log_filepath"`
	LogLevel       string            `json:"log_level"`
	Name           string            `json:"name"`
	Options        options           `json:"options"`
	OutputFilepath options           `json:"output_filepath"`
}
