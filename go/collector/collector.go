package collector

type Settings struct {
	BaseURL string
	CABundlePath string
	Timeout int
}

type OpenRequest struct {
	ProbeASN string `json:"probe_asn"`
	ProbeCC string `json:"probe_cc"`
	SoftwareName string `json:"software_name"`
	SoftwareVersion string `json:"software_version"`
	TestName string `json:"test_name"`
	TestVersion string `json:"test_version"`
}

type OpenResponse struct {
	Good bool
	ReportID string
	Logs []string
}

func Open(request OpenRequest) OpenResponse {
	return OpenResponse{}
}
