package gendata

import "serve/servicepp/ppcomm"

var combineConfig = ppcomm.Combine{
	{
		Name:  "普通5次不中奖",
		Count: 10,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: getIdsNoReward(),
				Count:     5,
			},
		},
	}, {
		Name:  "普通8次不中奖",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: getIdsNoReward(),
				Count:     8,
			},
		},
	}, {
		Name:  "普通12次不中奖",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: getIdsNoReward(),
				Count:     12,
			},
		},
	}, {
		Name:  "1次普通（10，15】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 15, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（15，20】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(15, 20, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（20，30】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（30，50】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(30, 50, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（40，50】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(40, 50, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（50，80】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(50, 80, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（80，100】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(80, 100, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（100，150】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(100, 150, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（150，200】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(150, 200, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（200，250】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(200, 250, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（250，300】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(250, 300, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（300，350】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(300, 350, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（350，400】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(350, 400, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（400，450】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(400, 450, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（450，500】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(450, 500, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（550，600】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(550, 600, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（600，650】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(600, 650, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次普通（650，700】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(650, 700, false, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	},

	{
		Name:  "1次小游戏（10，20】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（20，30】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（30，50】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(30, 50, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（50，80】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(50, 80, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（80，100】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(80, 100, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（100，150】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(100, 150, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（150，200】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(150, 200, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（200，250】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(200, 250, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（250，300】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(250, 300, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（300，350】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(300, 350, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（350，400】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(350, 400, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（400，450】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(400, 450, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（450，500】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(450, 500, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（500，550】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(500, 550, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（550，600】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(550, 600, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（650，700】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(650, 700, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     8,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	},

	{
		Name:  "1次小游戏（10，20】+1次小游戏（20，30】+12次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(20, 30, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     12,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（20，30】+1次小游戏（30，50】+12次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(30, 50, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     12,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（50，80】+1次小游戏（80，100】+12次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(50, 80, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(80, 100, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     12,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（10，20】+1次小游戏（150，200】+12次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(150, 200, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     12,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（10，20】+1次小游戏（200，250】+12次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(200, 250, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     12,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（20，30】+1次小游戏（250，300】+12次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(250, 300, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     12,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（10，20】+1次小游戏（300，350】+12次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(300, 350, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     12,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	}, {
		Name:  "1次小游戏（20，30】+1次小游戏（350，400】+12次不中奖+4次剩余随机",
		Count: 20,
		Meta: []ppcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(350, 400, true, ppcomm.GameTypeNormal),
				Count:     1,
			}, {
				BucketIds: getIdsNoReward(),
				Count:     12,
			}, {
				BucketIds: getIdsFree(),
				Count:     4,
			},
		},
	},
}

func GetCombine() ppcomm.Combine {
	for i := range combineConfig {
		combineConfig[i].ID = i
	}
	return combineConfig
}
