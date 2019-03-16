// Package collector contains a client implementation of the OONI
// collector protocol. We implement v2.0.0.
package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Config contains the collector configuration
type Config struct {
	// BaseURL is the collector base URL
	BaseURL string
}

// Template is the template for opening a report
type Template struct {
	// ProbeASN is the probe's autonomous system number
	ProbeASN string `json:"probe_asn"`

	// ProbeCC is the probe's country code
	ProbeCC string `json:"probe_cc"`

	// SoftwareName is the app name
	SoftwareName string `json:"software_name"`

	// SoftwareVersion is the app version
	SoftwareVersion string `json:"software_version"`

	// TestName is the test name
	TestName string `json:"test_name"`

	// TestVersion is the test version
	TestVersion string `json:"test_version"`
}

// Report is an open report
type Report struct {
	// ID is the report ID
	ID string `json:"report_id"`

	// Conf is the configuration used to contact the collector
	Conf Config
}

func post(ctx context.Context, c Config, p string, b []byte) ([]byte, error) {
	URL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	URL.Path = p
	request, err := http.NewRequest("POST", URL.String(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("The request failed")
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

// Open opens a new report with the given context, configuration, and
// report template. Returns an open report or an error.
func Open(ctx context.Context, conf Config, rt Template) (Report, error) {
	report := Report{Conf: conf}
	data, err := json.Marshal(rt)
	if err != nil {
		return report, err
	}
	data, err = post(ctx, conf, "/report", data)
	if err != nil {
		return report, err
	}
	err = json.Unmarshal(data, &report)
	return report, err
}

type updateRequest struct {
	// Format is the data format
	Format string `json:"format"`

	// Content is the actual report
	Content interface{} `json:"content"`
}

type updateResponse struct {
	// ID is the measurement ID
	ID string `json:"measurement_id"`
}

// Update updates a report with the given context and opaque measurement
// content. Returns the measurement ID or an error.
func (r Report) Update(ctx context.Context, c interface{}) (string, error) {
	ureq := updateRequest{
		Format:  "json",
		Content: c,
	}
	data, err := json.Marshal(ureq)
	if err != nil {
		return "", err
	}
	data, err = post(ctx, r.Conf, fmt.Sprintf("/report/%s", r.ID), data)
	if err != nil {
		return "", err
	}
	var ures updateResponse
	err = json.Unmarshal(data, &ures)
	if err != nil {
		return "", err
	}
	return ures.ID, nil
}

// Close returns the result of closing an open report.
func (r Report) Close(ctx context.Context) error {
	_, err := post(ctx, r.Conf, fmt.Sprintf("/report/%s/close", r.ID), nil)
	return err
}
