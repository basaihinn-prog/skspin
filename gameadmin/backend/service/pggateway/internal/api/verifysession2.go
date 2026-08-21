package api

import (
	"io"
	"net/http"
)

func verifySession2(w http.ResponseWriter, r *http.Request) {
	retStr := `{"dt":{"oj":{"jid":1},"pid":"sruaXZqqSl","pcd":"123456","tk":"A09447AE-DEAA-461E-95AA-0BB256FAEF02","st":1,"geu":"game-api/piggy-gold/","lau":"game-api/lobby/","bau":"web-api/game-proxy/","cc":"THB","cs":"฿","nkn":"123456","gm":[{"gid":39,"msdt":1538637872000,"medt":1538637872000,"st":1,"amsg":""}],"uiogc":{"bb":1,"grtp":0,"gec":1,"cbu":0,"cl":0,"bf":0,"mr":0,"phtr":0,"vc":0,"bfbsi":0,"bfbli":0,"il":0,"rp":0,"gc":0,"ign":0,"tsn":0,"we":0,"gsc":0,"bu":0,"pwr":0,"hd":0,"et":0,"np":0,"igv":0,"as":0,"asc":0,"std":0,"hnp":0,"ts":0,"smpo":0,"grt":0,"ivs":1,"ir":0,"hn":1},"ec":[],"occ":{"rurl":"","tcm":"","tsc":0,"ttp":0,"tlb":"","trb":""},"gcv":"1.1.0.8","ioph":"bab4f5c16fb0"},"err":null}`

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, retStr)
}

func pg39GameInfo(w http.ResponseWriter, r *http.Request) {
	retStr := `{"dt":{"fb":null,"wt":{"mw":5.0,"bw":10.0,"mgw":20.0,"smgw":35.0},"maxwm":null,"cs":[1.0,5.0,30.0,100.0],"ml":[1,2,3,4,5,6,7,8,9,10],"mxl":1,"bl":0.00,"inwe":false,"iuwe":false,"ls":{"si":{"wp":null,"lw":null,"frl":[2,4,6,5,6,2,2,7,3],"pc":null,"wm":null,"tnbwm":null,"gwt":-1,"fb":null,"ctw":0.0,"pmt":null,"cwc":0,"fstc":null,"pcwc":0,"rwsp":null,"hashr":null,"ml":10,"cs":1.0,"rl":[4,6,7],"sid":"1831210061087641089","psid":"1831210061087641089","st":1,"nst":1,"pf":1,"aw":0.00,"wid":0,"wt":"C","wk":"0_C","wbn":null,"wfg":null,"blb":970.00,"blab":960.00,"bl":0.00,"tb":10.00,"tbb":10.00,"tw":0.00,"np":-10.00,"ocr":null,"mr":null,"ge":[1,11]}},"cc":"THB"},"err":null}`

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, retStr)
}
