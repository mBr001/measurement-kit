// Package bouncer contains a bouncer client implementation. This code
// implements v2.0.0 of the buncer specification.
package bouncer

import (
	"context"
	"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Config contains the bouncer configuration.
type Config struct {
	// BaseURL is the base URL to use
	BaseURL string
}

// Entry is an entry returned by a bouncer query.
type Entry struct {
	// Address is the address of a bouncer entry.
	Address string `json:"address"`

	// Type is the type of a bouncer entry.
	Type string `json:"type"`

	// Front is the front to use with "cloudfront" type entries.
	Front string `json:"front"`
}

// TODO(bassosimone): if the v2.0.0 spec is approved then we should
// change the code to remove the result indirection.

type result struct {
	Results []Entry `json:"results"`
}

func get(ctx context.Context, config Config, path string) ([]Entry, error) {
	URL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}
	URL.Path = path
	request, err := http.NewRequest("GET", URL.String(), nil)
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
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result result
	err = json.Unmarshal(data, &result)
	return result.Results, err
}

// GetCollectors queries the bouncer specified in |config| for collectors
// using the provided |context|. Returns a list of entries or an error.
func GetCollectors(ctx context.Context, config Config) ([]Entry, error) {
	return get(ctx, config, "/api/v1/collectors")
}

// GetTestHelpers is like GetCollectors but for test helpers.
func GetTestHelpers(ctx context.Context, config Config) ([]Entry, error) {
	return get(ctx, config, "/api/v1/test-helpers")
}
