package comm

import (
	"go.mongodb.org/mongo-driver/bson"
	"serve/comm/db"
	"serve/comm/slotsmongo"
	"strings"
	"time"
)

type GetBetLogListNewResults struct {
	List  []*slotsmongo.DocBetLogAviator
	Count int64
}

type GetBetLogAviatorListNewParams struct {
	ID               string `bson:"_id"` // optional
	Pid              int64
	AppID            string // optional
	Uid              string // optional
	UserName         string //用户名
	CurrencyKey      string //币种
	RoomId           string
	Bet              int64 // 下注
	Win              int64 // 输赢
	Balance          int64 // 余额
	GameId           string
	CashOutDate      int64     `json:"cashOutDate,omitempty"` // 提现时间
	Payout           float64   `json:"payout"`
	Profit           int64     `json:"profit,omitempty"` // 利润
	RoundBetId       string    `json:"roundBetId,omitempty"`
	RoundId          int64     `json:"roundId,omitempty"`
	MaxMultiplier    float64   `json:"maxMultiplier,omitempty"`
	Frb              bool      // 是否免费
	InsertTime       time.Time `bson:"InsertTime"` // 数据插入时间
	LogType          int8
	FinishType       int8
	ManufacturerName string `json:"manufacturerName" bson:"manufacturerName"`
	CashOutDateSort  int8
	BetId            int
	LastRoundBetId   string

	PageSize  int64
	PageIndex int64
}

func GetBetLogListNew(ps GetBetLogAviatorListNewParams, ret *GetBetLogListNewResults) (err error) {

	query := bson.M{}
	queryStr := "select * from aviatorgamelogs where"
	countStr := "select count(*) from aviatorgamelogs where"
	appidlist := []string{ps.AppID}
	isFirst := true
	argStr := ""
	argList := []any{}

	if ps.Pid > 0 {
		query["Pid"] = ps.Pid
		if isFirst {
			argStr += " Pid = ? "
			isFirst = false
		} else {
			argStr += " and Pid = ? "
		}
		argList = append(argList, ps.Pid)
	}

	if len(ps.AppID) > 0 {
		query["AppID"] = ps.AppID
		if isFirst {
			argStr += " AppID = ? "
			isFirst = false
		} else {
			argStr += " and AppID = ? "
		}
		argList = append(argList, ps.AppID)
	}

	if ps.GameId != "ALL" && ps.GameId != "" {
		if isFirst {
			argStr += " GameID = ?"
			isFirst = false
		} else {
			argStr += " and GameID = ?"
		}
		argList = append(argList, ps.GameId)
	}

	if ps.RoundId != 0 {
		if isFirst {
			argStr += " RoundID = ?"
			isFirst = false
		} else {
			argStr += " and RoundID = ?"
		}
		argList = append(argList, ps.RoundId)
	}

	if ps.LogType != 0 {
		if isFirst {
			argStr += " LogType = ?"
			isFirst = false
		} else {
			argStr += " and LogType = ?"
		}
		argList = append(argList, ps.LogType)
	}

	if ps.FinishType != 0 {
		if isFirst {
			argStr += " FinishType = ?"
			isFirst = false
		} else {
			argStr += " and FinishType = ?"
		}
		argList = append(argList, ps.FinishType)
	}

	if len(ps.ManufacturerName) != 0 {
		if isFirst {
			argStr += " ManufacturerName = ?"
			isFirst = false
		} else {
			argStr += " and ManufacturerName = ?"
		}
		argList = append(argList, ps.ManufacturerName)
	}

	if ps.BetId > 0 {
		query["BetId"] = ps.BetId
		if isFirst {
			argStr += " BetId = ? "
			isFirst = false
		} else {
			argStr += " and BetId = ? "
		}
		argList = append(argList, ps.BetId)
	}

	if ps.Bet != 0 {
		if isFirst {
			argStr += " Bet >= ?"
			isFirst = false
		} else {
			argStr += " and Bet >= ?"
		}
		argList = append(argList, ps.Bet)
	}

	if ps.Win != 0 {
		if isFirst {
			argStr += " Win >= ?"
			isFirst = false
		} else {
			argStr += " and Win >= ?"
		}
		argList = append(argList, ps.Win)
	}

	if ps.Balance != 0 {
		if isFirst {
			argStr += " Balance >= ?"
			isFirst = false
		} else {
			argStr += " and Balance >= ?"
		}
		argList = append(argList, ps.Balance)
	}

	sorts := bson.D{}
	sorts = append(sorts, bson.E{Key: "_id", Value: -1})
	if ps.CashOutDateSort == 1 {
		sorts = append(sorts, bson.E{Key: "CashOutDate", Value: -1})
	}

	//RDate := ps.EndTime + ps.StartTime
	//switch RDate {
	//case 0:
	//case ps.EndTime:
	//	if isFirst {
	//		argStr += " InsertTime <= ?"
	//		isFirst = false
	//	} else {
	//		argStr += " and InsertTime <= ?"
	//	}
	//	argList = append(argList, ps.EndTime)
	//case ps.StartTime:
	//	if isFirst {
	//		argStr += " InsertTime >= ?"
	//		isFirst = false
	//	} else {
	//		argStr += " and InsertTime >= ?"
	//	}
	//	argList = append(argList, ps.StartTime)
	//default:
	//	if isFirst {
	//		argStr += " InsertTime >= ? and InsertTime <= ?"
	//		isFirst = false
	//	} else {
	//		argStr += " and InsertTime >= ? and InsertTime <= ?"
	//	}
	//	argList = append(argList, ps.StartTime, ps.EndTime)
	//}
	conditions := []string{}
	for _, appid := range appidlist {
		argList = append(argList, appid)
		conditions = append(conditions, "?")
	}
	if isFirst {
		argStr += " AppID in (" + strings.Join(conditions, ",") + ")"
	} else {
		argStr += " and AppID in (" + strings.Join(conditions, ",") + ")"
	}

	if ps.LastRoundBetId != "0" {
		argStr += `id < (
    SELECT id FROM aviatorgamelogs WHERE RoundBetId = ? ORDER BY id DESC LIMIT 1`
		argList = append(argList, ps.LastRoundBetId)
	}

	list := []*slotsmongo.DocBetLogAviator{}
	conn, err := db.ClickHouseCollection("")
	if err != nil {
		return
	}
	err = conn.QueryRow(countStr+argStr, argList...).Scan(&ret.Count)
	if err != nil {
		return err
	}
	pagestart := (ps.PageIndex - 1) * ps.PageSize
	pageend := pagestart + ps.PageSize
	if pagestart > ret.Count {
		return
	}

	if pageend > ret.Count {
		pageend = ret.Count
	}
	if err != nil {
		return err
	}
	argList = append(argList, ps.PageSize, (ps.PageIndex-1)*ps.PageSize)
	rows, err := conn.Query(queryStr+argStr+" order by id desc LIMIT ? OFFSET ?", argList...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		betlog := &slotsmongo.DocBetLogAviator{}
		if err = rows.Scan(
			&betlog.ID,
			&betlog.Pid,
			&betlog.AppID,
			&betlog.Uid,
			&betlog.UserName,
			&betlog.CurrencyKey,
			&betlog.RoomId,
			&betlog.Bet,
			&betlog.Win,
			&betlog.Balance,
			&betlog.CashOutDate,
			&betlog.Payout,
			&betlog.Profit,
			&betlog.GameID,
			&betlog.RoundBetId,
			&betlog.RoundId,
			&betlog.MaxMultiplier,
			&betlog.Frb,
			&betlog.InsertTime,
			&betlog.RoundMaxMultiplier,
			&betlog.LogType,
			&betlog.FinishType,
			&betlog.ManufacturerName,
			&betlog.BetId,
		); err != nil {
			return err
		}
		list = append(list, betlog)
	}
	//fmt.Println(err)
	if err != nil {
		return
	}
	ret.List = list

	return err
}
