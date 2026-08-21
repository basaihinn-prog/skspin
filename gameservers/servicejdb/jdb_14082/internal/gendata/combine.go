package gendata

import "serve/servicejdb/jdbcomm"

var combineConfig = jdbcomm.Combine{
	{
		Name:  "普通5次不中奖",
		Count: 10,
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: getIdsNoReward(),
				Count:     5,
			},
		},
	}, {
		Name:  "普通8次不中奖",
		Count: 20,
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: getIdsNoReward(),
				Count:     8,
			},
		},
	}, {
		Name:  "普通12次不中奖",
		Count: 20,
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: getIdsNoReward(),
				Count:     12,
			},
		},
	}, {
		Name:  "1次普通（10，15】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 15, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(15, 20, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, false),
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
		Name:  "1次普通（30，40】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(30, 40, false),
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
		Name:  "1次普通（40，60】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(40, 60, false),
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
		Name:  "1次普通（60，80】+8次不中奖+4次剩余随机",
		Count: 20,
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(60, 80, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(80, 100, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(100, 150, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(150, 200, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(200, 250, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(250, 300, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(300, 350, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(350, 400, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(400, 450, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(450, 500, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(550, 600, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(600, 650, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(650, 700, false),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(30, 50, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(50, 80, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(80, 100, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(100, 150, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(150, 200, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(200, 250, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(250, 300, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(300, 350, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(350, 400, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(400, 450, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(450, 500, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(500, 550, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(550, 600, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(650, 700, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(20, 30, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, true),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(30, 50, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(50, 80, true),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(80, 100, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(150, 200, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(200, 250, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, true),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(250, 300, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(10, 20, true),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(300, 350, true),
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
		Meta: []jdbcomm.CombineMeta{
			{
				BucketIds: GetBucketIds(20, 30, true),
				Count:     1,
			}, {
				BucketIds: GetBucketIds(350, 400, true),
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

func GetCombine() jdbcomm.Combine {
	for i := range combineConfig {
		combineConfig[i].ID = i
	}
	return combineConfig
}
