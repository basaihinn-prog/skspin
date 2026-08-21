package activity

import (
	common_proto "serve/fish_comm/common/proto"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	activity_proto "serve/service_fish/application/activity/proto"
	"serve/service_fish/domain/rank"
	"time"

	"github.com/gogo/protobuf/proto"
)

type activity struct {
	secWebSocketKey string
	hostExtId       string
	hostId          string
	memberId        string
	status          string
	isStop          chan bool
	isGuest         bool
	Uuid            string                   `json:"UUID"`
	Url             string                   `json:"URL"`
	Time            []int64                  `json:"Time"`
	Name            []string                 `json:"Name"`
	MinBet          []uint64                 `json:"MinBet"`
	Limits          []map[string]interface{} `json:"Limits"`
	Kind            string                   `json:"Kind"`
}

func (a *activity) run() {

	now := Service.nowUTC()

	switch {
	case now < a.Time[0]:

		offset := a.Time[0] - now

		logger.Service.Zap.Infow("Activity Status Preview Not Yet",
			"GameUser", a.secWebSocketKey,
			"Uuid", a.Uuid,
			"Time", offset,
		)

		t := time.NewTimer(time.Duration(offset) * time.Millisecond)

		for {
			select {
			case <-t.C:
				logger.Service.Zap.Infow("Activity Status Preview",
					"GameUser", a.secWebSocketKey,
					"Uuid", a.Uuid,
					"Time", offset,
				)

				t.Stop()

				a.status = status_PREVIEW // TODO JOHNNY maybe data race write
				a.EMSGID_eActivity(activity_proto.Status_PREVIEW)

				Service.Check(a.secWebSocketKey, a.hostExtId, a.hostId, a.memberId, a.isGuest)
				return
			}
		}

	case now >= a.Time[0] && now < a.Time[1]:

		offset := a.Time[1] - now

		logger.Service.Zap.Infow("Activity Status Preview",
			"GameUser", a.secWebSocketKey,
			"Uuid", a.Uuid,
			"Time", offset,
		)

		a.status = status_PREVIEW // TODO JOHNNY maybe data race write
		a.EMSGID_eActivity(activity_proto.Status_PREVIEW)

		t := time.NewTimer(time.Duration(offset) * time.Millisecond)

		for {
			select {
			case <-t.C:
				logger.Service.Zap.Infow("Activity Status Start",
					"GameUser", a.secWebSocketKey,
					"Uuid", a.Uuid,
					"Time", offset,
				)

				t.Stop()

				a.status = status_START // TODO JOHNNY maybe data race write
				a.EMSGID_eActivity(activity_proto.Status_START)

				Service.Check(a.secWebSocketKey, a.hostExtId, a.hostId, a.memberId, a.isGuest)
				return
			}
		}

	case now >= a.Time[1] && now < a.Time[2]:

		offset := a.Time[2] - now

		logger.Service.Zap.Infow("Activity Status Start",
			"GameUser", a.secWebSocketKey,
			"Uuid", a.Uuid,
			"Time", offset,
		)

		a.status = status_START // TODO JOHNNY maybe data race write
		a.EMSGID_eActivity(activity_proto.Status_START)

		t := time.NewTimer(time.Duration(offset) * time.Millisecond)

		for {
			select {
			case <-t.C:
				logger.Service.Zap.Infow("Activity Status Stop",
					"GameUser", a.secWebSocketKey,
					"Uuid", a.Uuid,
					"Time", offset,
				)

				t.Stop()

				a.status = status_STOP // TODO JOHNNY maybe data race write
				a.EMSGID_eActivity(activity_proto.Status_STOP)

				Service.Check(a.secWebSocketKey, a.hostExtId, a.hostId, a.memberId, a.isGuest)
				return
			}
		}

	case now >= a.Time[2] && now < a.Time[3]:

		offset := a.Time[3] - now

		logger.Service.Zap.Infow("Activity Status Stop",
			"GameUser", a.secWebSocketKey,
			"Uuid", a.Uuid,
			"Time", offset,
		)

		a.status = status_STOP // TODO JOHNNY maybe data race write
		a.EMSGID_eActivity(activity_proto.Status_STOP)

		t := time.NewTimer(time.Duration(offset) * time.Millisecond)

		for {
			select {
			case <-t.C:
				logger.Service.Zap.Infow("Activity Status Close",
					"GameUser", a.secWebSocketKey,
					"Uuid", a.Uuid,
					"Time", offset,
				)

				t.Stop()

				a.status = status_CLOSE // TODO JOHNNY maybe data race write
				a.EMSGID_eActivity(activity_proto.Status_CLOSE)

				Service.Delete(a.secWebSocketKey)
				Service.Check(a.secWebSocketKey, a.hostExtId, a.hostId, a.memberId, a.isGuest)
				return
			}
		}

	default:
		logger.Service.Zap.Infow("Activity Status Error",
			"GameUser", a.secWebSocketKey,
			"Uuid", a.Uuid,
		)

		a.status = status_CLOSE // TODO JOHNNY maybe data race write
		a.EMSGID_eActivity(activity_proto.Status_CLOSE)

		Service.Delete(a.secWebSocketKey)
	}
}

func (a *activity) EMSGID_eActivity(status activity_proto.Status) {
	dataSettings := &activity_proto.Settings{
		Msgid:      common_proto.EMSGID_eActivity,
		Activities: nil,
	}

	data := &activity_proto.Activity{
		Uuid:   a.Uuid,
		Url:    a.Url + "?host_ext_id=" + a.hostExtId,
		Time:   a.Time,
		Status: status,
	}

	switch a.Kind {
	case rank.KILL_RANKING:
		data.Kind = activity_proto.Kind_KILL_RANKING

	case rank.WIN_RANKING:
		data.Kind = activity_proto.Kind_WIN_RANKING
	}

	dataSettings.Activities = append(dataSettings.Activities, data)

	dataByte, _ := proto.Marshal(dataSettings)
	flux.Send(EMSGID_eActivity, a.Uuid, a.secWebSocketKey, dataByte)
}

func (a *activity) destroy() {
	logger.Service.Zap.Infow("Destroy Activity",
		"GameUser", a.secWebSocketKey,
		"MemberId", a.memberId,
		"HostId", a.hostId,
		"HostExtId", a.hostExtId,
		"Status", a.status,
	)
}
