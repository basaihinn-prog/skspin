package gamedata

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConvertConfig(t *testing.T) {
	tmp := Config{}
	tmp.BlackLocs = append(tmp.BlackLocs, "CN", "HK")
	content, _ := os.ReadFile("/data/game/bin/config/game_config.json")
	json.Unmarshal(content, &tmp)
	// ut.PrintJson()
	fmt.Println(tmp)

	buf, _ := yaml.Marshal(tmp)
	os.WriteFile("/data/game/bin/config/game_config.yaml", buf, 0644)

}
