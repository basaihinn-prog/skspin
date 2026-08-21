package main

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/comm/mux"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
)

func init() {
	mux.RegHttpWithSample("/AdminInfo/Demo/getGameList", "获取游戏列表", "AdminInfo", getDemoGameList, GetGameParam{})
	mux.RegHttpWithSample("/AdminInfo/Demo/getDemoManufacturerList", "获取厂商列表", "AdminInfo", getDemoManufacturerList, nil)
	mux.RegHttpWithSample("/AdminInfo/Demo/getDemoThemeList", "获取主题列表", "AdminInfo", getDemoThemeList, nil)
}

const (
	DEMO_FIND_ERROR_CODE = 10010 // 获取sql查询错误
	DEMO_MAP_ERROR_CODE  = 10020 // 获取数据映射错误
)

// 获取游戏请求参数
type GetGameParam struct {
	searchName      string
	themeFilterList []string
	classification  string
	sort            int64
}

type DemoGame struct {
	Animal           string `json:"animal"`
	Festival         string `json:"festival"`
	GameType         string `json:"gameType"`
	Hot              string `json:"hot"`
	MiniHot          string `json:"miniHot"`
	Id               string `json:"id"`
	Manufacturer     string `json:"manufacturer"`
	ManufacturerName string `json:"manufacturerName"`
	NewSort          string `json:"newSort"`
	Newgame          string `json:"newgame"`
	RunStatus        string `json:"runStatus"`
	Sort             string `json:"sort"`
	GameName         string `json:"gameName"`
	Icon             string `json:"icon"`
}

func getDemoGameList(r *http.Request, ps GetGameParam, response *mux.Response) (err error) {

	gameMap := map[string]*comm.Game2{}

	// 临时不上的游戏，但是需要展示在试玩站上
	GameNameConfig := map[string]string{
		"SPRIBE-1":  "Mines",
		"SPRIBE-2":  "Aviator",
		"SPRIBE-3":  "Goal",
		"SPRIBE-4":  "HiLo",
		"SPRIBE-5":  "HotLine",
		"SPRIBE-6":  "Keno",
		"SPRIBE-7":  "Keno80",
		"SPRIBE-8":  "MiniRoulette",
		"SPRIBE-9":  "Plinko",
		"SPRIBE-10": "Dice",
		"SPRIBE-11": "Limbo",
	}

	// 获取后台的游戏列表
	err, gameList := GetGameListFormatName(r.Header.Get("Lang"), nil)
	if err != nil {
		response.Error = "获取游戏列表失败"
		response.Code = DEMO_FIND_ERROR_CODE
		return err
	}

	for _, game := range gameList {
		gameMap[game.ID] = game
	}

	var responseData []*DemoGame

	for _, game := range DemoGameList {
		gameId := strings.ToLower(game.ManufacturerName) + "_" + game.Id

		if gameMap[gameId] == nil {
			game.GameName = GameNameConfig[strings.ToUpper(game.ManufacturerName)+"-"+game.Id]
			game.Icon = setting.IconAddr + "/BHdownload/" + strings.ToUpper(game.ManufacturerName) + "-" + game.Id + ".webp"
		} else {
			game.GameName = gameMap[gameId].OriginName
			game.Icon = gameMap[gameId].IconUrl
		}

		responseData = append(responseData, game)

	}

	response.Data = responseData
	response.Code = 0
	return err
}

func filterGameData(game *DemoGame, ps GetGameParam) {

	if ps.sort > 1 {

		if strings.ToLower(ps.classification) == "mini" {

		} else {

		}

	}
}

func getDemoManufacturerList(r *http.Request, ps GetGameParam, response *mux.Response) (err error) {
	var DemoManufacturerList []*GameManufacturer
	coll := db.Collection2("game", "GameManufacturers")
	find, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		response.Error = "获取游戏厂商错误"
		response.Code = DEMO_FIND_ERROR_CODE
		return err
	}
	err = find.All(context.TODO(), &DemoManufacturerList)
	if err != nil {
		response.Error = "获取游戏厂商错误"
		response.Code = DEMO_MAP_ERROR_CODE
		return err
	}

	response.Data = DemoManufacturerList
	response.Code = 0
	return err
}

func getDemoThemeList(r *http.Request, ps GetGameParam, response *mux.Response) (err error) {

	response.Data = DemoThemeList
	response.Code = 0
	return err
}
