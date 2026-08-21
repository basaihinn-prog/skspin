package db

import (
	"context"
	"math/rand"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type RtpData struct {
	ID                    int     `bson:"_id"`
	Rtp                   int     `bson:"Rtp"`
	DiscardTimes          float64 `bson:"DiscardTimes"`
	HasGameDiscardTimes   float64 `bson:"HasGameDiscardTimes"`
	BuyDiscardTimes       float64 `bson:"BuyDiscardTimes"`
	SuperBuyDiscardTimes1 float64 `bson:"SuperBuyDiscardTimes1"`
	SuperBuyDiscardTimes2 float64 `bson:"SuperBuyDiscardTimes2"`
	RealRtp               float64 `bson:"RealRtp"`
	ZeroCount             int     `bson:"ZeroCount"`
	SkipZeroCount         int     `bson:"SkipZeroCount"` //要满足目标RTP，需要跳过SkipMaxCount个times较小的样本
	SkipMaxCount          int     `bson:"SkipMaxCount"`  //要满足目标RTP，需要跳过SkipMaxCount个times较大的样本
	Type                  int     `bson:"type"`
}

type RtpDataMgr struct {
	RtpData []*RtpData
}

const DEFAULT_RTP int = -1 // 等于未设置

var GRtpDataMgr *RtpDataMgr = &RtpDataMgr{}

var (
	rtpdata_projection = D(
		"Rtp", 1,
		"DiscardTimes", 1,
		"HasGameDiscardTimes", 1,
		"BuyDiscardTimes", 1,
		"SuperBuyDiscardTimes1", 1,
		"SuperBuyDiscardTimes2", 1,
		"RealRtp", 1,
		"ZeroCount", 1,
		"SkipZeroCount", 1,
		"SkipMaxCount", 1,
		"type", 1,
	)
)

func (mgr *RtpDataMgr) LoadRtpData() (err error) {
	var (
		rtpData []*RtpData
		opts    *options.FindOptions
	)
	opts = options.Find().SetProjection(rtpdata_projection).SetSort(D("type", 1, "Rtp", 1))

	coll := Collection("Rtp")
	cur, _ := coll.Find(context.TODO(), D(), opts)

	var docs []RtpData
	err = cur.All(context.TODO(), &docs)
	if err != nil {
		return
	}
	for _, doc := range docs {
		rtpData = append(rtpData, &doc)
	}
	// sort.Slice(rtpData, func(i, j int) bool {
	// 	x := rtpData[i]
	// 	y := rtpData[j]
	// 	return x.Rtp < y.Rtp
	// })
	mgr.RtpData = rtpData
	return
}
func (mgr *RtpDataMgr) GetFirstRtpData(Type int) (outInf *RtpData) {
	rtpData := mgr.RtpData
	count := len(rtpData)
	if count == 0 {
		return
	}
	for i := 0; i < count; i++ {
		inf := rtpData[i]
		if inf.Type != Type {
			continue
		}
		outInf = inf
		return
	}
	return
}
func (mgr *RtpDataMgr) GetLastRtpData(Type int) (outInf *RtpData) {
	rtpData := mgr.RtpData
	count := len(rtpData)
	if count == 0 {
		return
	}
	for i := count - 1; i >= 0; i-- {
		inf := rtpData[i]
		if inf.Type != Type {
			continue
		}
		outInf = inf
		return
	}
	return
}

func (mgr *RtpDataMgr) GetPresetRtp(userRtp, Type int) (presetRtp int, ok bool) {
	presetRtp = DEFAULT_RTP
	rtpData := mgr.RtpData
	count := len(rtpData)
	if count == 0 {
		return
	}
	if userRtp == DEFAULT_RTP {
		return
	}

	var rtpData0 *RtpData = mgr.GetFirstRtpData(Type)
	if rtpData0 != nil {
		if userRtp <= rtpData0.Rtp {
			presetRtp = rtpData0.Rtp
			ok = true
			return
		}
	}
	var rtpDataLast *RtpData = mgr.GetLastRtpData(Type)
	if rtpDataLast != nil {
		if userRtp > rtpDataLast.Rtp {
			presetRtp = rtpDataLast.Rtp
			ok = true
			return
		}
	}

	// 相等时就取这个档
	for i := count - 1; i >= 0; i-- {
		inf := rtpData[i]
		if inf.Type != Type {
			continue
		}
		if inf.Rtp == userRtp {
			presetRtp = rtpData[i].Rtp
			ok = true
			return
		} else if inf.Rtp > userRtp && i-1 >= 0 && rtpData[i-1].Rtp <= userRtp {
			presetRtp = rtpData[i-1].Rtp
			ok = true
			return
		}
	}

	// // 相等时向下取一档
	// for i := count - 1; i >= 0; i-- {
	// 	inf := rtpData[i]
	// 	if inf.Rtp == userRtp && i-1 >= 0 && rtpData[i-1].Rtp < userRtp {
	// 		presetRtp = rtpData[i-1].Rtp
	// 		ok = true
	// 		return
	// 	} else if inf.Rtp > userRtp && i-1 >= 0 && rtpData[i-1].Rtp < userRtp {
	// 		presetRtp = rtpData[i-1].Rtp
	// 		ok = true
	// 		return
	// 	}
	// }

	return
}

func (mgr *RtpDataMgr) GetRtpData(Type, presetRtp int) (outInf *RtpData) {
	rtpData := mgr.RtpData
	count := len(rtpData)
	if count == 0 {
		return
	}
	for i := count - 1; i >= 0; i-- {
		inf := rtpData[i]
		if inf.Type != Type {
			continue
		}
		if inf.Rtp == presetRtp {
			outInf = inf
			return
		}
	}
	return
}

func (mgr *RtpDataMgr) CanSkipZero(nowSkipCount int, needSkipCount int, zeroCount int) bool {
	if nowSkipCount < needSkipCount && zeroCount > 0 && rand.Intn(zeroCount) < needSkipCount {
		return true
	}
	return false
}

func GetGotoPoolPercent(userRtp int, rtp int) int {
	if userRtp == DEFAULT_RTP || rtp == DEFAULT_RTP {
		return 0
	}
	if userRtp > 100 && userRtp > rtp {
		// >100 rtp 增加额外投注比例到奖池
		return userRtp - rtp
	} else {
		return 0
	}
}
func GetRtpMngKey(preset_rtp, Type int) int {
	return preset_rtp*100 + Type
}
