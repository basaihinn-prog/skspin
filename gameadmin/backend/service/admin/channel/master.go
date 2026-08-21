package channel

import (
	"game/comm/mq"
	"log/slog"
)

var TaskQueue = make(chan RtpSettings, 200)

type RtpSettings struct {
	GameId         string
	RewardPercent  int
	NoAwardPercent int
	PlayerId       int64
}
type PersonRtpSettings struct {
	RewardPercent  int   `json:"reward_percent"`
	NoAwardPercent int   `json:"no_award_percent"`
	PlayerId       int64 `json:"player_id"`
}

func PushPlayersetPlayerSettings() {
	for {
		task := <-TaskQueue
		PushMsg(task)
	}
	//for play := range TaskQueue {
	//	//select {
	//	//case :
	//}
}

func PushMsg(task RtpSettings) {
	var per = PersonRtpSettings{
		RewardPercent:  0,
		NoAwardPercent: 0,
		PlayerId:       task.PlayerId,
	}
	sub := "/player/setPlayerSettings_" + task.GameId
	err := mq.PublishMsg(sub, per)
	if err != nil {
		slog.Info("Error push msg /player/setPlayerSettings_ : ", per)
	}
}
