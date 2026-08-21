package slotsmongo

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"serve/comm/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BucketType struct {
	Gte         int
	Lte         int
	Type        int
	Probability int
}

type BucketTypes struct {
	Type   int
	Little map[string]int
	Big    map[string]int
}

//1	心跳型 BucketHeartBeat  删除小奖60%
//2	波动型 BucketWave      删除小奖40%
//3	仿正版（默认）BucketGov  不处理
//4	混合型 BucketMix       删除大奖20%
//5	稳定型 BucketStable     删除大奖40%
//6	高中奖率 BucketHighAward   删除大奖60%
//7	超高中奖率 BucketSuperHighAward   删除大奖80%

func InitData(tableName, code string) {
	if code == "pg_24" || code == "pg_63" {
		InitDatas(tableName, []map[string]BucketTypes{
			{
				"BucketHeartBeat": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         4,
						"Lte":         0,
						"Probability": 10,
					},
				},
				"BucketWave": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         4,
						"Lte":         0,
						"Probability": 10,
					},
				},
				"BucketMix": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         4,
						"Lte":         0,
						"Probability": 20,
					},
				},
				"BucketStable": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         4,
						"Lte":         0,
						"Probability": 40,
					},
				},
				"BucketHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         4,
						"Lte":         0,
						"Probability": 60,
					},
				},
				"BucketSuperHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         4,
						"Lte":         0,
						"Probability": 80,
					},
				},
			},
		})
		return
	}

	if code == "pg_39" {
		InitDatasPg39(tableName, map[string]BucketType{
			"BucketHeartBeat": {
				Gte:         4,
				Lte:         0,
				Type:        0,
				Probability: 10,
			},
			"BucketWave": {
				Gte:         4,
				Lte:         0,
				Type:        0,
				Probability: 10,
			},
			"BucketMix": {
				Gte:         4,
				Lte:         0,
				Type:        0,
				Probability: 20,
			},
			"BucketStable": {
				Gte:         4,
				Lte:         0,
				Type:        0,
				Probability: 40,
			},
			"BucketHighAward": {
				Gte:         4,
				Lte:         0,
				Type:        0,
				Probability: 60,
			},
			"BucketSuperHighAward": {
				Gte:         4,
				Lte:         0,
				Type:        0,
				Probability: 80,
			},
		})
		return
	}

	if code == "pp_vs15godsofwar" {
		InitDatas(tableName, []map[string]BucketTypes{
			{
				"BucketHeartBeat": {
					Type: 0,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         2,
						"Lte":         10,
						"Probability": 70,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 10,
					},
				},
				"BucketWave": {
					Type: 0,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         2,
						"Lte":         10,
						"Probability": 50,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 10,
					},
				},
				"BucketMix": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 20,
					},
				},
				"BucketStable": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 40,
					},
				},
				"BucketHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 60,
					},
				},
				"BucketSuperHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 80,
					},
				},
			},
			{
				"BucketHeartBeat": {
					Type: 5,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         112,
						"Lte":         120,
						"Probability": 70,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         120,
						"Lte":         0,
						"Probability": 10,
					},
				},
				"BucketWave": {
					Type: 5,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         112,
						"Lte":         120,
						"Probability": 50,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         120,
						"Lte":         0,
						"Probability": 10,
					},
				},
				"BucketMix": {
					Type: 5,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         120,
						"Lte":         0,
						"Probability": 20,
					},
				},
				"BucketStable": {
					Type: 5,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         120,
						"Lte":         0,
						"Probability": 40,
					},
				},
				"BucketHighAward": {
					Type: 5,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         120,
						"Lte":         0,
						"Probability": 60,
					},
				},
				"BucketSuperHighAward": {
					Type: 5,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         120,
						"Lte":         0,
						"Probability": 80,
					},
				},
			},
		})
		return
	}

	if Types1(code) {
		InitDatas(tableName, []map[string]BucketTypes{
			{
				"BucketHeartBeat": {
					Type: 0,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         2,
						"Lte":         10,
						"Probability": 60,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 5,
					},
				},
				"BucketWave": {
					Type: 0,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         2,
						"Lte":         10,
						"Probability": 40,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 5,
					},
				},
				"BucketMix": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 30,
					},
				},
				"BucketStable": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 60,
					},
				},
				"BucketHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 70,
					},
				},
				"BucketSuperHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 85,
					},
				},
			},
		})
		return
	}

	if Types2(code) {
		InitDatas(tableName, []map[string]BucketTypes{
			{
				"BucketHeartBeat": {
					Type: 0,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         2,
						"Lte":         10,
						"Probability": 60,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 5,
					},
				},
				"BucketWave": {
					Type: 0,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         2,
						"Lte":         10,
						"Probability": 40,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 5,
					},
				},
				"BucketMix": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 10,
					},
				},
				"BucketStable": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 20,
					},
				},
				"BucketHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 30,
					},
				},
				"BucketSuperHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 40,
					},
				},
			},
		})
		return
	}

	if Types3(code) {
		InitDatas(tableName, []map[string]BucketTypes{
			{
				"BucketHeartBeat": {
					Type: 0,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         2,
						"Lte":         10,
						"Probability": 60,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 5,
					},
				},
				"BucketWave": {
					Type: 0,
					Little: map[string]int{
						"Yes":         1,
						"Gte":         2,
						"Lte":         10,
						"Probability": 40,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 5,
					},
				},
				"BucketMix": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 15,
					},
				},
				"BucketStable": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 20,
					},
				},
				"BucketHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 25,
					},
				},
				"BucketSuperHighAward": {
					Type: 0,
					Little: map[string]int{
						"Yes":         0,
						"Gte":         0,
						"Lte":         0,
						"Probability": 0,
					},
					Big: map[string]int{
						"Yes":         1,
						"Gte":         10,
						"Lte":         0,
						"Probability": 30,
					},
				},
			},
		})
		return
	}

	if Types4(code) {
		num := 7
		if code == `jili_51_mc` {
			num = 5
		}
		for typee := 0; typee < num; typee++ {
			//if typee < 3 {
			//	continue
			//}
			slog.Info(fmt.Sprintf(`当前桶子type为%v`, typee))
			InitDatas1(tableName, []map[string]BucketTypes{
				{
					"BucketHeartBeat": {
						Type: typee,
						Little: map[string]int{
							"Yes":         1,
							"Gte":         typee*8 + 1,
							"Lte":         typee*8 + 2,
							"Probability": 40,
						},
						Big: map[string]int{
							"Yes":         1,
							"Gte":         typee*8 + 3,
							"Lte":         typee * 8,
							"Probability": 10,
						},
					},
					"BucketWave": {
						Type: typee,
						Little: map[string]int{
							"Yes":         1,
							"Gte":         typee*8 + 1,
							"Lte":         typee*8 + 2,
							"Probability": 30,
						},
						Big: map[string]int{
							"Yes":         1,
							"Gte":         typee*8 + 3,
							"Lte":         typee * 8,
							"Probability": 10,
						},
					},
					"BucketMix": {
						Type: typee,
						Little: map[string]int{
							"Yes":         0,
							"Gte":         typee * 8,
							"Lte":         typee * 8,
							"Probability": 0,
						},
						Big: map[string]int{
							"Yes":         1,
							"Gte":         typee*8 + 3,
							"Lte":         typee * 8,
							"Probability": 15,
						},
					},
					"BucketStable": {
						Type: typee,
						Little: map[string]int{
							"Yes":         0,
							"Gte":         typee * 8,
							"Lte":         typee * 8,
							"Probability": 0,
						},
						Big: map[string]int{
							"Yes":         1,
							"Gte":         typee*8 + 3,
							"Lte":         typee * 8,
							"Probability": 20,
						},
					},
					"BucketHighAward": {
						Type: typee,
						Little: map[string]int{
							"Yes":         0,
							"Gte":         typee * 8,
							"Lte":         typee * 8,
							"Probability": 0,
						},
						Big: map[string]int{
							"Yes":         1,
							"Gte":         typee*8 + 3,
							"Lte":         typee * 8,
							"Probability": 25,
						},
					},
					"BucketSuperHighAward": {
						Type: typee,
						Little: map[string]int{
							"Yes":         0,
							"Gte":         typee * 8,
							"Lte":         typee * 8,
							"Probability": 0,
						},
						Big: map[string]int{
							"Yes":         1,
							"Gte":         typee*8 + 3,
							"Lte":         typee * 8,
							"Probability": 30,
						},
					},
				},
			})
		}
		return
	}

	InitDatas(tableName, []map[string]BucketTypes{
		{
			"BucketHeartBeat": {
				Type: 0,
				Little: map[string]int{
					"Yes":         1,
					"Gte":         2,
					"Lte":         10,
					"Probability": 70,
				},
				Big: map[string]int{
					"Yes":         1,
					"Gte":         10,
					"Lte":         0,
					"Probability": 10,
				},
			},
			"BucketWave": {
				Type: 0,
				Little: map[string]int{
					"Yes":         1,
					"Gte":         2,
					"Lte":         10,
					"Probability": 50,
				},
				Big: map[string]int{
					"Yes":         1,
					"Gte":         10,
					"Lte":         0,
					"Probability": 10,
				},
			},
			"BucketMix": {
				Type: 0,
				Little: map[string]int{
					"Yes":         0,
					"Gte":         0,
					"Lte":         0,
					"Probability": 0,
				},
				Big: map[string]int{
					"Yes":         1,
					"Gte":         10,
					"Lte":         0,
					"Probability": 20,
				},
			},
			"BucketStable": {
				Type: 0,
				Little: map[string]int{
					"Yes":         0,
					"Gte":         0,
					"Lte":         0,
					"Probability": 0,
				},
				Big: map[string]int{
					"Yes":         1,
					"Gte":         10,
					"Lte":         0,
					"Probability": 40,
				},
			},
			"BucketHighAward": {
				Type: 0,
				Little: map[string]int{
					"Yes":         0,
					"Gte":         0,
					"Lte":         0,
					"Probability": 0,
				},
				Big: map[string]int{
					"Yes":         1,
					"Gte":         10,
					"Lte":         0,
					"Probability": 60,
				},
			},
			"BucketSuperHighAward": {
				Type: 0,
				Little: map[string]int{
					"Yes":         0,
					"Gte":         0,
					"Lte":         0,
					"Probability": 0,
				},
				Big: map[string]int{
					"Yes":         1,
					"Gte":         10,
					"Lte":         0,
					"Probability": 80,
				},
			},
		},
	})
	return
}

func Types1(status string) (yes bool) {
	//1	心跳型 BucketHeartBeat  删除小奖70%
	//2	波动型 BucketWave      删除小奖50%
	//3	仿正版（默认）BucketGov  不处理
	//4	混合型 BucketMix       删除大奖30%
	//5	稳定型 BucketStable     删除大奖60%
	//6	高中奖率 BucketHighAward   删除大奖80%
	//7	超高中奖率 BucketSuperHighAward   删除大奖95%
	tmp := []string{
		"jili_21_ols",
		"jili_35_ols2",
		"jili_109_fg",
		"jili_183_gj",
		"jili_223_fgp",
		"jili_300_fg3",
		"jili_16_kk2",
		"jili_44_fivestar",
		"jili_45_cbt",
		"jili_303_mp2",
		"jili_214_ka",
		"jili_299_pw",
		"jili_136_samba",
		"jili_208_phoenix",
		"jili_225_ctk",
		"jili_146_fb",
		"tada_21_ols",
		"tada_35_ols2",
		"tada_109_fg",
		"tada_183_gj",
		"tada_223_fgp",
		"tada_300_fg3",
		"tada_16_kk2",
		"tada_44_fivestar",
		"tada_45_cbt",
		"tada_303_mp2",
		"tada_214_ka",
		"tada_299_pw",
		"tada_136_samba",
		"tada_208_phoenix",
		"tada_225_ctk",
		"tada_146_fb",
		"pg_33",
		"pg_64",
		"pg_68",
		"pg_94",
		"pg_97",
		"pg_98",
		"pg_108",
		"pg_126",
		"pg_1682240",
		"pg_1717688",
		"pg_1879752",
		"pp_vs1dragon8",
		"pp_vs1fufufu",
		"pp_vs5aztecgems",
		"pp_vs5jjwild",
		"pp_vs5joker",
		"pp_vs5jokerdice",
		"pp_vs9aztecgemsdx",
		"pp_vs10eyestorm",
		"pp_vs10bhallbnza2",
		"pp_vs10bburger",
		"pp_vs10bxmasbnza",
		"pp_vs10cowgold",
		"pp_vs10jokerhot",
		"pp_vs10firestrike",
		"pp_vs10floatdrg",
		"pp_vs10mayangods",
		"pp_vs20doghouse2",
		"pp_vs20sh",
		"pp_vs20cjcluster",
		"pp_vs20forge",
		"pp_vswaysbufking",
		"pp_vswayswildwest",
	}
	for _, s := range tmp {
		if s == status {
			yes = true
		}
	}
	return
}

func Types2(status string) (yes bool) {
	//1	心跳型 BucketHeartBeat  删除小奖70%
	//2	波动型 BucketWave      删除小奖50%
	//3	仿正版（默认）BucketGov  不处理
	//4	混合型 BucketMix       删除大奖10%
	//5	稳定型 BucketStable     删除大奖20%
	//6	高中奖率 BucketHighAward   删除大奖30%
	//7	超高中奖率 BucketSuperHighAward   删除大奖40%
	tmp := []string{
		"pp_vs20fruitparty",
		"pp_vs25pandagold",
		"pp_vs25ultwolgol",
		"pp_vs20goldfever",
		"pp_vs25holiday",
		"pp_vs9hotroll",
		"pp_vs20aladdinsorc",
		"pp_vs20egypttrs",
		"pp_vs20fparty2",
		"pg_1489936",
		"pg_1513328",
		"pg_1529867",
		"pg_1568554",
		"pg_1580541",
		"pg_1601012",
		"pg_124",
		"pg_128",
		"pg_129",
	}
	for _, s := range tmp {
		if s == status {
			yes = true
		}
	}
	return
}

func Types3(status string) (yes bool) {
	//1	心跳型 BucketHeartBeat  删除小奖70%
	//2	波动型 BucketWave      删除小奖50%
	//3	仿正版（默认）BucketGov  不处理
	//4	混合型 BucketMix       删除大奖10%
	//5	稳定型 BucketStable     删除大奖20%
	//6	高中奖率 BucketHighAward   删除大奖30%
	//7	超高中奖率 BucketSuperHighAward   删除大奖40%
	tmp := []string{
		"pg_91",
		"pg_1418544",
		"pg_1655268",
		"pp_vs25btygold",
		"pp_vswaysmonkey",
		"pp_vswayswest",
		"jdb_14052",
		"hacksaw_1209",
		"hacksaw_1172",
		"hacksaw_1160",
		"hacksaw_1144",
	}
	for _, s := range tmp {
		if s == status {
			yes = true
		}
	}
	return
}

func Types4(status string) (yes bool) {
	//1	心跳型 BucketHeartBeat  删除小奖70%
	//2	波动型 BucketWave      删除小奖50%
	//3	仿正版（默认）BucketGov  不处理
	//4	混合型 BucketMix       删除大奖10%
	//5	稳定型 BucketStable     删除大奖20%
	//6	高中奖率 BucketHighAward   删除大奖30%
	//7	超高中奖率 BucketSuperHighAward   删除大奖40%
	tmp := []string{
		"jili_51_mc",
		"jili_302_mcp",
	}
	for _, s := range tmp {
		if s == status {
			yes = true
		}
	}
	return
}

func InitDatas(tableName string, Data []map[string]BucketTypes) {
	coll := db.Collection(tableName)
	update := bson.M{"$set": bson.M{"BucketHeartBeat": 1, "BucketWave": 1, "BucketGov": 1, "BucketMix": 1, "BucketStable": 1, "BucketHighAward": 1, "BucketSuperHighAward": 1}}
	_, err := coll.UpdateMany(context.TODO(), bson.M{}, update)
	if err != nil {
		panic(err)
	}

	for _, v := range Data {
		for kk, vv := range v {
			var cancelGatherNum int

			//1.计算小奖数量
			if vv.Little["Yes"] == 1 {
				filter := bson.M{
					"$and": []bson.M{
						{"type": vv.Type},
						{"selected": true},
						{"bucketid": bson.M{"$gte": vv.Little["Gte"]}},
						{"bucketid": bson.M{"$lte": vv.Little["Lte"]}},
					},
				}
				sortOptions := options.Find().SetProjection(bson.M{"_id": true})
				cursor, err := coll.Find(context.TODO(), filter, sortOptions)
				if err != nil {
					panic(err)
				}

				//2.根据比例计算删除小奖数量
				var ids []primitive.ObjectID
				for cursor.Next(context.TODO()) {
					var result struct {
						Id primitive.ObjectID `bson:"_id,omitempty"`
					}
					if cursor.Decode(&result) != nil {
						panic(err)
					}
					ids = append(ids, result.Id)
				}
				var UpIds []primitive.ObjectID
				cancelGatherNumLittle := (len(ids) * vv.Little["Probability"]) / 100
				UpIds = ids[:cancelGatherNumLittle]
				if len(UpIds) != 0 {
					//3.删除小奖
					updateFilter := bson.M{"_id": bson.M{"$in": UpIds}}
					updateData := bson.M{"$set": bson.M{kk: float32(0)}}
					_, err = coll.UpdateMany(context.TODO(), updateFilter, updateData)
					if err != nil {
						panic(err)
					}
					cancelGatherNum += cancelGatherNumLittle
				}
			}

			//4.计算大奖数量
			if vv.Big["Yes"] == 1 {
				filter := bson.M{
					"$and": []bson.M{
						{"type": vv.Type},
						{"selected": true},
						{"bucketid": bson.M{"$gte": vv.Big["Gte"]}},
					},
				}
				sortOptions := options.Find().SetProjection(bson.M{"_id": true}).SetSort(bson.D{{"times", -1}})
				cursor, err := coll.Find(context.TODO(), filter, sortOptions)
				if err != nil {
					panic(err)
				}

				//5.根据比例计算删除大奖数量
				var ids []primitive.ObjectID
				for cursor.Next(context.TODO()) {
					var result struct {
						Id primitive.ObjectID `bson:"_id,omitempty"`
					}
					if cursor.Decode(&result) != nil {
						panic(err)
					}
					ids = append(ids, result.Id)
				}
				var UpIds []primitive.ObjectID
				cancelGatherNumBig := (len(ids) * vv.Big["Probability"]) / 100
				UpIds = ids[:cancelGatherNumBig]
				if len(UpIds) != 0 {
					//6.删除大奖
					updateFilter := bson.M{"_id": bson.M{"$in": UpIds}}
					updateData := bson.M{"$set": bson.M{kk: float32(0)}}
					_, err = coll.UpdateMany(context.TODO(), updateFilter, updateData)
					if err != nil {
						panic(err)
					}
					cancelGatherNum += cancelGatherNumBig
				}
			}

			//7.计算原数据数量
			OldCountFilter := db.D("selected", true, "type", 0)
			OldCount, err := coll.CountDocuments(context.TODO(), OldCountFilter)
			if err != nil {
				panic(err)
			}

			//8.计算原数据time
			OldMatchStage := bson.D{{"$match", bson.D{
				{"type", vv.Type},
				{"selected", true},
			}}}
			OldGroupStage := bson.D{{"$group", bson.D{
				{"_id", nil},
				{"totalAmount", bson.D{{"$sum", "$times"}}},
			}}}
			OldCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{OldMatchStage, OldGroupStage})
			if err != nil {
				panic(err)
			}
			var OldResults []struct {
				TotalAmount float64 `bson:"totalAmount"`
			}
			err = OldCursor.All(context.TODO(), &OldResults)
			if err != nil {
				panic(err)
			}
			var OldSum float64
			if len(OldResults) > 0 {
				OldSum = OldResults[0].TotalAmount
			} else {
				panic("无旧times数据")
			}

			//9.计算新数据time
			NewMatchStage := bson.D{{"$match", bson.D{
				{kk, 1},
				{"type", vv.Type},
				{"selected", true},
			}}}
			NewGroupStage := bson.D{{"$group", bson.D{
				{"_id", nil},
				{"totalAmount", bson.D{{"$sum", "$times"}}},
			}}}
			NewCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{NewMatchStage, NewGroupStage})
			if err != nil {
				panic(err)
			}
			var NewResults []struct {
				TotalAmount float64 `bson:"totalAmount"`
			}
			err = NewCursor.All(context.TODO(), &NewResults)
			if err != nil {
				panic(err)
			}
			var NewSum float64
			if len(NewResults) > 0 {
				NewSum = NewResults[0].TotalAmount
			} else {
				panic("无新times数据")
			}

			Num := int(float64(OldCount) - NewSum/OldSum*float64(OldCount))
			tmp := Num - cancelGatherNum
			fmt.Println("字段", kk, "老数据条数", OldCount, "旧数据和", OldSum, "新数据和", NewSum, "计算值", Num, "大小删除数", cancelGatherNum, "删除条数", tmp)
			if tmp > 0 {
				///////////////针对pp_vs15godsofwar///////////////
				bucketid := 0
				if vv.Type == 5 {
					bucketid = 110
				}
				///////////////针对pp_vs15godsofwar///////////////
				findFilter := bson.M{
					"$and": []bson.M{
						{kk: 1},
						{"type": vv.Type},
						{"bucketid": bucketid},
						{"selected": true},
					},
				}
				projection := bson.M{"_id": true}
				cursors, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(projection))
				if err != nil {
					panic(err)
				}
				var idss []primitive.ObjectID
				for cursors.Next(context.TODO()) {
					var result struct {
						Id primitive.ObjectID `bson:"_id,omitempty"`
					}
					err = cursors.Decode(&result)
					if err != nil {
						panic(err)
					}
					idss = append(idss, result.Id)
				}

				var cancelGathers []primitive.ObjectID
				fmt.Println("idss", len(idss), "tmp", tmp)
				for i := 0; i < tmp; i++ {
					popNum := rand.Int64N(int64(len(idss)))
					var cancelResult primitive.ObjectID
					idss, cancelResult = popAtIndexs(idss, popNum)
					cancelGathers = append(cancelGathers, cancelResult)
				}

				updateFilters := bson.M{"_id": bson.M{"$in": cancelGathers}}
				updateDatas := bson.M{"$set": bson.M{kk: float32(0)}}

				_, err = coll.UpdateMany(context.TODO(), updateFilters, updateDatas)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func InitDatas1(tableName string, Data []map[string]BucketTypes) {
	coll := db.Collection(tableName)
	update := bson.M{"$set": bson.M{"BucketHeartBeat": 1, "BucketWave": 1, "BucketGov": 1, "BucketMix": 1, "BucketStable": 1, "BucketHighAward": 1, "BucketSuperHighAward": 1}}
	_, err := coll.UpdateMany(context.TODO(), bson.M{}, update)
	if err != nil {
		panic(err)
	}

	for _, v := range Data {
		for kk, vv := range v {
			var cancelGatherNum int

			//1.计算小奖数量
			if vv.Little["Yes"] == 1 {
				filter := bson.M{
					"$and": []bson.M{
						{"type": vv.Type},
						{"selected": true},
						{"bucketid": bson.M{"$gte": vv.Little["Gte"]}},
						{"bucketid": bson.M{"$lte": vv.Little["Lte"]}},
					},
				}
				sortOptions := options.Find().SetProjection(bson.M{"_id": true})
				cursor, err := coll.Find(context.TODO(), filter, sortOptions)
				if err != nil {
					panic(err)
				}

				//2.根据比例计算删除小奖数量
				var ids []primitive.ObjectID
				for cursor.Next(context.TODO()) {
					var result struct {
						Id primitive.ObjectID `bson:"_id,omitempty"`
					}
					if cursor.Decode(&result) != nil {
						panic(err)
					}
					ids = append(ids, result.Id)
				}
				var UpIds []primitive.ObjectID
				cancelGatherNumLittle := (len(ids) * vv.Little["Probability"]) / 100
				UpIds = ids[:cancelGatherNumLittle]
				if len(UpIds) != 0 {
					//3.删除小奖
					updateFilter := bson.M{"_id": bson.M{"$in": UpIds}}
					updateData := bson.M{"$set": bson.M{kk: float32(0)}}
					_, err = coll.UpdateMany(context.TODO(), updateFilter, updateData)
					if err != nil {
						panic(err)
					}
					cancelGatherNum += cancelGatherNumLittle
				}
			}

			//4.计算大奖数量
			if vv.Big["Yes"] == 1 {
				filter := bson.M{
					"$and": []bson.M{
						{"type": vv.Type},
						{"selected": true},
						{"bucketid": bson.M{"$gte": vv.Big["Gte"]}},
					},
				}
				sortOptions := options.Find().SetProjection(bson.M{"_id": true}).SetSort(bson.D{{"times", -1}})
				cursor, err := coll.Find(context.TODO(), filter, sortOptions)
				if err != nil {
					panic(err)
				}

				//5.根据比例计算删除大奖数量
				var ids []primitive.ObjectID
				for cursor.Next(context.TODO()) {
					var result struct {
						Id primitive.ObjectID `bson:"_id,omitempty"`
					}
					if cursor.Decode(&result) != nil {
						panic(err)
					}
					ids = append(ids, result.Id)
				}
				var UpIds []primitive.ObjectID
				cancelGatherNumBig := (len(ids) * vv.Big["Probability"]) / 100
				UpIds = ids[:cancelGatherNumBig]
				if len(UpIds) != 0 {
					//6.删除大奖
					updateFilter := bson.M{"_id": bson.M{"$in": UpIds}}
					updateData := bson.M{"$set": bson.M{kk: float32(0)}}
					_, err = coll.UpdateMany(context.TODO(), updateFilter, updateData)
					if err != nil {
						panic(err)
					}
					cancelGatherNum += cancelGatherNumBig
				}
			}

			//7.计算原数据数量
			OldCountFilter := db.D("selected", true, "type", 0)
			OldCount, err := coll.CountDocuments(context.TODO(), OldCountFilter)
			if err != nil {
				panic(err)
			}

			//8.计算原数据time
			OldMatchStage := bson.D{{"$match", bson.D{
				{"type", vv.Type},
				{"selected", true},
			}}}
			OldGroupStage := bson.D{{"$group", bson.D{
				{"_id", nil},
				{"totalAmount", bson.D{{"$sum", "$times"}}},
			}}}
			OldCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{OldMatchStage, OldGroupStage})
			if err != nil {
				panic(err)
			}
			var OldResults []struct {
				TotalAmount float64 `bson:"totalAmount"`
			}
			err = OldCursor.All(context.TODO(), &OldResults)
			if err != nil {
				panic(err)
			}
			var OldSum float64
			if len(OldResults) > 0 {
				OldSum = OldResults[0].TotalAmount
			} else {
				panic("无旧times数据")
			}

			//9.计算新数据time
			NewMatchStage := bson.D{{"$match", bson.D{
				{kk, 1},
				{"type", vv.Type},
				{"selected", true},
			}}}
			NewGroupStage := bson.D{{"$group", bson.D{
				{"_id", nil},
				{"totalAmount", bson.D{{"$sum", "$times"}}},
			}}}
			NewCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{NewMatchStage, NewGroupStage})
			if err != nil {
				panic(err)
			}
			var NewResults []struct {
				TotalAmount float64 `bson:"totalAmount"`
			}
			err = NewCursor.All(context.TODO(), &NewResults)
			if err != nil {
				panic(err)
			}
			var NewSum float64
			if len(NewResults) > 0 {
				NewSum = NewResults[0].TotalAmount
			} else {
				panic("无新times数据")
			}

			Num := int(float64(OldCount) - NewSum/OldSum*float64(OldCount))
			tmp := Num - cancelGatherNum
			fmt.Println("字段", kk, "老数据条数", OldCount, "旧数据和", OldSum, "新数据和", NewSum, "计算值", Num, "大小删除数", cancelGatherNum, "删除条数", tmp)
			if tmp > 0 {
				findFilter := bson.M{
					"$and": []bson.M{
						{kk: 1},
						{"type": vv.Type},
						{"bucketid": vv.Type * 8},
						{"selected": true},
					},
				}
				projection := bson.M{"_id": true}
				cursors, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(projection))
				if err != nil {
					panic(err)
				}
				var idss []primitive.ObjectID
				for cursors.Next(context.TODO()) {
					var result struct {
						Id primitive.ObjectID `bson:"_id,omitempty"`
					}
					err = cursors.Decode(&result)
					if err != nil {
						panic(err)
					}
					idss = append(idss, result.Id)
				}

				var cancelGathers []primitive.ObjectID
				fmt.Println("idss", len(idss), "tmp", tmp)
				if len(idss) < tmp {
					slog.Error("idss必须大于tmp，现在数据不够删")
					return
				}
				for i := 0; i < tmp; i++ {
					popNum := rand.Int64N(int64(len(idss)))
					var cancelResult primitive.ObjectID
					idss, cancelResult = popAtIndexs(idss, popNum)
					cancelGathers = append(cancelGathers, cancelResult)
				}

				updateFilters := bson.M{"_id": bson.M{"$in": cancelGathers}}
				updateDatas := bson.M{"$set": bson.M{kk: float32(0)}}

				_, err = coll.UpdateMany(context.TODO(), updateFilters, updateDatas)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func InitDatasPg39(tableName string, Data map[string]BucketType) {
	coll := db.Collection(tableName)
	update := bson.M{"$set": bson.M{"BucketHeartBeat": 1, "BucketWave": 1, "BucketGov": 1, "BucketMix": 1, "BucketStable": 1, "BucketHighAward": 1, "BucketSuperHighAward": 1}}
	_, err := coll.UpdateMany(context.TODO(), bson.M{}, update)
	if err != nil {
		panic(err)
	}

	for k, v := range Data {
		fmt.Println(k)

		filter := bson.M{
			"$and": []bson.M{
				//{"type": v.Type},
				{"selected": true},
				{"bucketid": bson.M{"$gte": v.Gte}},
				{"bucketid": bson.M{"$lte": v.Lte}},
			},
		}
		sortOptions := options.Find().SetProjection(bson.M{"_id": true})

		if v.Lte == 0 {
			filter = bson.M{
				"$and": []bson.M{
					//{"type": v.Type},
					{"selected": true},
					{"bucketid": bson.M{"$gte": v.Gte}},
				},
			}
			sortOptions = options.Find().SetProjection(bson.M{"_id": true}).SetSort(bson.D{{"times", -1}})
		}

		cursor, err := coll.Find(context.TODO(), filter, sortOptions)
		if err != nil {
			panic(err)
		}

		var ids []primitive.ObjectID
		for cursor.Next(context.TODO()) {
			var result struct {
				Id primitive.ObjectID `bson:"_id,omitempty"`
			}
			if cursor.Decode(&result) != nil {
				panic(err)
			}
			ids = append(ids, result.Id)
		}

		var UpIds []primitive.ObjectID
		cancelGatherNum := (len(ids) * v.Probability) / 100

		if v.Lte == 0 {
			UpIds = ids[:cancelGatherNum]
		} else {
			for i := 0; i < cancelGatherNum; i++ {
				popNum := rand.Int64N(int64(len(ids)))
				var id primitive.ObjectID
				ids, id = popAtIndexs(ids, popNum)
				UpIds = append(UpIds, id)
			}
		}

		updateFilter := bson.M{"_id": bson.M{"$in": UpIds}}
		updateData := bson.M{"$set": bson.M{k: float32(0)}}
		_, err = coll.UpdateMany(context.TODO(), updateFilter, updateData)
		if err != nil {
			panic(err)
		}

		/////////////////////////////////////////////////////////////////

		OldCountFilter := db.D("selected", true, "type", 0)
		OldCount, err := coll.CountDocuments(context.TODO(), OldCountFilter)
		if err != nil {
			panic(err)
		}

		////////////////////////////////

		OldMatchStage := bson.D{{"$match", bson.D{
			//{"type", v.Type},
			{"selected", true},
		}}}
		OldGroupStage := bson.D{{"$group", bson.D{
			{"_id", nil},
			{"totalAmount", bson.D{{"$sum", "$aw"}}},
		}}}
		OldCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{OldMatchStage, OldGroupStage})
		if err != nil {
			panic(err)
		}
		var OldResults []struct {
			TotalAmount float64 `bson:"totalAmount"`
		}
		err = OldCursor.All(context.TODO(), &OldResults)
		if err != nil {
			panic(err)
		}
		var OldSum float64
		if len(OldResults) > 0 {
			OldSum = OldResults[0].TotalAmount
		} else {
			panic("无旧times数据")
		}

		////////////////////////////////

		NewMatchStage := bson.D{{"$match", bson.D{
			{k, 1},
			//{"type", v.Type},
			{"selected", true},
		}}}
		NewGroupStage := bson.D{{"$group", bson.D{
			{"_id", nil},
			{"totalAmount", bson.D{{"$sum", "$aw"}}},
		}}}
		NewCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{NewMatchStage, NewGroupStage})
		if err != nil {
			panic(err)
		}
		var NewResults []struct {
			TotalAmount float64 `bson:"totalAmount"`
		}
		err = NewCursor.All(context.TODO(), &NewResults)
		if err != nil {
			panic(err)
		}
		var NewSum float64
		if len(NewResults) > 0 {
			NewSum = NewResults[0].TotalAmount
		} else {
			panic("无新times数据")
		}

		////////////////////////////////

		Num := int(float64(OldCount) - NewSum/OldSum*float64(OldCount))
		tmp := Num - cancelGatherNum
		if tmp > 0 {
			findFilter := bson.M{
				"$and": []bson.M{
					{k: 1},
					//{"type": v.Type},
					{"bucketid": 0},
					{"selected": true},
				},
			}
			projection := bson.M{"_id": true}
			cursors, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(projection))
			if err != nil {
				panic(err)
			}

			var idss []primitive.ObjectID
			for cursors.Next(context.TODO()) {
				var result struct {
					Id primitive.ObjectID `bson:"_id,omitempty"`
				}
				err = cursors.Decode(&result)
				if err != nil {
					panic(err)
				}
				idss = append(idss, result.Id)
			}

			var cancelGathers []primitive.ObjectID
			for i := 0; i < Num-cancelGatherNum; i++ {
				popNum := rand.Int64N(int64(len(idss)))
				var cancelResult primitive.ObjectID
				idss, cancelResult = popAtIndexs(idss, popNum)
				cancelGathers = append(cancelGathers, cancelResult)
			}

			updateFilters := bson.M{"_id": bson.M{"$in": cancelGathers}}
			updateDatas := bson.M{"$set": bson.M{k: float32(0)}}

			_, err = coll.UpdateMany(context.TODO(), updateFilters, updateDatas)
			if err != nil {
				panic(err)
			}
		}
	}
}

func popAtIndexs(slice []primitive.ObjectID, index int64) ([]primitive.ObjectID, primitive.ObjectID) {
	poppedElement := slice[index]
	return append(slice[:index], slice[index+1:]...), poppedElement
}
