package fish

import (
	fish_proto "serve/service_fish/domain/fish/proto"
	"serve/service_fish/models"
	"sync"
	"time"

	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/logger"
)

type Fish struct {
	GameId       string
	RoomUuid     string
	FishUuid     string
	MathModuleId string
	SeatId       int32 // for bot used
	TypeId       int32
	fishLife     int32
	Scene        int
	RtpId        string
	createTime   int64
	extraData    interface{}
	path         []*fish_proto.BezierInfo
	die          chan bool
	out          chan bool
	mutex        sync.Mutex
	group        *Group
}

type Fishies struct {
	Fishies []*Fish
}

func New(
	gameId, roomUuid, fishUuid, mathModuleId string,
	seatId, typeId, fishLife int32,
	scene int,
	rtpId string,
	extraData interface{},
	fishPath []*fish_proto.BezierInfo,
	groupId string,
	separation int32) *Fish {
	return &Fish{
		GameId:       gameId,
		RoomUuid:     roomUuid,
		FishUuid:     fishUuid,
		MathModuleId: mathModuleId,
		SeatId:       seatId,
		TypeId:       typeId,
		fishLife:     fishLife,
		Scene:        scene,
		RtpId:        rtpId,
		createTime:   time.Now().Unix(),
		extraData:    extraData,
		path:         fishPath,
		die:          make(chan bool, 1024),
		out:          make(chan bool, 1024),
		mutex:        sync.Mutex{},
		group: &Group{
			GroupId:    groupId,
			Separation: separation,
		},
	}
}

func (f *Fish) skills() {
	switch f.GameId {
	case models.PSF_ON_00001, models.PGF_ON_00001:
		psf_on_00001_Skills(f)

	case models.PSF_ON_00002, models.PSF_ON_20002:
		psf_on_00002_Skills(f)

	case models.PSF_ON_00003:
		psf_on_00003_Skills(f)

	case models.PSF_ON_00004:
		psf_on_00004_Skills(f)

	case models.PSF_ON_00005:
		psf_on_00005_Skills(f)

	case models.PSF_ON_00006:
		psf_on_00006_Skills(f)

	case models.PSF_ON_00007:
		psf_on_00007_Skills(f)

	case models.RKF_H5_00001:
		rkf_h5_00001_Skills(f)

	default:
		logger.Service.Zap.Errorw(Fish_GAME_ID_INVALID,
			"GameRoomUuid", f.RoomUuid,
			"FishUuid", f.FishUuid,
			"FishType", f.TypeId,
		)

		// TODO JOHNNY game room need implement
		errorcode.Service.Fatal(f.RoomUuid, Fish_GAME_ID_INVALID)
	}
}

func (f *Fish) destroy() {
	close(f.out)
	close(f.die)
}
