package platpg

import (
	"encoding/json"
	"game/comm/define"
	"game/comm/mux"
	"path"
)

type IPlat interface {
	LaunchGame(uid string, game, lang string) (url string, err error)
	FundTransferIn(uid string, amount float64) (status string)

	GetBalance(uid string) (balance float64, err error)
	FundTransferOut(uid string) (amount float64, status string)

	GetGameList() (games HotGames, err error)

	Regist(uid string) (err error)
}

// var Plats = map[string]IPlat{}

// func Start() {
// 	for k, plat := range Plats {
// 		regPlat("plat/"+k, plat)
// 	}
// }

func regPlat(ns string, plat IPlat) {
	type LaunchGamePs struct {
		UID  string
		Game string
		Lang string
	}
	type LaunchGameRet struct {
		Url string
	}

	mux.RegHttpWithSample(path.Join("/", ns, "LaunchGame"), "拉起游戏", ns, func(ps LaunchGamePs, ret *LaunchGameRet) (err error) {
		// ba plat.GetBalance(ps.UID)
		ret.Url, err = plat.LaunchGame(ps.UID, ps.Game, ps.Lang)
		return
	}, LaunchGamePs{"123456", "39", "en"})

	type FundTransferInPs struct {
		UID    string
		Amount float64
	}
	type FundTransferInRet struct {
		Status string
	}
	mux.RegHttpWithSample(path.Join("/", ns, "FundTransferIn"), "带入", ns, func(ps FundTransferInPs, ret *FundTransferInRet) (err error) {
		ret.Status = plat.FundTransferIn(ps.UID, ps.Amount)
		return
	}, FundTransferInPs{"123456", 1000.0})

	type FundTransferOutPs struct {
		UID string
	}
	type FundTransferOutRet struct {
		Amount float64
		Status string
	}
	mux.RegHttpWithSample(path.Join("/", ns, "FundTransferOut"), "带出", ns, func(ps FundTransferOutPs, ret *FundTransferOutRet) (err error) {
		ret.Amount, ret.Status = plat.FundTransferOut(ps.UID)
		return
	}, FundTransferOutPs{"123456"})

	type GetBalancePs struct {
		UID string
	}
	type GetBalanceRet struct {
		Balance float64
	}
	mux.RegHttpWithSample(path.Join("/", ns, "GetBalance"), "获取玩家余额", ns, func(ps GetBalancePs, ret *GetBalanceRet) (err error) {
		ret.Balance, err = plat.GetBalance(ps.UID)
		return
	}, GetBalancePs{"123456"})

	type GetGameListPs struct {
		// GameType string
	}
	// type GetGameListRet struct {
	// 	Games HotGames
	// }
	type GetGameListRet = json.RawMessage

	var gamelist json.RawMessage
	mux.RegHttpWithSample(path.Join("/", ns, "GetGameList"), "获取游戏列表", ns, func(ps GetGameListPs, ret *GetGameListRet) (err error) {
		if len(gamelist) != 0 {
			*ret = gamelist
			return
		}
		games, err := plat.GetGameList()
		if err != nil {
			return
		}

		buf, _ := json.Marshal(define.M{
			"List": games,
		})
		gamelist = json.RawMessage(buf)
		*ret = gamelist
		return
	}, GetGameListPs{})

	type RegistPs struct {
		UID string
	}
	type RegistRet struct {
	}
	mux.RegHttpWithSample(path.Join("/", ns, "Regist"), "注册", ns, func(ps RegistPs, ret *RegistRet) (err error) {
		err = plat.Regist(ps.UID)
		return
	}, RegistPs{"123456"})

}
