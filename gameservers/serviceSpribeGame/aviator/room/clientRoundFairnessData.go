package room

import (
	"context"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
)

func ClientRoundFairnessData(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientChangeProfileImage Err", "json.Unmarshal", err)
	//}
	p, _ := data.GetSFSObject("p")
	tempRoundId, _ := p.GetDouble("roundId")
	roundId := int64(tempRoundId)
	coll := db.Collection2(lazy.ServiceName, "fairnessData")
	query := bson.M{
		"roundID": roundId,
		"roomID":  r.Name,
	}
	var fairnessData *FairnessData
	coll.FindOne(context.TODO(), query).Decode(&fairnessData)
	content := &Content{
		C: "roundFairnessResponse",
		P: P{
			Code:     comm.Succ,
			Fairness: fairnessData,
		},
	}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(body)
	marshal, _ := GetClientRoundFairnessDataRsp(&body).ToBinary()
	c.WriteMessage(messageType, marshal)
}
func GetClientRoundFairnessDataRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	fairness := ut.NewSFSObject()
	playerSeeds := ut.NewSFSArray()
	if len(initRsp.Content.P.Fairness.PlayerSeeds) > 0 {
		for _, seed := range initRsp.Content.P.Fairness.PlayerSeeds {
			temp := ut.NewSFSObject()
			temp.PutString("seed", seed.Seed)
			temp.PutString("profileImage", seed.ProfileImage)
			temp.PutString("username", seed.Username)
			playerSeeds.Add(temp, ut.SFS_OBJECT, true)
		}
	}
	fairness.PutSFSArray("playerSeeds", playerSeeds)
	fairness.PutLong("partSeedDecimalNumber", initRsp.Content.P.Fairness.PartSeedDecimalNumber)
	fairness.PutLong("roundStartDate", initRsp.Content.P.Fairness.RoundStartDate)
	fairness.PutInt("roundId", int32(initRsp.Content.P.Fairness.RoundID))
	fairness.PutDouble("result", initRsp.Content.P.Fairness.Result)
	fairness.PutDouble("number", initRsp.Content.P.Fairness.Number)
	fairness.PutString("seedSHA256", initRsp.Content.P.Fairness.SeedSHA256)
	fairness.PutString("partSeedHexNumber", initRsp.Content.P.Fairness.PartSeedHexNumber)
	fairness.PutString("serverSeed", initRsp.Content.P.Fairness.ServerSeed)

	pp.PutSFSObject("fairness", fairness)

	p.PutSFSObject("p", pp)
	p.PutString("c", "roundFairnessResponse")

	so.AddCreatePAC(p, 1, 13)
	return so
}
