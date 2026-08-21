package fish

import (
	"fmt"
	"reflect"
	PSF_ON_00001 "serve/service_fish/domain/fish/PSF-ON-00001"
	PSF_ON_00002 "serve/service_fish/domain/fish/PSF-ON-00002"
	PSF_ON_00003 "serve/service_fish/domain/fish/PSF-ON-00003"
	PSF_ON_00004 "serve/service_fish/domain/fish/PSF-ON-00004"
	PSF_ON_00005 "serve/service_fish/domain/fish/PSF-ON-00005"
	PSF_ON_00006 "serve/service_fish/domain/fish/PSF-ON-00006"
	PSF_ON_00007 "serve/service_fish/domain/fish/PSF-ON-00007"
	RKF_H5_00001 "serve/service_fish/domain/fish/RKF-H5-00001"
	fish_proto "serve/service_fish/domain/fish/proto"
	"serve/service_fish/models"
	common_proto "serve/service_fish/models/proto"
	"strconv"
	"time"

	"serve/fish_comm/common"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
)

const (
	EMSGID_eAllFishCall         = "EMSGID_eAllFishCall"
	EMSGID_eAllFishRecall       = "EMSGID_eAllFishRecall"
	EMSGID_eBroadcastFishIn     = "EMSGID_eBroadcastFishIn"
	EMSGID_eBroadcastFishOut    = "EMSGID_eBroadcastFishOut"
	EMSGID_eBroadcastRoundStart = "EMSGID_eBroadcastRoundStart"
	EMSGID_eBroadcastRoundEnd   = "EMSGID_eBroadcastRoundEnd"
	EMSGID_eBroadcastFishChange = "EMSGID_eBroadcastFishChange"
	EMSGID_eBroadcastGroupIn    = "EMSGID_eBroadcastGroupIn"
	ActionChangeNextScene       = "ActionChangeNextScene"
	ActionFishExistCheck        = "ActionFishExistCheck"
	ActionFishExist             = "ActionFishExist"
	ActionFishNotExist          = "ActionFishNotExist"
	ActionMultiFishExistCheck   = "ActionMultiFishExistCheck"
	ActionMultiFishiesCheck     = "ActionMultiFishiesCheck"
	ActionFishDie               = "ActionFishDie"
	actionFishIn                = "actionFishIn"
	actionFishOut               = "actionFishOut"
	actionFishStart             = "actionFishStart"
	actionFishReset             = "actionFishReset"
	actionFishDestroy           = "actionFishDestroy"
	actionGroupIn               = "actionGroupIn"
)

type fishPool struct {
	roomUuid              string
	gameId                string
	mathModuleId          string // not used now
	in                    chan *flux.Action
	scriptsScene          map[int]map[string]interface{}
	scriptsSceneResetTime map[int]int
	goroutinePool         map[chan bool]bool
	isReady               chan bool
}

func newFishPool(roomUuid, gameId, mathModuleId string, scriptsScene map[int]map[string]interface{}, scriptsSceneResetTime map[int]int) *fishPool {
	f := &fishPool{
		roomUuid:              roomUuid,
		gameId:                gameId,
		mathModuleId:          mathModuleId,
		in:                    make(chan *flux.Action, common.Service.ChanSize),
		scriptsScene:          scriptsScene,
		scriptsSceneResetTime: scriptsSceneResetTime,
		goroutinePool:         make(map[chan bool]bool),
		isReady:               make(chan bool, 1024),
	}
	go f.run()

	flux.Send(ActionChangeNextScene, f.hashId(), f.roomUuid, 0)
	logger.Service.Zap.Infow("Start Fish Pool",
		"GameRoomUuid", roomUuid,
		"Chan", fmt.Sprintf("%p", f.in),
	)
	return f
}

func (f *fishPool) run() {
	flux.Register(f.hashId(), f.in)
	fishInPool := make(map[string]*Fish)
	f.isReady <- true

	for {
		select {
		case action, ok := <-f.in:
			if !ok {
				return
			}
			f.handleAction(action, fishInPool)
		}
	}
}

func (f *fishPool) handleAction(action *flux.Action, fishInPool map[string]*Fish) {
	switch action.Key().Name() {
	case actionFishIn:
		if action.Key().To() == f.hashId() {
			fish := action.Payload()[0].(*Fish)
			dataBuffer := action.Payload()[1]
			fishInPool[fish.FishUuid] = fish
			flux.Send(EMSGID_eBroadcastFishIn, f.hashId(), f.roomUuid, dataBuffer)
		}

	case actionFishOut:
		if action.Key().To() == f.hashId() {
			fishUuid := action.Payload()[0].(string)

			if f, ok := fishInPool[fishUuid]; ok {
				f.destroy()
				delete(fishInPool, fishUuid)
			}
		}

	case actionGroupIn:
		if action.Key().To() == f.hashId() {
			groupFish := action.Payload()[0].([]*Fish)
			dataBuffer := action.Payload()[1]

			for _, v := range groupFish {
				fishInPool[v.FishUuid] = v
			}
			flux.Send(EMSGID_eBroadcastGroupIn, f.hashId(), f.roomUuid, dataBuffer)
		}

	case actionFishStart:
		if action.Key().To() == f.hashId() {
			f.cleanFishPool(fishInPool)
			f.cleanGoroutinePool()
			f.startScriptScene(0)
		}

	case actionFishReset:
		if action.Key().To() == f.hashId() {
			scene := action.Payload()[0].(int)
			f.eRoundEnd(scene + 1)
			f.cleanGoroutinePool()
			f.cleanFishPool(fishInPool)
			time.Sleep(5 * time.Second)
			f.startScriptScene(scene + 1)
		}

	case actionFishDestroy:
		if action.Key().To() == f.hashId() {
			f.cleanGoroutinePool()
			f.cleanFishPool(fishInPool)
			flux.UnRegister(f.hashId(), f.in)
		}

	case ActionFishExistCheck:
		hitFish := action.Payload()[0].(*Fish)
		mercenaryType := action.Payload()[1].(int32)
		blacklistValue := action.Payload()[2].(uint32)
		depositMultiple := action.Payload()[3].(uint64)

		if action.Key().To() == f.hashId() {
			if v, ok := fishInPool[hitFish.FishUuid]; ok {
				if v.TypeId == hitFish.TypeId {
					hitFish.Scene = v.Scene
					flux.Send(ActionFishExist, f.hashId(), action.Key().From(), hitFish, mercenaryType, blacklistValue, depositMultiple)
				} else {
					logger.Service.Zap.Warnw("Can't find fish type.",
						"GameRoomUuid", hitFish.RoomUuid,
						"FishUuid", hitFish.FishUuid,
						"FishType", hitFish.TypeId,
					)
					flux.Send(ActionFishNotExist, f.hashId(), action.Key().From(), hitFish, mercenaryType, blacklistValue, depositMultiple)
				}
			} else {
				logger.Service.Zap.Warnw("Can't find fish UUID.",
					"GameRoomUuid", hitFish.RoomUuid,
					"FishUuid", hitFish.FishUuid,
					"FishType", hitFish.TypeId,
				)
				flux.Send(ActionFishNotExist, f.hashId(), action.Key().From(), hitFish, mercenaryType, blacklistValue, depositMultiple)
			}
		}

	case ActionMultiFishExistCheck:
		hitFishies := action.Payload()[0].(*Fishies)
		mercenaryType := action.Payload()[1].(int32)
		blacklistValue := action.Payload()[2].(uint32)
		depositMultiple := action.Payload()[3].(uint64)
		checkFishiesExist := make([]bool, len(hitFishies.Fishies))

		if action.Key().To() == f.hashId() {
			for i := 0; i < len(hitFishies.Fishies); i++ {
				if v, ok := fishInPool[hitFishies.Fishies[i].FishUuid]; ok {
					if v.TypeId == hitFishies.Fishies[i].TypeId {
						hitFishies.Fishies[i].Scene = v.Scene
						checkFishiesExist[i] = true
					} else {
						logger.Service.Zap.Warnw("Can't find fish type",
							"GameRoomUuid", hitFishies.Fishies[i].RoomUuid,
							"FishUuid", hitFishies.Fishies[i].FishUuid,
							"FishType", hitFishies.Fishies[i].TypeId,
						)
						checkFishiesExist[i] = false
					}
				} else {
					logger.Service.Zap.Warnw("Can't find fish UUID",
						"GameRoomUuid", hitFishies.Fishies[i].RoomUuid,
						"FishUuid", hitFishies.Fishies[i].FishUuid,
						"FishType", hitFishies.Fishies[i].TypeId,
					)
					checkFishiesExist[i] = false
				}
			}
		}

		flux.Send(ActionMultiFishiesCheck, f.hashId(), action.Key().From(), hitFishies, checkFishiesExist, mercenaryType, blacklistValue, depositMultiple)

	case EMSGID_eAllFishCall:
		secWebSocketKey := action.Key().From()
		data := &fish_proto.AllFishCall{}

		if err := proto.Unmarshal(action.Payload()[0].([]byte), data); err != nil {
			logger.Service.Zap.Errorw(Fish_ALL_FISH_CALL_PROTO_INVALID,
				"GameUser", secWebSocketKey,
				"GameRoomUuid", data.RoomUuid,
			)
			errorcode.Service.Fatal(secWebSocketKey, Fish_ALL_FISH_CALL_PROTO_INVALID)
			return
		}

		if f.roomUuid == data.RoomUuid {
			f.eAllFishRecall(secWebSocketKey, fishInPool)
		}

	case ActionFishDie:
		if action.Key().To() == f.hashId() {
			hitFish := action.Payload()[0].(*Fish)

			if f, ok := fishInPool[hitFish.FishUuid]; ok {
				f.die <- true
				f.destroy()
				delete(fishInPool, hitFish.FishUuid)
			}
		}
	}
}

func (f *fishPool) cleanFishPool(fishInPool map[string]*Fish) {
	for fishUuid, v := range fishInPool {
		v.out <- true
		v.destroy()
		delete(fishInPool, fishUuid)
	}
}

func (f *fishPool) cleanGoroutinePool() {
	for destroy := range f.goroutinePool {
		destroy <- true
		close(destroy)
		delete(f.goroutinePool, destroy)
	}
}

func (f *fishPool) startScriptScene(scene int) {
	if scene >= len(f.scriptsScene) {
		flux.Send(actionFishStart, f.hashId(), f.hashId(), 0)
		return
	}

	f.eRoundStart()

	destroy := make(chan bool, 1024)
	f.goroutinePool[destroy] = true
	go f.runScriptsSceneResetTime(scene, destroy)
	f.runScriptsScene(scene)
}

func (f *fishPool) runScriptsSceneResetTime(scene int, destroy chan bool) {
	t := time.NewTimer(time.Duration(f.scriptsSceneResetTime[scene]) * time.Second)

	defer t.Stop()

	for {
		select {
		case <-t.C:
			logger.Service.Zap.Debugw("Fish Pool Script Reset",
				"GameRoomUuid", f.roomUuid,
				"Scene", scene,
			)
			flux.Send(actionFishReset, f.hashId(), f.hashId(), scene)
			return

		case is, ok := <-destroy:
			if !ok || is {
				return
			}
		}
	}
}

func (f *fishPool) runScriptsScene(scene int) {
	ss := f.scriptsScene[scene]

	for i := 1; i <= len(ss); i++ {
		destroy := make(chan bool, 1024)
		f.goroutinePool[destroy] = true

		if v, ok := ss["Script_"+strconv.Itoa(i)]; ok {
			go f.runScriptEpisode(v.(map[string]interface{}), scene, destroy)
		}
	}
}

func (f *fishPool) runScriptEpisode(scriptData map[string]interface{}, scene int, destroy chan bool) {
	t := time.NewTimer(time.Duration(int(scriptData["Time"].(float64))) * time.Second)

	defer t.Stop()

	for {
		select {
		case <-t.C:
			var fishPath []interface{}
			var rtpId = ""

			switch f.gameId {
			case models.PSF_ON_00001, models.PGF_ON_00001:
				fishPath = PSF_ON_00001.FishPath.Data(int(scriptData["PathID"].(float64)))

			case models.PSF_ON_00002, models.PSF_ON_20002:
				fishPath = PSF_ON_00002.FishPath.Data(int(scriptData["PathID"].(float64)))

			case models.PSF_ON_00003:
				fishPath = PSF_ON_00003.FishPath.Data(int(scriptData["PathID"].(float64)))

			case models.PSF_ON_00004:
				fishPath = PSF_ON_00004.FishPath.Data(int(scriptData["PathID"].(float64)))
				rtpId = scriptData["RTP"].(string)

			case models.PSF_ON_00005:
				fishPath = PSF_ON_00005.FishPath.Data(int(scriptData["PathID"].(float64)))

			case models.PSF_ON_00006:
				fishPath = PSF_ON_00006.FishPath.Data(int(scriptData["PathID"].(float64)))
				rtpId = scriptData["RTP"].(string)

			case models.PSF_ON_00007:
				fishPath = PSF_ON_00007.FishPath.Data(int(scriptData["PathID"].(float64)))
				rtpId = scriptData["RTP"].(string)

			case models.RKF_H5_00001:
				fishPath = RKF_H5_00001.FishPath.Data(int(scriptData["PathID"].(float64)))

			default:
				logger.Service.Zap.Errorw(Fish_GAME_ID_INVALID,
					"GameRoomUuid", f.roomUuid,
					"Chan", fmt.Sprintf("%p", f.in),
					"GameId", f.gameId,
					"MathModuleId", f.mathModuleId,
				)
				errorcode.Service.Fatal(f.roomUuid, Fish_GAME_ID_INVALID)
				return
			}

			if scriptData["IsGroup"].(bool) {
				f.fishGroup(fishPath, int(scriptData["GroupID"].(float64)), scene, rtpId)
			} else {
				f.fishSingle(fishPath, int(scriptData["FishID"].(float64)), scene, rtpId)
			}
			return

		case is, ok := <-destroy:
			if !ok || is {
				return
			}
		}
	}
}

func (f *fishPool) fishGroup(fishPath []interface{}, groupId, scene int, rtpId string) {
	var group map[string]interface{}
	var fishId int

	switch f.gameId {
	case models.PSF_ON_00001, models.PGF_ON_00001:
		group, fishId = PSF_ON_00001.Groups.Data(groupId)

	case models.PSF_ON_00002, models.PSF_ON_20002:
		group, fishId = PSF_ON_00002.Groups.Data(groupId)

	case models.PSF_ON_00003:
		group, fishId = PSF_ON_00003.Groups.Data(groupId)

	case models.PSF_ON_00004:
		group, fishId = PSF_ON_00004.Groups.Data(groupId)

	case models.PSF_ON_00005:
		group, fishId = PSF_ON_00005.Groups.Data(groupId)

	case models.PSF_ON_00006:
		group, fishId = PSF_ON_00006.Groups.Data(groupId)

	case models.PSF_ON_00007:
		group, fishId = PSF_ON_00007.Groups.Data(groupId)

	case models.RKF_H5_00001:
		group, fishId = RKF_H5_00001.Groups.Data(groupId)

	default:
		logger.Service.Zap.Errorw(Fish_GAME_ID_INVALID,
			"GameRoomUuid", f.roomUuid,
			"Chan", fmt.Sprintf("%p", f.in),
			"GameId", f.gameId,
			"MathModuleId", f.mathModuleId,
		)
		errorcode.Service.Fatal(f.roomUuid, Fish_GAME_ID_INVALID)
		return
	}

	if group != nil && fishId == -1 {
		// 目前只有[PSF_ON_00004]捕魚大排檔有在使用
		switch f.gameId {
		case models.PSF_ON_00004, models.PSF_ON_00006, models.PSF_ON_00007:
			fishId := group["FishID"].(float64)
			fishAmount := int(group["FishAmount"].(float64))

			// 一群魚(群游)
			if f.getGroupType(f.gameId, groupId) == 1 {
				separation := int32(group["Separation"].(float64))
				f.eBroadcastGroupIn(fishPath, int32(fishId), fishAmount, scene, separation, rtpId)
			}

			// 一群魚(單筆魚資訊給)
			if f.getGroupType(f.gameId, groupId) == 0 {
				interval := int(group["Interval"].(float64) * 1000)
				for i := 0; i < fishAmount; i++ {
					f.fishSingle(fishPath, int(fishId), scene, rtpId)
					time.Sleep(time.Duration(interval) * time.Millisecond)
				}
			}

		// PSF_ON_00001, PSF_ON_00002, PSF_ON_00003, PSF_ON_00005
		default:
			fishId := group["FishID"].(float64)
			fishAmount := int(group["FishAmount"].(float64))
			interval := int(group["Interval"].(float64) * 1000)

			for i := 0; i < fishAmount; i++ {
				f.fishSingle(fishPath, int(fishId), scene, rtpId)
				time.Sleep(time.Duration(interval) * time.Millisecond)
			}
		}
	}

	if group == nil && fishId > -1 {
		f.fishSingle(fishPath, fishId, scene, rtpId)
	}
}

func (f *fishPool) fishSingle(fishPath []interface{}, fishId, scene int, rtpId string) {
	f.eBroadcastFishIn(fishPath, int32(fishId), scene, rtpId)
}

func (f *fishPool) eRoundStart() {
	rs := &fish_proto.RoundStart{
		Msgid:    common_proto.EMSGID_eRoundStart,
		RoomUuid: f.roomUuid,
	}
	dataBuffer, _ := proto.Marshal(rs)
	flux.Send(EMSGID_eBroadcastRoundStart, f.hashId(), f.roomUuid, dataBuffer)
}

func (f *fishPool) eRoundEnd(scene int) {
	if scene >= len(f.scriptsScene) {
		scene = 0
	}

	data := &fish_proto.RoundEnd{
		Msgid:    common_proto.EMSGID_eRoundEnd,
		RoomUuid: f.roomUuid,
	}

	var fishScene int

	switch f.gameId {
	case models.PSF_ON_00001, models.PGF_ON_00001:
		// 3 scenes with 9 scripts
		data.NextScene = int32(scene % 3)
		fishScene = scene % 3

	case models.PSF_ON_00002, models.PSF_ON_20002:
		// 3 scenes with 12 scripts
		data.NextScene = int32(scene % 3)
		fishScene = scene % 3

	case models.PSF_ON_00003:
		// 3 scenes with 9 scripts
		data.NextScene = int32(scene % 3)
		fishScene = scene % 3

	case models.PSF_ON_00004:
		// 3 scenes with 12 scripts
		data.NextScene = int32(scene % 3)
		fishScene = scene % 3

	case models.PSF_ON_00005:
		// 3 scenes with 12 scripts
		data.NextScene = int32(scene % 3)
		fishScene = scene % 3

	case models.PSF_ON_00006, models.PSF_ON_00007:
		// 3 scenes with 12 scripts
		data.NextScene = int32(scene % 3)
		fishScene = scene % 3

	case models.RKF_H5_00001:
		// 3 scenes with 12 scripts
		data.NextScene = int32(scene % 3)
		fishScene = scene % 3

	default:
		logger.Service.Zap.Errorw(Fish_GAME_ID_INVALID,
			"GameRoomUuid", f.roomUuid,
			"Chan", fmt.Sprintf("%p", f.in),
			"GameId", f.gameId,
			"MathModuleId", f.mathModuleId,
		)
		errorcode.Service.Fatal(f.roomUuid, Fish_GAME_ID_INVALID)
		return
	}

	dataBuffer, _ := proto.Marshal(data)
	flux.Send(EMSGID_eBroadcastRoundEnd, f.hashId(), f.roomUuid, dataBuffer)
	flux.Send(ActionChangeNextScene, f.hashId(), f.roomUuid, fishScene)
}

func (f *fishPool) eBroadcastGroupIn(fishPath []interface{}, fishId int32, fishAmount, scene int, separation int32, rtpId string) {
	var fishLife int32
	var fishExtraData interface{}
	var fishScene int
	groupUUID := uuid.New().String()

	switch f.gameId {
	case models.PSF_ON_00004:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00004.Objects, int(fishId), scene, 3)

	case models.PSF_ON_00006:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00006.Objects, int(fishId), scene, 3)

	case models.PSF_ON_00007:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00007.Objects, int(fishId), scene, 3)

	default:
		logger.Service.Zap.Errorw(Fish_GAME_ID_INVALID,
			"GameRoomUuid", f.roomUuid,
			"Chan", fmt.Sprintf("%p", f.in),
			"GameId", f.gameId,
			"MathModuleId", f.mathModuleId,
		)
		errorcode.Service.Fatal(f.roomUuid, Fish_GAME_ID_INVALID)
		return
	}

	bi := f.bezierInfo(fishPath, fishLife)
	var groupFish []*fish_proto.FishInfo
	var groupFishInfo []*Fish

	for i := 0; i < fishAmount; i++ {
		fishUuid := uuid.New().String()

		fish := New(
			f.gameId,
			f.roomUuid,
			fishUuid,
			f.mathModuleId,
			-1,
			fishId,
			fishLife,
			fishScene,
			rtpId,
			fishExtraData,
			[]*fish_proto.BezierInfo{bi},
			groupUUID,
			separation,
		)
		fish.skills()

		fish.mutex.Lock()
		fi := &fish_proto.FishInfo{
			Uuid:        fishUuid,
			SymbolId:    fish.TypeId,
			Life:        fishLife,
			RtpId:       rtpId,
			ElapsedTime: 0,
			Group: &fish_proto.GroupInfo{
				Uuid:       groupUUID,
				Separation: separation,
			},
		}
		fish.mutex.Unlock()

		fi.Path = append(fi.Path, bi)
		go f.eBroadcastFishOut(fish)

		groupFish = append(groupFish, fi)
		groupFishInfo = append(groupFishInfo, fish)
	}

	data := &fish_proto.BroadcastGroupIn{
		Msgid:     common_proto.EMSGID_eBroadcastGroupIn,
		RoomUuid:  f.roomUuid,
		GroupFish: groupFish,
	}
	dataBuffer, _ := proto.Marshal(data)
	flux.Send(actionGroupIn, f.hashId(), f.hashId(), groupFishInfo, dataBuffer)
}

func (f *fishPool) getObjectData(object interface{}, fishId, scene, sceneSize int) (fishLife int32, fishExtraData interface{}, fishScene int) {
	objectData := reflect.ValueOf(object)
	inputValue := []reflect.Value{reflect.ValueOf(fishId)}

	speed := objectData.MethodByName("Speed").Call(inputValue)
	fishLife = speed[0].Interface().(int32)

	fishExtraData = objectData.MethodByName("ExtraData").Call(inputValue)[0].Interface()
	fishScene = scene % sceneSize

	return fishLife, fishExtraData, fishScene
}

func (f *fishPool) eBroadcastFishIn(fishPath []interface{}, fishId int32, scene int, rtpId string) {
	fishUuid := uuid.New().String()
	var fishLife int32
	var fishExtraData interface{}
	var fishScene int

	switch f.gameId {
	case models.PSF_ON_00001, models.PGF_ON_00001:
		// 3 scenes with 9 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00001.Objects, int(fishId), scene, 3)

	case models.PSF_ON_00002, models.PSF_ON_20002:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00002.Objects, int(fishId), scene, 3)

	case models.PSF_ON_00003:
		// 3 scenes with 9 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00003.Objects, int(fishId), scene, 3)

	case models.PSF_ON_00004:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00004.Objects, int(fishId), scene, 3)

	case models.PSF_ON_00005:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00005.Objects, int(fishId), scene, 3)

	case models.PSF_ON_00006:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00006.Objects, int(fishId), scene, 3)

	case models.PSF_ON_00007:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(PSF_ON_00007.Objects, int(fishId), scene, 3)

	case models.RKF_H5_00001:
		// 3 scenes with 12 scripts
		fishLife, fishExtraData, fishScene = f.getObjectData(RKF_H5_00001.Objects, int(fishId), scene, 3)

	default:
		logger.Service.Zap.Errorw(Fish_GAME_ID_INVALID,
			"GameRoomUuid", f.roomUuid,
			"Chan", fmt.Sprintf("%p", f.in),
			"GameId", f.gameId,
			"MathModuleId", f.mathModuleId,
		)
		errorcode.Service.Fatal(f.roomUuid, Fish_GAME_ID_INVALID)
		return
	}

	bi := f.bezierInfo(fishPath, fishLife)

	fish := New(
		f.gameId,
		f.roomUuid,
		fishUuid,
		f.mathModuleId,
		-1,
		fishId,
		fishLife,
		fishScene,
		rtpId,
		fishExtraData,
		[]*fish_proto.BezierInfo{bi},
		"",
		-1,
	)
	fish.skills()

	fish.mutex.Lock()
	fi := &fish_proto.FishInfo{
		Uuid:        fishUuid,
		SymbolId:    fish.TypeId,
		Life:        fishLife,
		RtpId:       rtpId,
		ElapsedTime: 0,
		Group: &fish_proto.GroupInfo{
			Uuid: "",
		},
	}
	fish.mutex.Unlock()

	fi.Path = append(fi.Path, bi)

	go f.eBroadcastFishOut(fish)

	data := &fish_proto.BroadcastFishIn{
		Msgid:    common_proto.EMSGID_eBroadcastFishIn,
		RoomUuid: f.roomUuid,
		Fish:     fi,
	}

	dataBuffer, _ := proto.Marshal(data)
	flux.Send(actionFishIn, f.hashId(), f.hashId(), fish, dataBuffer)
}

func (f *fishPool) bezierInfo(path []interface{}, fishLife int32) *fish_proto.BezierInfo {
	return &fish_proto.BezierInfo{
		Time: fishLife,
		StartPoint: &fish_proto.BezierInfo_Point{
			X: int32(path[0].(float64)),
			Y: int32(path[1].(float64)),
		},
		EndPoint: &fish_proto.BezierInfo_Point{
			X: int32(path[6].(float64)),
			Y: int32(path[7].(float64)),
		},
		Op1Point: &fish_proto.BezierInfo_Point{
			X: int32(path[2].(float64)),
			Y: int32(path[3].(float64)),
		},
		Op2Point: &fish_proto.BezierInfo_Point{
			X: int32(path[4].(float64)),
			Y: int32(path[5].(float64)),
		},
	}
}

func (f *fishPool) eBroadcastFishOut(fish *Fish) {
	data := &fish_proto.BroadcastFishOut{
		Msgid:    common_proto.EMSGID_eBroadcastFishOut,
		RoomUuid: f.roomUuid,
		Uuid:     fish.FishUuid,
	}

	t := time.NewTimer(time.Duration(fish.fishLife) * time.Second)

	defer t.Stop()

	for {
		select {
		case is, ok := <-fish.out:

			// We needn't send EMSGID_eBroadcastFishOut when actionFishReset.
			if !ok || is {
				logger.Service.Zap.Debugw("Fish Reset",
					"GameRoomUuid", f.roomUuid,
					"FishUuid", fish.FishUuid,
				)
				return
			}

		case is, ok := <-fish.die:
			if !ok {
				return
			}

			if is {
				data.Type = 1
				dataByte, _ := proto.Marshal(data)

				flux.Send(EMSGID_eBroadcastFishOut, f.hashId(), f.roomUuid, dataByte)

				logger.Service.Zap.Infow("Fish Died",
					"GameRoomUuid", f.roomUuid,
					"FishUuid", fish.FishUuid,
				)
				return
			}

		case <-t.C:
			data.Type = 0
			flux.Send(actionFishOut, f.hashId(), f.hashId(), fish.FishUuid)
			dataBuffer, _ := proto.Marshal(data)
			flux.Send(EMSGID_eBroadcastFishOut, f.hashId(), f.roomUuid, dataBuffer)

			logger.Service.Zap.Debugw("Fish Timeout",
				"GameRoomUuid", f.roomUuid,
				"FishUuid", fish.FishUuid,
			)
			return
		}
	}
}

func (f *fishPool) eAllFishRecall(secWebSocketKey string, fishInPool map[string]*Fish) {
	data := &fish_proto.AllFishRecall{
		Msgid:    common_proto.EMSGID_eAllFishRecall,
		RoomUuid: f.roomUuid,
	}

	for k, v := range fishInPool {
		fi := &fish_proto.FishInfo{
			Uuid:        k,
			SymbolId:    v.TypeId,
			Life:        v.fishLife,
			ElapsedTime: time.Now().Unix() - v.createTime,
			Path:        v.path,
			Group: &fish_proto.GroupInfo{
				Uuid:       v.group.GroupId,
				Separation: v.group.Separation,
			},
		}
		data.AllFish = append(data.AllFish, fi)
	}

	dataByte, _ := proto.Marshal(data)
	flux.Send(EMSGID_eAllFishRecall, f.hashId(), secWebSocketKey, dataByte)
}

func (f *fishPool) hashId() string {
	return Service.HashId(f.roomUuid)
}

func (f *fishPool) destroy() {
	logger.Service.Zap.Infow("Fish Pool Destroy",
		"GameRoomUuid", f.roomUuid,
		"Chan", fmt.Sprintf("%p", f.in),
		"GameId", f.gameId,
		"MathModuleId", f.mathModuleId,
	)

	flux.Send(actionFishDestroy, f.hashId(), f.hashId(), "")
	close(f.isReady)
}

func (f *fishPool) getGroupType(gameId string, groupId int) (groupType int32) {
	switch gameId {
	case models.PSF_ON_00004:
		return PSF_ON_00004.Groups.GroupType(groupId)
	case models.PSF_ON_00006:
		return PSF_ON_00006.Groups.GroupType(groupId)
	case models.PSF_ON_00007:
		return PSF_ON_00007.Groups.GroupType(groupId)
	default:
		return 0
	}
}
