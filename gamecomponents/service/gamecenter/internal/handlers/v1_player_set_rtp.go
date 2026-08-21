package handlers

import (
	"fmt"
	"game/comm"
	"game/comm/define"
	"game/comm/mq"
	"game/comm/mux"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/setRtp",
		Handler:      v1_player_set_rtp,
		Desc:         "设置玩家RTP",
		Kind:         "api/v1",
		ParamsSample: v1PlayerSetRtpPs{"abc", 80, "pg_89", 90, 95, 100, 1000000},
		Class:        "operator",
		GetArg0:      getArg0,
	})
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/operator/setRtp",
		Handler:      v1_operetor_set_rtp,
		Desc:         "设置商户RTP",
		Kind:         "api/v1",
		ParamsSample: v1OperatorSetRtp{80, 90, "pg_89", 95},
		Class:        "operator",
		GetArg0:      getArg0,
	})
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/operator/setExtraRtp",
		Handler:      v1_operetor_set_extra_rtp,
		Desc:         "设置商户额外RTP",
		Kind:         "api/v1",
		ParamsSample: v1OperatorSetExtraRtp{"faketrans", 1, 100, 20},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

var (
	minScore = int64(1)
	maxScore = int64(1000000)
	minMult  = 30
	maxMult  = 10000
)
var rtpList1 = []int{10, 20, 30, 40, 50, 60, 70, 80, 85, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99}
var rtpList2 = []int{93, 94, 95, 96, 97, 98, 99, 100, 102, 105, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 220, 240, 260, 280, 300}
var rtpList3 = []int{10, 20, 30, 40, 50, 60, 70, 80, 85, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 102, 105, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 220, 240, 260, 280, 300}
var rtpList4 = []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 85, 90, 95}
var rtpListMap = map[int][]int{
	1: rtpList1,
	2: rtpList2,
	3: rtpList3,
	4: rtpList4,
}
var patternList = []int{1, 2, 3, 4, 5, 6, 7}

type v1PlayerSetRtpPs struct {
	UserID            string
	Rtp               int
	GameId            string
	RemoveRTP         int
	BuyRTP            int
	PersonWinMaxMult  int   //个人最高盈利倍数
	PersonWinMaxScore int64 //个人最高盈利分值
}

type v1OperatorSetRtp struct {
	GamePattern int
	Rtp         int
	GameId      string
	BuyRTP      int
}

type v1OperatorSetExtraRtp struct {
	AppID string `json:"AppID"`
	// 保护开关
	IsProtection int `json:"IsProtection""`
	// 保护次数
	ProtectionRotateCount int `json:"ProtectionRotateCount"`
	// 保护内获取比率
	ProtectionRewardPercentLess int `json:"ProtectionRewardPercentLess"`
}

type v1PlayerSetRtpRet struct {
	PidList []int64
}

type v1OperatorSetRtpRet struct {
	AppId string
}

func v1_player_set_rtp(app *operator.MemApp, ps v1PlayerSetRtpPs, ret *v1PlayerSetRtpRet) (err error) {
	if app.PlayerRTPSettingOff == 0 {
		return define.NewErrCode("RTP modification permission is not enabled", 1021)
	}
	userIds := strings.Split(ps.UserID, ",")
	appIDandPlayerIDs := ""
	pidList := make([]int64, len(userIds))
	memplr := &operator.MemPlr{}
	for i, userId := range userIds {
		memplr, err = operator.AppMgr.GetPlr2(app, userId)
		if err != nil {
			return
		}
		if i == 0 {
			appIDandPlayerIDs = fmt.Sprintf("%s:%d", app.AppID, memplr.Pid)
		} else {
			appIDandPlayerIDs = fmt.Sprintf("%s,%s:%d", appIDandPlayerIDs, app.AppID, memplr.Pid)
		}
		pidList[i] = memplr.Pid
	}

	if !gamedata.Contains(rtpListMap[int(app.HighRTPOff)], ps.Rtp) || (ps.RemoveRTP != 0 && !gamedata.Contains(rtpListMap[3], ps.RemoveRTP)) || !gamedata.Contains(rtpListMap[4], ps.BuyRTP) {
		err = define.NewErrCode("incorrect rtp value", 1015)
		return
	}

	if ps.GameId == "" {
		err = define.NewErrCode("GameId is empty", 1004)
		return
	}

	if ps.PersonWinMaxMult != 0 {
		if ps.PersonWinMaxMult < minMult || ps.PersonWinMaxMult > maxMult {
			return define.NewErrCode("PersonWinMaxMult value is error", 1029)
		}
	} else {
		ps.PersonWinMaxMult = 100
	}
	if ps.PersonWinMaxScore != 0 {
		if ps.PersonWinMaxScore < minScore || ps.PersonWinMaxScore > maxScore {
			return define.NewErrCode("PersonWinMaxScore value is error", 1028)
		}
	} else {
		ps.PersonWinMaxScore = 1000000
	}

	operator := comm.Operator_V2{}
	gcdb.CollAdminOperator.FindOne(bson.M{"AppID": app.AppID}, &operator)
	var resInfo struct {
		code  int
		error string
	}
	if operator.BuyRTPOff == 0 && ps.BuyRTP != 0 {
		return define.NewErrCode("BuyRTP modification permission is not enabled.", 1027)
	}

	err = mq.Invoke("/AdminInfo/Interior/CreatePlayerRTP", map[string]any{
		"AppIDandPlayerID":  appIDandPlayerIDs,
		"AutoRemoveRTP":     ps.RemoveRTP,
		"ContrllRTP":        ps.Rtp,
		"GameID":            ps.GameId,
		"BuyRTP":            ps.BuyRTP,
		"PersonWinMaxMult":  ps.PersonWinMaxMult,
		"PersonWinMaxScore": ps.PersonWinMaxScore,
	}, &resInfo)

	if err != nil {
		return
	}
	if resInfo.code != 0 {
		err = define.NewErrCode(resInfo.error, resInfo.code)
		return
	}
	ret.PidList = pidList
	return
}

func v1_operetor_set_rtp(app *operator.MemApp, ps v1OperatorSetRtp, ret *v1OperatorSetRtpRet) (err error) {
	if !gamedata.Contains(rtpListMap[int(app.HighRTPOff)], ps.Rtp) || !gamedata.Contains(rtpListMap[4], ps.BuyRTP) {
		err = define.NewErrCode("incorrect rtp value", 1015)
		return
	}

	if !gamedata.Contains(patternList, ps.GamePattern) {
		err = define.NewErrCode("incorrect game pattern value", 1020)
		return
	}

	if ps.GameId == "" {
		err = define.NewErrCode("GameId is empty", 1004)
	}
	var resInfo struct {
		code  int
		error string
	}
	operator := comm.Operator_V2{}
	gcdb.CollAdminOperator.FindOne(bson.M{"AppID": app.AppID}, &operator)
	if operator.BuyRTPOff == 0 && ps.BuyRTP != 0 {
		return define.NewErrCode("BuyRTP modification permission is not enabled.", 1027)
	}

	if operator.RTPOff == 0 {
		return define.NewErrCode("RTP modification permission is not enabled", 1021)
	}

	err = mq.Invoke("/AdminInfo/Interior/UpdateBatchGameRTP", map[string]any{
		"AppID":       app.AppID,
		"RTP":         ps.Rtp,
		"GamePattern": ps.GamePattern,
		"GameList":    ps.GameId,
		"Operator":    operator,
		"BuyRTP":      ps.BuyRTP,
		"BuyRTPOff":   operator.BuyRTPOff,
	}, &resInfo)

	if err != nil {
		return
	}
	if resInfo.code != 0 {
		err = define.NewErrCode(resInfo.error, resInfo.code)
		return
	}
	ret.AppId = app.AppID
	return
}

func v1_operetor_set_extra_rtp(app *operator.MemApp, ps v1OperatorSetExtraRtp, ret *v1OperatorSetRtpRet) (err error) {
	if ps.AppID != app.AppID {
		err = define.NewErrCode("AppID is not correct", 1002)
		return
	}
	sysconfig := &comm.SystemConfig{}
	err = gcdb.CollSystemConfig.FindOne(bson.M{}, &sysconfig)
	if err != nil {
		return err
	}
	if ps.IsProtection != 0 && ps.IsProtection != 1 {
		err = define.NewErrCode("incorrect protection switch value", 1024)
		return
	}
	if ps.ProtectionRewardPercentLess <= 0 || ps.ProtectionRewardPercentLess > 20 {
		err = define.NewErrCode("incorrect extra rtp value", 1015)
		return
	}
	if ps.ProtectionRotateCount <= 0 || ps.ProtectionRotateCount > 100 {
		err = define.NewErrCode("incorrect protection rotate count value", 1025)
		return
	}

	var resInfo struct {
		code  int
		error string
	}

	err = mq.Invoke("/AdminInfo/Interior/setExtraRtpConfig", map[string]any{
		"AppID":                       app.AppID,
		"IsProtection":                ps.IsProtection,
		"ProtectionRotateCount":       ps.ProtectionRotateCount,
		"ProtectionRewardPercentLess": ps.ProtectionRewardPercentLess,
	}, &resInfo)

	if err != nil {
		return
	}
	if resInfo.code != 0 {
		err = define.NewErrCode(resInfo.error, resInfo.code)
		return
	}
	ret.AppId = app.AppID
	return
}
