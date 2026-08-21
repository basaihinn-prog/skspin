package api

import (
	"fmt"
	"game/comm"
	"game/comm/define"
	"game/duck/ut2/jwtutil"
	"strconv"
)

// https://api.pg-demo.com/web-api/auth/session/v2/verifySession?traceId=RDVVOC12
type M = map[string]any
type D = []any

// var apiurlbase = "https://api.kafa010.com"

func verifySession(ps *PGParams, ret *M) (err error) {
	// btt=1&vc=0&pf=1&l=th&gi=39&tk=6F1A7757-27A1-4516-9425-C682D35AB052&otk=abcd1234abcd123432531532111kkafa

	tk := ps.Form.Get("tk")
	pid, err := jwtutil.ParseToken(tk)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	// session := sessionMng.Get(tk)
	// if session == nil {
	// 	err = mux.NewErrCode("Invalid session", 1302)
	// 	return
	// }

	// test env
	// M{"oj":M{"jid":1},"pid":"2upgWqlysS","pcd":"testuserzy1","tk":"0ED8607C-986B-4280-9AB2-1B48F67AECD3","st":1,"geu":"https://api.pg-demo.com/game-api/ult-striker/","lau":"https://api.pg-demo.com/game-api/lobby/","bau":"https://api.pg-demo.com/web-api/game-proxy/","cc":"THB","cs":"฿","nkn":"testuserzy1","gm":[{"gid":1489936,"msdt":1695700323000,"medt":1695700323000,"st":1,"amsg":""}],"uiogc":{"bb":1,"grtp":0,"gec":1,"cbu":0,"cl":0,"bf":1,"mr":0,"phtr":0,"vc":0,"bfbsi":2,"bfbli":3,"il":0,"rp":0,"gc":0,"ign":0,"tsn":0,"we":0,"gsc":1,"bu":0,"pwr":0,"hd":0,"et":0,"np":0,"igv":0,"as":0,"asc":0,"std":0,"hnp":0,"ts":0,"smpo":0,"grt":0,"ivs":1,"ir":0,"gvs":0,"hn":1},"ec":[{"n":"132bb011e7","v":"10","il":0,"om":0,"uie":{"ct":"1"}},{"n":"5e3d8c75c3","v":"6","il":0,"om":0,"uie":{"ct":"1"}}],"occ":{"rurl":"","tcm":"You are playing Demo.","tsc":10,"ttp":300,"tlb":"Continue","trb":"Quit"},"ioph":"7b500db83d9f"}

	// prod env
	// {"dt":{"oj":{"jid":1},"pid":"sruaXZqqSl","pcd":"123456","tk":"3331816B-A666-4DD0-916A-57A07C5A8462","st":1,"geu":"game-api/piggy-gold/","lau":"game-api/lobby/","bau":"web-api/game-proxy/","cc":"THB","cs":"฿","nkn":"123456","gm":[{"gid":39,"msdt":1538637872000,"medt":1538637872000,"st":1,"amsg":""}],"uiogc":{"bb":1,"grtp":0,"gec":1,"cbu":0,"cl":0,"bf":0,"mr":0,"phtr":0,"vc":0,"bfbsi":0,"bfbli":0,"il":0,"rp":0,"gc":0,"ign":0,"tsn":0,"we":0,"gsc":0,"bu":0,"pwr":0,"hd":0,"et":0,"np":0,"igv":0,"as":0,"asc":0,"std":0,"hnp":0,"ts":0,"smpo":0,"grt":0,"ivs":1,"ir":0,"hn":1},"ec":[],"occ":{"rurl":"","tcm":"","tsc":0,"ttp":0,"tlb":"","trb":""},"gcv":"1.1.0.8","ioph":"c287917ae070"},"err":null}

	// TODO from config
	// s := M{"oj": M{"jid": 1}, "pid": "0", "pcd": "", "tk": "tokenxxx", "st": 1, "geu": "/game-api/piggy-gold/", "lau": "/game-api/lobby/", "bau": "/web-api/game-proxy/", "cc": "PGC", "cs": "", "nkn": "", "uiogc": M{"bb": 1, "grtp": 0, "gec": 1, "cbu": 0, "cl": 0, "bf": 0, "mr": 0, "phtr": 0, "vc": 0, "bfbsi": 0, "bfbli": 0, "il": 0, "rp": 0, "gc": 0, "ign": 0, "tsn": 0, "we": 0, "gsc": 1, "bu": 0, "pwr": 0, "hd": 0, "et": 0, "np": 0, "igv": 0, "as": 0, "asc": 0, "std": 0, "hnp": 0, "ts": 0, "smpo": 0, "grt": 0, "ivs": 0, "ir": 0,  "hn": 1}, "ec": D{}, "occ": M{"rurl": "", "tcm": "You are playing Demo.", "tsc": 1000000, "ttp": 43200, "tlb": "Continue", "trb": "Quit"}, "gcv": "1.1.0.8", "ioph": "297572482eb4"}

	// cc := "THB"
	// cs := "฿"
	// cc = "PHP"
	// cs = "₱"

	params := comm.GetEXParams(pid, ps.GameId)
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
	bf := 1                    //购买fg
	if ps.GameId == "pg_104" { //亡灵特殊处理
		bf = 0
	}
	gid := ps.GetInt("gi")
	pidstr := strconv.Itoa(int(pid))
	geu := fmt.Sprintf("game-api/%v/", ps.Get("gi"))
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
		"gm":    D{M{"gid": gid, "msdt": 1538637872000, "medt": 1538637872000, "st": 1, "amsg": ""}},
		"uiogc": M{"bb": 1, "grtp": 0, "gec": 1, "cbu": 0, "cl": 0, "bf": bf, "mr": 0, "phtr": 0, "vc": 0, "bfbsi": 0, "bfbli": 0, "il": 0, "rp": 0, "gc": params.ShowNameAndTimeOff, "ign": params.ShowNameAndTimeOff, "tsn": 0, "we": 0, "gsc": 0, "bu": 0, "pwr": 0, "hd": 0, "et": 0, "np": 0, "igv": 0, "as": 1000, "asc": params.StopLoss, "std": 0, "hnp": 0, "ts": 0, "smpo": 0, "grt": 0, "ivs": 1, "ir": 0, "hn": 1, "sp": params.CarouselOff},
		"ec":    D{},
		"occ":   M{"rurl": "", "tcm": "", "tsc": 0, "ttp": 0, "tlb": "", "trb": ""},
		"gcv":   "1.1.0.8",
		"ioph":  "c287917ae070",
		"jc":    M{"as": 1000, "asc": 1, "bf": 0, "et": 0, "hnp": 0, "np": 0, "smpo": 0, "sp": params.CarouselOff, "std": 0, "ts": 1},
	}

	// s["tk"] = tk
	// gi, _ := strconv.Atoi(ps.Form.Get("gi"))
	// s["gm"] = D{M{"gid": gi, "msdt": 1569831980000, "medt": 1569832080000, "st": 1, "amsg": ""}}
	// s["geu"] = apiurlbase + "/game-api/piggy-gold/"
	// s["geu"] = fmt.Sprintf("/game-api/%v/", ps.Get("gi"))
	// s["lau"] = apiurlbase + "/game-api/lobby/"
	// s["bau"] = apiurlbase + "/web-api/game-proxy/"
	// s["cc"] = "THB"
	// s["cs"] = "฿"
	// s["pid"] = strconv.Itoa(int(pid))
	// s["pcd"] = strconv.Itoa(int(pid))
	// s["uiogc"] = M{"bb": 1, "grtp": 0, "gec": 1, "cbu": 0, "cl": 0, "bf": 1, "mr": 0, "phtr": 0, "vc": 0, "bfbsi": 2, "bfbli": 3, "il": 0, "rp": 0, "gc": 0, "ign": 0, "tsn": 0, "we": 0, "gsc": 1, "bu": 0, "pwr": 0, "hd": 0, "et": 0, "np": 0, "igv": 0, "as": 0, "asc": 0, "std": 0, "hnp": 0, "ts": 0, "smpo": 0, "grt": 0, "ivs": 1, "ir": 0,  "hn": 1}

	//TODO fetch to game svr
	//s["uiogc"] = M{"bb": 1, "grtp": 1, "gec": 0, "cbu": 0, "cl": 0, "bf": 1, "mr": 0, "phtr": 0, "vc": 0, "bfbsi": 1, "bfbli": 1, "il": 0, "rp": 0, "gc": 1, "ign": 1, "tsn": 0, "we": 0, "gsc": 0, "bu": 0, "pwr": 0, "hd": 0, "et": 0, "np": 0, "igv": 0, "as": 0, "asc": 0, "std": 0, "hnp": 0, "ts": 0, "smpo": 0, "grt": 0, "ivs": 1, "ir": 0, "hn": 1}

	*ret = s
	return
}
