package pp

import (
	"fmt"
	"strings"

	"serve/comm/define"
	"serve/comm/plats/platcomm"
)

type Game struct {
	GameID          string
	GameName        string
	TypeDescription string
	Technology      string
	Platform        string
	DataType        string
}
type GetGameListResult struct {
	GameList []*Game
}

func (_ pp) GetGameList() (games platcomm.HotGames, err error) {
	ps := map[string]string{
		"options": "GetDataTypes",
	}

	var result GetGameListResult

	err = invoke("/getCasinoGames/", ps, &result)
	if err != nil {
		return
	}

	cfg := ppConfig

	// "Game_url": "https://happycasino.prerelease-env.biz/gs2c/openGame.do?tc=TQ2MUqFWt1By6ySmMFXr8hWK11LG7i8vUA9SJvdmhrEhnl1eaaINte3yfpnoPank&stylename=hppc_hemu&dummy="
	for _, game := range result.GameList {
		if game.DataType == "RNG" && strings.Contains(game.Platform, "MOBILE") {
			hg := &platcomm.HotGame{
				Plat: "PP",
				ID:   game.GameID,
				Name: game.GameName,
				Type: define.GameType_Slot,
				Icon: fmt.Sprintf(cfg.IconUrl, game.GameID),
			}

			games = append(games, hg)
		}
	}

	return
}
