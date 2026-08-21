package gamedata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertConfig(t *testing.T) {
	// tmp := Config{}
	// tmp.BlackLocs = append(tmp.BlackLocs, "CN", "HK")
	// content, _ := os.ReadFile("/data/game/bin/config/game_config.json")
	// json.Unmarshal(content, &tmp)
	// // ut.PrintJson()
	// fmt.Println(tmp)

	// buf, _ := yaml.Marshal(tmp)
	// os.WriteFile("/data/game/bin/config/game_config.yaml", buf, 0644)

}

func TestLimit(t *testing.T) {
	m := TransferLimitMap{
		"default": 100,
		"fake":    80,
	}

	assert.Equal(t, 80.0, m.GetLimit("fake"))
	assert.Equal(t, 100.0, m.GetLimit("fake1"))

	m = nil
	assert.Equal(t, 1e8, m.GetLimit("fake"))
}
