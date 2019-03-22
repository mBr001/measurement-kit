package geolookup

import (
	"context"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
)

func get(ctx context.Context, URL string) ([]byte, error) {
	request, err := http.NewRequest("GET", URL, nil)
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

type response struct {
	XMLName     xml.Name `xml:"Response"`
	IP          string   `xml:"Ip"`
	CountryCode string   `xml:"CountryCode"`
}

var validCC = regexp.MustCompile(`^[A-Z]{2}$`)

// lookupIPAndCC lookups the probe IP and probe country code (CC) and stores
// them inside of result, on success; returns an error, on failure.
func lookupIPAndCC(ctx context.Context, result *Result) error {
	data, err := get(ctx, "https://geoip.ubuntu.com/lookup")
	if err != nil {
		return err
	}
	v := response{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	if !validCC.MatchString(v.CountryCode) {
		return errors.New("Invalid country code")
	}
	if net.ParseIP(v.IP) == nil {
		return errors.New("Invalid IP address")
	}
	result.ProbeIP = v.IP
	result.ProbeCC = v.CountryCode
	return nil
}
