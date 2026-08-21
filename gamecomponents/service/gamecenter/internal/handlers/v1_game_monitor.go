package handlers

import (
	"fmt"
	"game/comm"
	"game/comm/define"
	"game/comm/mq"
	"game/comm/mux"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/operator"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/game/monitor",
		Handler:      v1_game_monitor,
		Desc:         "设置游戏监控",
		Kind:         "api/v1",
		ParamsSample: setGameMonitorParams{0, 0, 50, 5, 5, []*comm.MoniterRtpRangeValue{{RangeMinValue: 1, RangeMaxValue: 5, NewbieValue: 0, NotNewbieValue: 0}, {RangeMinValue: 1, RangeMaxValue: 5, NewbieValue: 0, NotNewbieValue: 0}}, []*comm.MoniterRtpRangeValue{{RangeMinValue: 1, RangeMaxValue: 5, NewbieValue: 0, NotNewbieValue: 0}, {RangeMinValue: 1, RangeMaxValue: 5, NewbieValue: 0, NotNewbieValue: 0}}},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

var moniterNewbieNumList = []int{50, 60, 70, 80, 90, 100, 150, 200, 300, 400, 500}
var moniterRTPErrorValueList = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
var moniterNumCycleList = []int{1, 5, 10, 20, 30, 50}
var moniterTypeList = []int{0, 1}

type setGameMonitorParams struct {
	MoniterType                int                          `json:"MoniterType"`
	IsMoniter                  int                          `json:"IsMoniter"`
	MoniterNewbieNum           int                          `json:"MoniterNewbieNum"`
	MoniterRTPErrorValue       int                          `json:"MoniterRTPErrorValue"`
	MoniterNumCycle            int                          `json:"MoniterNumCycle"`
	MoniterAddRTPRangeValue    []*comm.MoniterRtpRangeValue `json:"MoniterAddRTPRangeValue"`
	MoniterReduceRTPRangeValue []*comm.MoniterRtpRangeValue `json:"MoniterReduceRTPRangeValue"`
}

func v1_game_monitor(app *operator.MemApp, ps setGameMonitorParams, ret *comm.Empty) (err error) {
	if !gamedata.Contains(moniterTypeList, ps.MoniterType) || !gamedata.Contains(moniterTypeList, ps.IsMoniter) {
		err = define.NewErrCode("incorrect MoniterType or IsMoniter value", 1030)
		return
	}
	if !gamedata.Contains(moniterNewbieNumList, ps.MoniterNewbieNum) {
		err = define.NewErrCode("incorrect MoniterNewbieNum value", 1031)
		return
	}
	if !gamedata.Contains(moniterRTPErrorValueList, ps.MoniterRTPErrorValue) {
		err = define.NewErrCode("incorrect MoniterRTPErrorValue value", 1032)
		return
	}
	if !gamedata.Contains(moniterNumCycleList, ps.MoniterNumCycle) {
		err = define.NewErrCode("incorrect MoniterNumCycle value", 1033)
		return
	}
	if len(ps.MoniterAddRTPRangeValue) == 0 || len(ps.MoniterReduceRTPRangeValue) == 0 {
		err = define.NewErrCode("incorrect MoniterAddRTPRangeValue or MoniterReduceRTPRangeValue is empty", 1034)
		return
	}
	ps.MoniterAddRTPRangeValue, err = checkMoniterRangeValue(ps.MoniterRTPErrorValue, ps.MoniterAddRTPRangeValue)
	if err != nil {
		return
	}
	ps.MoniterReduceRTPRangeValue, err = checkMoniterRangeValue(ps.MoniterRTPErrorValue, ps.MoniterReduceRTPRangeValue)
	if err != nil {
		return
	}

	err = mq.Invoke("/AdminInfo/Interior/SetGameMonitor", map[string]any{
		"AppID":                      app.AppID,
		"MoniterType":                ps.MoniterType,
		"IsMoniter":                  ps.IsMoniter,
		"MoniterNewbieNum":           ps.MoniterNewbieNum,
		"MoniterRTPErrorValue":       ps.MoniterRTPErrorValue,
		"MoniterNumCycle":            ps.MoniterNumCycle,
		"MoniterAddRTPRangeValue":    ps.MoniterAddRTPRangeValue,
		"MoniterReduceRTPRangeValue": ps.MoniterReduceRTPRangeValue,
	}, &ret)
	if err != nil {
		return
	}

	return
}

func checkMoniterRangeValue(moniterRTPErrorValue int, moniterRangeValue []*comm.MoniterRtpRangeValue) (newMoniterRangeValue []*comm.MoniterRtpRangeValue, err error) {
	NextMinValue := int64(1)
	for i, v := range moniterRangeValue {
		if v.NewbieValue < 0 || v.NotNewbieValue < 0 || v.NewbieValue > 100 || v.NotNewbieValue > 100 {
			err = define.NewErrCode(fmt.Sprintf("incorrect MoniterAddRTPRangeValue value NewbieValue %v or NotNewbieValue %v ", v.NewbieValue, v.NotNewbieValue), 1035)
			return
		}
		v.RangeMinValue = NextMinValue
		v.RangeMaxValue = v.RangeMinValue + int64(moniterRTPErrorValue)
		if i == len(moniterRangeValue)-1 {
			v.RangeMaxValue = 9999999
		}
		NextMinValue = v.RangeMaxValue
	}
	return moniterRangeValue, err
}
