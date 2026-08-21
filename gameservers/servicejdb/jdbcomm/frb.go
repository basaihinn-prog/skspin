package jdbcomm

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"serve/comm/mq"
)

type NextFRBPs struct {
	AppID  string
	UserID string
	GameID string

	// Remove this record after fetch success
	Remove bool
}

type FRBPlayer struct {
	ID             primitive.ObjectID `bson:"_id"`
	AppID          string             `bson:"AppID"`
	UserID         string             `bson:"UserID"`
	GameID         string             `bson:"GameID"`
	BonusCode      string             `bson:"BonusCode"`
	TotalBet       float64            `bson:"TotalBet"`
	Rounds         int                `bson:"Rounds"`
	ExpirationDate int64              `bson:"ExpirationDate"`
}

func (frb *FRBPlayer) EVStart(line int) string {
	// "FR1~0.50,20,33.00,,;FR0~1.00,20,10,0,0,1729740027,1,,"
	ev := "FR0~%.2f,%d,%d,0,0,%d,1,,"

	return fmt.Sprintf(ev, frb.TotalBet/float64(line), line, frb.Rounds, frb.ExpirationDate)
}

type NextFRBRet = FRBPlayer

func NextFRB(ps NextFRBPs) (ret *NextFRBRet, err error) {
	var v NextFRBRet
	err = mq.Invoke("/gamecenter/frb/next", ps, &v)
	if err != nil {
		return
	}

	ret = &v
	return
}

//    FreeRound: {
//         Events: "ev",
//         RoundsLeft: "frn",
//         TotalWin: "fra",
//         RoundType: "frt",
//         TurboSpinMode: "tsm",
//         BonusCode: "frbc",
//         BonusCodeParam: "bonus_code",
//         Event: {
//             Start: "FR0",
//             Finish: "FR1",
//             Error: "FR2"
//         }
//     },

/**
VSProtocolParser.ParseVsFreeRoundEvents = function(nameValues) {
    var isFromInit = VSProtocolParser.firstFRBparse;
    if (nameValues[GameProtocolDictionary.FreeRound.Events] == undefined)
        return null;
    var evts = [];
    var items = nameValues[GameProtocolDictionary.FreeRound.Events].split(";");
    for (var i = 0; i < items.length; ++i) {
        var item = items[i].split("~");
        if (item.length > 1) {
            var type = item[0];
            var args = item[1].split(",");
            if (XT.GetBool(Vars.DontShowFRBEndWindowOnInit) && isFromInit && type == GameProtocolDictionary.FreeRound.Event.Finish)
                continue;
            if (type == GameProtocolDictionary.FreeRound.Event.Start || type == GameProtocolDictionary.FreeRound.Event.Finish || type == GameProtocolDictionary.FreeRound.Event.Error) {
                var e = new VsFreeRoundEvent;
                e.Bet = _number.otod(args[0]);
                e.Lines = _number.otoi(args[1]);
                switch (type) {
                case GameProtocolDictionary.FreeRound.Event.Start:
                    e.Type = VsFreeRoundEvent.EventType.Start;
                    e.RoundsLeft = _number.otoi(args[2]);
                    e.TurboSpinMode = _number.otoi(args[3]) == 1 ? true : false;
                    e.PlayLaterAvailable = args.length > 4 && _number.otoi(args[4]) == 1 ? true : false;
                    e.EndDateTimestamp = args.length > 5 ? _number.otod(args[5]) : -1;
                    e.IsFreeRoundPending = args.length > 6 && _number.otoi(args[6]) == 0 ? false : true;
                    e.PromoLocalizedName = args.length > 7 && args[7].length > 0 ? decodeURIComponent(atob(args[7]).split("").map(function(c) {
                        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2)
                    }).join("")) : "";
                    e.PromoType = args.length > 8 && args[8].length > 0 ? args[8] : "";
                    break;
                case GameProtocolDictionary.FreeRound.Event.Finish:
                    e.Type = VsFreeRoundEvent.EventType.Finish;
                    break;
                case GameProtocolDictionary.FreeRound.Event.Error:
                    e.Type = VsFreeRoundEvent.EventType.Error;
                    break
                }
                evts.push(e)
            }
        }
    }
    return evts.length > 0 ? evts : null
}
*/
