package fish

import (
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	fish_proto "serve/service_fish/domain/fish/proto"
	"serve/service_fish/domain/probability"
	common_fish_proto "serve/service_fish/models/proto"
	"time"

	"github.com/gogo/protobuf/proto"
)

func psf_on_00005_Skills(f *Fish) {
	switch f.TypeId {
	case 500:
		psf_on_00005_eBroadcastFishChange(f)
		go psf_on_00005_GiantWhaleChanger(f)
	}
}

func psf_on_00005_GiantWhaleChanger(f *Fish) {
	logger.Service.Zap.Debugw("GiantWhale start.",
		"GameRoomUuid", f.RoomUuid,
		"FishUuid", f.FishUuid,
		"FishType", f.TypeId,
	)

	var changedTimeCount int32 = 0
	t := f.extraData.([]interface{})[0].(map[string]interface{})["ChangeFaceTime"].(float64)
	ticker := time.NewTicker(time.Duration(int64(t)) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if changedTimeCount >= f.fishLife {
				logger.Service.Zap.Debugw("GiantWhale stop.",
					"GameRoomUuid", f.RoomUuid,
					"FishUuid", f.FishUuid,
					"FishType", f.TypeId,
				)
				return
			}

			psf_on_00005_eBroadcastFishChange(f)
			changedTimeCount += int32(t)

		case _, ok := <-f.die:
			if !ok {
				logger.Service.Zap.Debugw("GiantWhale stop.",
					"GameRoomUuid", f.RoomUuid,
					"FishUuid", f.FishUuid,
					"FishType", f.TypeId,
				)
				return
			}

		case _, ok := <-f.out:
			if !ok {
				logger.Service.Zap.Debugw("GiantWhale stop.",
					"GameRoomUuid", f.RoomUuid,
					"FishUuid", f.FishUuid,
					"FishType", f.TypeId,
				)
				return
			}

		}
	}
}

func psf_on_00005_eBroadcastFishChange(f *Fish) {
	for {
		hitResult := probability.Service.Calc(f.GameId, f.MathModuleId, "", 500, -1, 0, 0)
		symbolId := int32(hitResult.FishTypeId)

		if f.TypeId == 500 {
			f.TypeId = symbolId
			break
		}

		if symbolId != f.TypeId {
			f.mutex.Lock()
			f.TypeId = symbolId
			f.mutex.Unlock()

			data := &fish_proto.BroadcastFishChange{
				Msgid:    common_fish_proto.EMSGID_eBroadcastFishChange,
				RoomUuid: f.RoomUuid,
				Uuid:     f.FishUuid,
				SymbolId: symbolId,
			}

			logger.Service.Zap.Debugw("GiantWhale changed.",
				"GameRoomUuid", data.RoomUuid,
				"FishUuid", data.Uuid,
				"FishType", data.SymbolId,
			)

			dataBuffer, _ := proto.Marshal(data)
			flux.Send(EMSGID_eBroadcastFishChange, Service.HashId(f.RoomUuid), f.RoomUuid, dataBuffer)
			break
		}
	}
}
