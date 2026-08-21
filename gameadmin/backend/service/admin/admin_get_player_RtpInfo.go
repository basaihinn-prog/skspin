package main

func init() {
	//RegMsgProc("/AdminInfo/CreatePlayerRTP", "新建玩家RTP设置", "AdminInfo", CreatePlayerinfo, CreatePlayerRTP{})
}

//func CreatePlayerinfo(ctx *Context, ps CreatePlayerRTP, ret *comm.Empty) (err error) {
//	// 1.发布地址   /player/setPlayerSettings_%s
//	// 发布结构体
//	//
//	//	type PersonRtpSettings struct {
//	//		RewardPercent  int   `json:"reward_percent"`
//	//		NoAwardPercent int   `json:"no_award_percent"`
//	//		PlayerId       int64 `json:"player_id"`
//	//	}
//	game_Filter := bson.M{}
//	gameid := strings.Split(ps.GameID, ",")
//	game_Filter["_id"] = bson.M{
//		"$in": gameid,
//	}
//	gamelist := []*comm.Game{}
//
//	_ = NewOtherDB("game").Collection("Games").FindAll(game_Filter, &gamelist)
//
//	for _, g := range gamelist {
//
//		rw, nw := GetGameNwandAw(g.ID, float64(ps.ContrllRTP))
//		arw, anw := GetGameNwandAw(g.ID, float64(ps.AutoRemoveRTP))
//		pRTP, err := getBetByWinRTP(ps.AppID, ps.Uid)
//		rTPControl := &comm.PlayerRTPControlModel{
//			AppID:              ps.AppID,
//			GameID:             g.ID,
//			CreateAt:           time.Now(),
//			Uid:                ps.Uid,
//			GameName:           g.Name,
//			PlayerRTP:          pRTP,
//			ControlRTP:         ps.ContrllRTP,
//			RewardPercent:      rw,
//			NoAwardPercent:     nw,
//			AutoRemoveRTP:      ps.AutoRemoveRTP,
//			AutoRewardPercent:  arw,
//			AutoNoAwardPercent: anw,
//			Status:             1,
//		}
//		updata_rTPControl := PlayerRTPControlModel{}
//		err = CollPlayerRTPControl.FindOne(bson.M{"GameID": g.ID, "Uid": ps.Uid}, &updata_rTPControl)
//		//if err != nil {
//		//	return err
//		//}
//		if err != mongo.ErrNoDocuments {
//			settingInfo := bson.M{
//				"CreateAt":           rTPControl.CreateAt,
//				"PlayerRTP":          rTPControl.PlayerRTP,
//				"ControlRTP":         rTPControl.ControlRTP,
//				"RewardPercent":      rTPControl.RewardPercent,
//				"NoAwardPercent":     rTPControl.NoAwardPercent,
//				"AutoRemoveRTP":      rTPControl.AutoRemoveRTP,
//				"AutoRewardPercent":  rTPControl.AutoRewardPercent,
//				"AutoNoAwardPercent": rTPControl.AutoNoAwardPercent,
//				"Status":             rTPControl.Status,
//			}
//			set := bson.M{"$set": settingInfo}
//			err = CollPlayerRTPControl.UpsertOne(bson.M{"GameID": g.ID, "Uid": ps.Uid}, set)
//			if err != nil {
//				return err
//			}
//		} else {
//			id, err := NextID(CollPlayerRTPControl, 0)
//			if err != nil {
//				return err
//			}
//			rTPControl.ID = id
//			err = CollPlayerRTPControl.InsertOne(rTPControl)
//			if err != nil {
//				return err
//			}
//		}
//
//		var prs = PersonRtpSettings{
//			RewardPercent:  int(rw),
//			NoAwardPercent: int(nw),
//			PlayerId:       ps.Uid,
//		}
//		sub := "/player/setPlayerSettings_" + g.ID
//		_ = mq.JsonNC.Publish(sub, prs)
//	}
//	return
//}
