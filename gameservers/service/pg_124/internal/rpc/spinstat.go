package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/mux"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/service/pg_124/internal/models"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.RegRpc("/pg_124/game/spinstat", "spin stat", "stat", spinstat, nil).SetOnlyDev()
}

func spinstat(ps ut.SpinStatPs, ret *ut.SpinStatRet) (err error) {
	ret.Balance_before, _ = slotsmongo.GetBalance(ps.Pid)
	spinMsg := define.PGParams{
		Form:   url.Values{},
		Header: map[string]any{},
	}
	spinMsg.Form.Set("cs", fmt.Sprintf("%v", ps.Bet))
	spinMsg.Form.Set("ml", fmt.Sprintf("%v", 1))
	spinMsg.Form.Set("id", fmt.Sprintf("%v", 0))
	if ps.BuyBonus {
		spinMsg.Form.Set("fb", fmt.Sprintf("2"))
	}
	var player *models.Player
	coll := db.Collection("players")
	err = coll.FindOne(context.TODO(), db.ID(ps.Pid)).Decode(&player)
	if err != nil {
		return
	}

	respon := &json.RawMessage{}
	for i := 0; i < ps.Count; i++ {
		ret.Count++
		for {
			err = spin(player, spinMsg, respon)
			if err != nil {
				return
			}
			bet := spinMsg.Header["stat_bet"].(int64)
			ret.Bet += int64(bet)
			win := spinMsg.Header["stat_win"].(int64)
			ret.Win += int64(win)
			lastId := spinMsg.Header["LastSid"].(string)
			if lastId == "" {
				break
			}

			spinMsg.Form.Set("id", fmt.Sprintf("%v", lastId))
			spinMsg.Header["stat_bet"] = int64(0)
			spinMsg.Header["stat_win"] = int64(0)
			params := strings.Split(lastId, "_")
			if len(params) != 3 || params[1] == params[2] {
				break
			}
		}
	}

	ret.Balance_after, _ = slotsmongo.GetBalance(ps.Pid)
	_, err = coll.ReplaceOne(context.TODO(), db.ID(ps.Pid), player, options.Replace().SetUpsert(true))

	return
}
