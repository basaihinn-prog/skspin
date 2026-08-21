package main

import (
	"encoding/json"
	"game/duck/lazy"
	"game/duck/logger"
	"time"
)

func initConfig() {
	w := lazy.ConfigManager

	w.Watch("comm_choushui.json", loadFlowTransferYeji)
}

// flow transfer YeJi start
type flow_transfer_yeji_config struct {
	BaiRen     int
	ZhaJinHua  int
	FishChuShu float64
	//流水换算业绩 start
	Slots         float64
	BaiJiaLe      float64
	HongHeiDaZhan float64
	LongHu        float64
	NiuMoWang     float64
	YaoQianShu    float64
	FeiQin        float64
	Hilo          float64
	//流水换算业绩 end
	FishNoPlayerInWarning  time.Duration
	HiloNoPlayerInWarning  time.Duration
	SlotsNoPlayerInWarning time.Duration
}

var FlowTransferYejiConfig = &flow_transfer_yeji_config{}

func loadFlowTransferYeji(buf []byte) error {
	var tmp flow_transfer_yeji_config
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	logger.Infof("flow_transfer_yeji_config:%+v\n", tmp)
	FlowTransferYejiConfig = &tmp
	return nil
}
