package platpg

import (
	"fmt"
	"game/comm/define"
	"testing"
)

func TestLaunch(t *testing.T) {
	var pg pg

	username := "testzy_RMB"
	invoke("/Player/v1/Create", define.M{
		"player_name": username,
		"nickname":    username,
		"currency":    "PHP",
	}, nil)
	pg.FundTransferIn(username, 1000.0)
	u, e := pg.LaunchGame(username, "39", "en")
	fmt.Println(u, e)
}
