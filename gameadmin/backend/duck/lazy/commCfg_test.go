package lazy

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGenCommcfg(t *testing.T) {
	cfg := commCfg{
		Alert: &alert{
			TelegramNotify:          "https://api.telegram.org/bot6968591568:AAGB1jKKJ5ZW22Y7tLvkatG2erC8M0PfOW4/sendMessage?chat_id=-1002019720165&text=",
			AlertGameThreshold:      10e4,
			AlertSingleWinThreshold: 2e4,
		},
	}

	data, _ := yaml.Marshal(cfg)
	os.WriteFile("comm_config.yaml", data, 0644)
}
