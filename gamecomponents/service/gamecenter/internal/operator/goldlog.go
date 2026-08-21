package operator

import (
	"context"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"strconv"
	"time"
)

type DocModifyGoldLog struct {
	ID         string    `bson:"_id"`
	Pid        int64     `bson:"pid"`        // 内部id
	InsertTime time.Time `bson:"InsertTime"` // 请求时间

	Change  int64  `bson:"Change"`  // 金币变化
	Balance int64  `bson:"Balance"` // 修改后余额
	Error   string `bson:"Error"`
	Action  string `bson:"Action"`
	ReqData any    `bson:"ReqData"`
}

func InsertModifyLog(doc *DocModifyGoldLog, plr *MemPlr) {
	return
	if len(doc.ID) == 0 {
		sf := ut.NewSnowflake()
		doc.ID = strconv.Itoa(int(sf.NextID()))
	}
	doc.InsertTime = time.Now()
	doc.Pid = plr.Pid

	slotsmongo.GetTTLLogColl("goldlog_"+plr.AppID, 180).InsertOne(context.TODO(), doc)
}

// func InsertModifyLog_old(logid primitive.ObjectID, plr *MemPlr /*req *gamepb.ModifyGoldReq,*/, change int64, comment string, beforeTime time.Time, balance int64, err error) {
// 	oneLog := &DocModifyGoldLog{
// 		ID:       logid,
// 		Pid:      plr.Pid,
// 		Change:   change,
// 		Comment:  comment,
// 		Balance:  balance,
// 		ReqTime:  beforeTime,
// 		RespTime: time.Now(),
// 	}

// 	if err == nil {
// 		oneLog.Status = gamepb.ModifyStatus_OK
// 	} else {
// 		oneLog.Status = gamepb.ModifyStatus_Err
// 		oneLog.Error = err.Error()
// 	}

// 	// gcdb.CollModifyGoldLog.InsertOne(oneLog)

// 	slotsmongo.GetTTLLogColl("goldlog", 30).InsertOne(context.TODO(), oneLog)
// }
