package psiphontunnel

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	config := Config{
		ConfigFilePath: "/tmp/psiphon.json",
		Timeout: 10,
		WorkDirPath: "/tmp/",
	}
	result := Run(config)
	fmt.Printf("%+v", result)
}
