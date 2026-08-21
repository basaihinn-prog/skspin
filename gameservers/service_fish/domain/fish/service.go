package fish

import (
	"fmt"
	"reflect"
	"serve/service_fish/domain/file"
	PSF_ON_00001 "serve/service_fish/domain/fish/PSF-ON-00001"
	PSF_ON_00002 "serve/service_fish/domain/fish/PSF-ON-00002"
	PSF_ON_00003 "serve/service_fish/domain/fish/PSF-ON-00003"
	PSF_ON_00004 "serve/service_fish/domain/fish/PSF-ON-00004"
	PSF_ON_00005 "serve/service_fish/domain/fish/PSF-ON-00005"
	PSF_ON_00006 "serve/service_fish/domain/fish/PSF-ON-00006"
	PSF_ON_00007 "serve/service_fish/domain/fish/PSF-ON-00007"
	RKF_H5_00001 "serve/service_fish/domain/fish/RKF-H5-00001"
	"serve/service_fish/models"
	"time"

	"serve/fish_comm/common"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
)

const (
	ActionFishStart   = "ActionFishStart"
	ActionFishStop    = "ActionFishStop"
	ActionFishRestart = "ActionFishRestart"
)

var Service = &service{
	Id: "FishService",
	in: make(chan *flux.Action, common.Service.ChanSize),
}

type service struct {
	Id string
	in chan *flux.Action
}

func init() {
	go Service.run()
	logger.Service.Zap.Infow("Service created.",
		"Service", Service.Id,
		"Chan", fmt.Sprintf("%p", Service.in),
	)
}

func (s *service) run() {
	flux.Register(s.Id, s.in)
	fishPools := make(map[string]*fishPool)

	for {
		select {
		case action, ok := <-s.in:
			if !ok {
				for k, v := range fishPools {
					v.destroy()
					delete(fishPools, k)
				}
				return
			}

			switch action.Key().Name() {
			case ActionFishStart:
				roomUuid := action.Key().From()
				gameId := action.Payload()[0].(string)
				mathModuleId := action.Payload()[1].(string)
				f := s.born(roomUuid, gameId, mathModuleId)

				if f != nil {
					if <-f.isReady {
						fishPools[roomUuid] = f
						time.Sleep(1 * time.Second)
						flux.Send(actionFishStart, s.Id, s.Id+roomUuid, 0)
					}
				}

			case ActionFishStop:
				roomUuid := action.Key().From()

				if v, ok := fishPools[roomUuid]; ok {
					v.destroy()
					delete(fishPools, roomUuid)
				}

			case ActionFishRestart:
				for k, v := range fishPools {
					v.destroy()
					delete(fishPools, k)

					f := s.born(k, v.gameId, v.mathModuleId)

					if f != nil {
						if <-f.isReady {
							fishPools[k] = f
							time.Sleep(1 * time.Second)
							flux.Send(actionFishStart, s.Id, s.Id+k, 0)
						}
					}
				}
			}
		}
	}
}

func (s *service) born(roomUuid, gameId, mathModuleId string) *fishPool {
	switch gameId {
	case models.PSF_ON_00001, models.PGF_ON_00001:
		return s.getScriptScene(roomUuid, gameId, mathModuleId, []interface{}{
			PSF_ON_00001.ScriptA1.ScriptInfo, PSF_ON_00001.ScriptA2.ScriptInfo, PSF_ON_00001.ScriptA3.ScriptInfo,
			PSF_ON_00001.ScriptB1.ScriptInfo, PSF_ON_00001.ScriptB2.ScriptInfo, PSF_ON_00001.ScriptB3.ScriptInfo,
			PSF_ON_00001.ScriptC1.ScriptInfo, PSF_ON_00001.ScriptC2.ScriptInfo, PSF_ON_00001.ScriptC3.ScriptInfo,
		})

	case models.PSF_ON_00002, models.PSF_ON_20002:
		return s.getScriptScene(roomUuid, gameId, mathModuleId, []interface{}{
			PSF_ON_00002.ScriptA1.ScriptInfo, PSF_ON_00002.ScriptA2.ScriptInfo, PSF_ON_00002.ScriptA3.ScriptInfo,
			PSF_ON_00002.ScriptB1.ScriptInfo, PSF_ON_00002.ScriptB2.ScriptInfo, PSF_ON_00002.ScriptB3.ScriptInfo,
			PSF_ON_00002.ScriptC1.ScriptInfo, PSF_ON_00002.ScriptC2.ScriptInfo, PSF_ON_00002.ScriptC3.ScriptInfo,
			PSF_ON_00002.ScriptD1.ScriptInfo, PSF_ON_00002.ScriptD2.ScriptInfo, PSF_ON_00002.ScriptD3.ScriptInfo,
		})

	case models.PSF_ON_00003:
		return s.getScriptScene(roomUuid, gameId, mathModuleId, []interface{}{
			PSF_ON_00003.ScriptA1.ScriptInfo, PSF_ON_00003.ScriptA2.ScriptInfo, PSF_ON_00003.ScriptA3.ScriptInfo,
			PSF_ON_00003.ScriptB1.ScriptInfo, PSF_ON_00003.ScriptB2.ScriptInfo, PSF_ON_00003.ScriptB3.ScriptInfo,
			PSF_ON_00003.ScriptC1.ScriptInfo, PSF_ON_00003.ScriptC2.ScriptInfo, PSF_ON_00003.ScriptC3.ScriptInfo,
		})

	case models.PSF_ON_00004:
		return s.getScriptScene(roomUuid, gameId, mathModuleId, []interface{}{
			PSF_ON_00004.ScriptA1.ScriptInfo, PSF_ON_00004.ScriptA2.ScriptInfo, PSF_ON_00004.ScriptA3.ScriptInfo,
			PSF_ON_00004.ScriptB1.ScriptInfo, PSF_ON_00004.ScriptB2.ScriptInfo, PSF_ON_00004.ScriptB3.ScriptInfo,
			PSF_ON_00004.ScriptC1.ScriptInfo, PSF_ON_00004.ScriptC2.ScriptInfo, PSF_ON_00004.ScriptC3.ScriptInfo,
			PSF_ON_00004.ScriptD1.ScriptInfo, PSF_ON_00004.ScriptD2.ScriptInfo, PSF_ON_00004.ScriptD3.ScriptInfo,
		})

	case models.PSF_ON_00005:
		return s.getScriptScene(roomUuid, gameId, mathModuleId, []interface{}{
			PSF_ON_00005.ScriptA1.ScriptInfo, PSF_ON_00005.ScriptA2.ScriptInfo, PSF_ON_00005.ScriptA3.ScriptInfo,
			PSF_ON_00005.ScriptB1.ScriptInfo, PSF_ON_00005.ScriptB2.ScriptInfo, PSF_ON_00005.ScriptB3.ScriptInfo,
			PSF_ON_00005.ScriptC1.ScriptInfo, PSF_ON_00005.ScriptC2.ScriptInfo, PSF_ON_00005.ScriptC3.ScriptInfo,
			PSF_ON_00005.ScriptD1.ScriptInfo, PSF_ON_00005.ScriptD2.ScriptInfo, PSF_ON_00005.ScriptD3.ScriptInfo,
		})

	case models.PSF_ON_00006:
		return s.getScriptScene(roomUuid, gameId, mathModuleId, []interface{}{
			PSF_ON_00006.ScriptA1.ScriptInfo, PSF_ON_00006.ScriptA2.ScriptInfo, PSF_ON_00006.ScriptA3.ScriptInfo,
			PSF_ON_00006.ScriptB1.ScriptInfo, PSF_ON_00006.ScriptB2.ScriptInfo, PSF_ON_00006.ScriptB3.ScriptInfo,
			PSF_ON_00006.ScriptC1.ScriptInfo, PSF_ON_00006.ScriptC2.ScriptInfo, PSF_ON_00006.ScriptC3.ScriptInfo,
			PSF_ON_00006.ScriptD1.ScriptInfo, PSF_ON_00006.ScriptD2.ScriptInfo, PSF_ON_00006.ScriptD3.ScriptInfo,
		})

	case models.PSF_ON_00007:
		return s.getScriptScene(roomUuid, gameId, mathModuleId, []interface{}{
			PSF_ON_00007.ScriptA1.ScriptInfo, PSF_ON_00007.ScriptA2.ScriptInfo, PSF_ON_00007.ScriptA3.ScriptInfo,
			PSF_ON_00007.ScriptB1.ScriptInfo, PSF_ON_00007.ScriptB2.ScriptInfo, PSF_ON_00007.ScriptB3.ScriptInfo,
			PSF_ON_00007.ScriptC1.ScriptInfo, PSF_ON_00007.ScriptC2.ScriptInfo, PSF_ON_00007.ScriptC3.ScriptInfo,
			PSF_ON_00007.ScriptD1.ScriptInfo, PSF_ON_00007.ScriptD2.ScriptInfo, PSF_ON_00007.ScriptD3.ScriptInfo,
		})

	case models.RKF_H5_00001:
		return s.getScriptScene(roomUuid, gameId, mathModuleId, []interface{}{
			RKF_H5_00001.ScriptA1.ScriptInfo, RKF_H5_00001.ScriptA2.ScriptInfo, RKF_H5_00001.ScriptA3.ScriptInfo,
			RKF_H5_00001.ScriptB1.ScriptInfo, RKF_H5_00001.ScriptB2.ScriptInfo, RKF_H5_00001.ScriptB3.ScriptInfo,
			RKF_H5_00001.ScriptC1.ScriptInfo, RKF_H5_00001.ScriptC2.ScriptInfo, RKF_H5_00001.ScriptC3.ScriptInfo,
			RKF_H5_00001.ScriptD1.ScriptInfo, RKF_H5_00001.ScriptD2.ScriptInfo, RKF_H5_00001.ScriptD3.ScriptInfo,
		})

	default:
		logger.Service.Zap.Errorw(Fish_GAME_ID_INVALID,
			"Service", s.Id,
			"Chan", fmt.Sprintf("%p", s.in),
			"GameRoomUuid", roomUuid,
			"GameId", gameId,
			"MathModuleId", mathModuleId,
		)

		// TODO JOHNNY game room need implement
		errorcode.Service.Fatal(roomUuid, Fish_GAME_ID_INVALID)
		return nil
	}
}

func (s *service) HashId(gameRoomUuid string) string {
	return s.Id + gameRoomUuid
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}

func (s *service) getScriptScene(roomUuid, gameId, mathModuleId string, scripts []interface{}) *fishPool {
	scriptsScene := make(map[int]map[string]interface{})
	scriptsSceneResetTime := make(map[int]int)

	for i := 0; i < len(scripts); i++ {
		script := reflect.ValueOf(scripts[i]).Interface().(file.ScriptInfo)
		scriptsScene[i] = script.Data()
		scriptsSceneResetTime[i] = script.ResetTime()
	}

	return newFishPool(roomUuid, gameId, mathModuleId, scriptsScene, scriptsSceneResetTime)
}
