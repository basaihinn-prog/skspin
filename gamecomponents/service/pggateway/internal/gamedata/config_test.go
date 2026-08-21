package gamedata

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	cfg := Config{
		BoBetDetailsTempUrl: "https://aaaaa.www.ccc",
		ReverseProxy: map[string]string{
			"abc": "123",
			"bbb": "456",
		},
	}

	out, _ := yaml.Marshal(cfg)
	os.Stdout.Write(out)
}
