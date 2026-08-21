package api

import (
	"fmt"
	"game/comm"
	"game/comm/define"
	"game/duck/ut2/jwtutil"
	"strconv"
)

// https://api.pg-demo.com/web-api/auth/session/v2/verifyOperatorPlayerSession?traceId=VXTNFQ12
func verifyOperatorPlayerSession(ps *PGParams, ret *M) (err error) {
	os := ps.Form.Get("os")
	// slog.Info("verifyOperatorPlayerSession", "os", os)
	pid, err := jwtutil.ParseToken(os)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	item := comm.GetCurrentItem(pid)
	cc, cs := item.Key, item.Symbol
	if item.CurrencyManufactureVisibleOff != nil {
		if _, ok := item.CurrencyManufactureVisibleOff[comm.PG]; ok {
			if item.CurrencyManufactureVisibleOff[comm.PG] == 0 { //关闭
				cc, cs = "", ""
			} else { //开启
			}
		}
	}
	// cs = "R"
	params := comm.GetEXParams(pid, ps.GameId)
	gid := ps.GetInt("gi")
	pidstr := strconv.Itoa(int(pid))
	geu := fmt.Sprintf("game-api/%v/", ps.Get("gi"))
	tk := os
	bf := 1                    //购买fg
	if ps.GameId == "pg_104" { //亡灵特殊处理
		bf = 0
	}
	//as字段控制自动旋转设为1000
	s := M{
		"oj":    M{"jid": 0},
		"pid":   pidstr,
		"pcd":   pidstr,
		"tk":    tk,
		"st":    1,
		"geu":   geu,
		"lau":   "game-api/lobby/",
		"bau":   "web-api/game-proxy/",
		"cc":    cc,
		"cs":    cs,
		"nkn":   "123456",
		"gm":    D{M{"gid": gid, "msdt": 1734752727375, "medt": 1744752727375, "st": 1, "amsg": "", "rtp": M{"df": M{"min": 96.22, "max": 96.22}}, "mxe": 2500, "meshr": 8960913}},
		"uiogc": M{"bb": 1, "grtp": 1, "gec": 0, "cbu": 0, "cl": 0, "bf": bf, "mr": 0, "phtr": 0, "vc": 0, "bfbsi": 1, "bfbli": 1, "il": 0, "rp": 0, "gc": params.ShowNameAndTimeOff, "ign": params.ShowNameAndTimeOff, "tsn": 0, "we": 0, "gsc": 0, "bu": 0, "pwr": 0, "hd": 0, "et": 0, "np": 0, "igv": 0, "as": 1000, "asc": params.StopLoss, "std": 0, "hnp": 0, "ts": 0, "smpo": 0, "grt": 0, "ivs": 1, "ir": 0, "hn": 1, "gvs": 0, "hwl": 1, "rcf": 0, "sbb": 1, "sp": params.CarouselOff, "swf": 0, "swfbli": 0, "swfbsi": 0},
		"ec":    D{},
		"occ":   M{"rurl": "", "tcm": "", "tsc": 0, "ttp": 0, "tlb": "", "trb": ""},
		"gcv":   "1.1.0.8",
		"ioph":  "9411a2022187",
		"jc":    M{"as": 1000, "asc": 1, "bf": 0, "et": 0, "hnp": 0, "hwl": 1, "np": 0, "rcf": 0, "sbb": 1, "smpo": 0, "sp": params.CarouselOff, "std": 0, "swf": 0, "ts": 1},
	}

	// btt=1&vc=2&pf=1&l=th&gi=39&os=000102030404e44708090b853da47f89&otk=c6d81c1d8bfb0e214f632e6185f11e71
	// s := M{"oj": M{"jid": 1}, "pid": "seJKemVLhF", "pcd": "123456", "tk": "7E99995F-4BF6-4CD8-9796-AC49E336BDB1", "st": 1, "geu": "https://api.kafa010.com/game-api/piggy-gold/", "lau": "https://api.kafa010.com/game-api/lobby/", "bau": "https://api.kafa010.com/web-api/game-proxy/", "cc": "PGC", "cs": "", "nkn": "123456", "gm": D{M{"gid": 39, "msdt": 1569831980000, "medt": 1569832080000, "st": 1, "amsg": ""}}, "uiogc": M{"bb": 1, "grtp": 0, "gec": 1, "cbu": 0, "cl": 0, "bf": 0, "mr": 0, "phtr": 0, "vc": 0, "bfbsi": 0, "bfbli": 0, "il": 0, "rp": 0, "gc": 0, "ign": 0, "tsn": 0, "we": 0, "gsc": 1, "bu": 0, "pwr": 0, "hd": 0, "et": 0, "np": 0, "igv": 0, "as": 0, "asc": 0, "std": 0, "hnp": 0, "ts": 0, "smpo": 0, "grt": 0, "ivs": 0, "ir": 0, "gvs": 0, "hn": 1}, "ec": D{M{"n": "132bb011e7", "v": "10", "il": 0, "om": 0, "uie": M{"ct": "1"}}, M{"n": "5e3d8c75c3", "v": "6", "il": 0, "om": 0, "uie": M{"ct": "1"}}}, "occ": M{"rurl": "", "tcm": "You are playing Demo.", "tsc": 10, "ttp": 300, "tlb": "Continue", "trb": "Quit"}, "gcv": "1.1.0.8", "ioph": "fb9e681fb2e9"}

	// tk := ps.Form.Get("tk")
	// TODO gen token
	// apiurlbase := gamedata.Get().ApiUrlBase

	// tk := strings.ToUpper(uuid.NewString())
	// tk = os
	// sessionMng.Set(tk, 123)
	// s["tk"] = tk
	// gi, _ := strconv.Atoi(ps.Form.Get("gi"))
	// s["gm"] = D{M{"gid": gi, "msdt": 1569831980000, "medt": 1569832080000, "st": 1, "amsg": ""}}
	// s["geu"] = apiurlbase + "/game-api/piggy-gold/"
	// s["geu"] = apiurlbase + fmt.Sprintf("/game-api/%v/", ps.Get("gi"))
	// s["lau"] = apiurlbase + "/game-api/lobby/"
	// s["bau"] = apiurlbase + "/web-api/game-proxy/"
	// s["cc"] = "PGC"
	// s["cs"] = "$"
	// s["cc"] = "THB"
	// s["cs"] = "฿"
	// s["pid"] = strconv.Itoa(int(pid))
	// s["uiogc"] = M{"bb": 1, "grtp": 0, "gec": 1, "cbu": 0, "cl": 0, "bf": 1, "mr": 0, "phtr": 0, "vc": 0, "bfbsi": 2, "bfbli": 3, "il": 0, "rp": 0, "gc": 0, "ign": 0, "tsn": 0, "we": 0, "gsc": 1, "bu": 0, "pwr": 0, "hd": 0, "et": 0, "np": 0, "igv": 0, "as": 0, "asc": 0, "std": 0, "hnp": 0, "ts": 0, "smpo": 0, "grt": 0, "ivs": 1, "ir": 0, "gvs": 0, "hn": 1}

	//TODO fetch to game svr
	// s["uiogc"] = M{"bb": 1, "grtp": 0, "gec": 1, "cbu": 0, "cl": 0, "bf": 1, "mr": 0, "phtr": 0, "vc": 0, "bfbsi": 2, "bfbli": 2, "il": 0, "rp": 0, "gc": 0, "ign": 0, "tsn": 0, "we": 0, "gsc": 0, "bu": 0, "pwr": 0, "hd": 0, "et": 0, "np": 0, "igv": 0, "as": 0, "asc": 0, "std": 0, "hnp": 0, "ts": 0, "smpo": 0, "grt": 0, "ivs": 1, "ir": 0, "hn": 1}

	*ret = s

	// slotsmongo.UpdatePlrLoginTime(pid)
	return
}
