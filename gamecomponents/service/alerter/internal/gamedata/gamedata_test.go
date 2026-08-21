package gamedata

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConvertConfig(t *testing.T) {
	tmp := Setting{
		Player: []SingleCond{
			{10000, 30000, 1.02},
			{10000, 50000, 1.02},
			{10000, 100000, 1.02},
		},
		SystemInternalSec: 300,
		SystemMin:         -10e4,
		SystemMax:         10e4,

		GameInternalSec: 300,
		GameMin:         -5e4,
		GameMax:         5e4,
	}

	buf, _ := yaml.Marshal(tmp)
	os.WriteFile("/data/game/bin/config/alerter_setting.yaml", buf, 0644)
}
