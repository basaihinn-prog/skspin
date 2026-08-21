package jdbcomm

//
//import "serve/comm/lazy"
//
//func tryNextFRB(plr *Player, retVars Variables) {
//	frb, _ := NextFRB(NextFRBPs{
//		AppID:  plr.AppID,
//		UserID: plr.Uid,
//		GameID: lazy.ServiceName,
//		Remove: true,
//	})
//
//	if frb != nil {
//		plr.FRBData = &FRBData{
//			Config: frb,
//			Frn:    frb.Rounds,
//		}
//
//		addFRBEV(plr, retVars)
//		setFRBFields(plr, retVars)
//	}
//}
//
//func addFRBEV(plr *Player, retVars Variables) {
//	ev := retVars.Str("ev")
//	if ev != "" {
//		ev += ";"
//	}
//
//	if plr.FRBData.Frn == 0 {
//		ev += plr.FRBData.EVFinish(lazy.Line)
//	} else {
//		ev += plr.FRBData.EVStart(lazy.Line)
//	}
//	retVars.Set("ev", ev)
//}
//
//func setFRBFields(plr *Player, retVars Variables) {
//	if plr.FRBData == nil {
//		return
//	}
//
//	if !retVars.Exist("fra") {
//		retVars.SetFloat("fra", plr.FRBData.Fra)
//	}
//	retVars.SetInt("frn", plr.FRBData.Frn)
//	retVars.Set("frt", "N")
//}
//
//func IsFRBField(field string) bool {
//	return field == "fra" || field == "frn" || field == "frt" || field == "ev"
//}
//
//func FRBOnRoundFinish(plr *Player, retVars Variables, tw float64) {
//	frbdata := plr.FRBData
//	if frbdata != nil {
//		frbdata.Fra += tw
//		setFRBFields(plr, retVars)
//
//		if frbdata.Frn == 0 {
//			addFRBEV(plr, retVars)
//			plr.FRBData = nil
//		} else if frbdata.ReplenishEV {
//			frbdata.ReplenishEV = false
//			addFRBEV(plr, retVars)
//		}
//	}
//	if plr.FRBData == nil {
//		tryNextFRB(plr, retVars)
//	}
//}
//
//func FRBOnInit(plr *Player, retVars Variables) {
//	isEnd, _ := plr.IsEndO()
//	if isEnd {
//		if !plr.FRBData.IsValid() {
//			plr.FRBData = nil
//		}
//
//		if plr.FRBData != nil {
//			addFRBEV(plr, retVars)
//			setFRBFields(plr, retVars)
//		} else {
//			tryNextFRB(plr, retVars)
//		}
//	} else {
//		if plr.FRBData != nil {
//			plr.FRBData.ReplenishEV = true
//		}
//	}
//}
//
//func FRBOnSpin(plr *Player, retVars Variables, isCompleted bool) {
//	if isCompleted {
//		FRBOnRoundFinish(plr, retVars, 0)
//	} else {
//		setFRBFields(plr, retVars)
//	}
//}
