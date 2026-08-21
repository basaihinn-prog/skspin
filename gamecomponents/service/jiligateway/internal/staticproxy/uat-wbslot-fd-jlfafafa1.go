package staticproxy

import (
	"game/comm/define"
	"net/http"
)

func init() {
	uat_wbslot_fd_mux.HandleFunc("/webservice/event/trigger", event_trigger)

	// https://uat-wbslot-fd-jlfafafa1.kafa010.com/webservice/event/error
	uat_wbslot_fd_mux.HandleFunc("/webservice/event/error", event_error)

	uat_wbslot_fd_mux.HandleFunc("/subagentservice/MakeUserSubAgent", makeUserSubAgent)

	uat_wbslot_fd_mux.HandleFunc("/rankingservice/user/GetMailList", getMailList)
	uat_wbslot_fd_mux.HandleFunc("/favoriteservice/OnLogin", favoriteserviceOnLogin)
	uat_wbslot_fd_mux.HandleFunc("/rankingservice/user/GetRankingListV2", getRankingList)
	uat_wbslot_fd_mux.HandleFunc("/vipservice/VIPGet", vipGet)
	uat_wbslot_fd_mux.HandleFunc("/smartnotice/notice/getReq", noticeGet)
	uat_wbslot_fd_mux.HandleFunc("/promotionservice/OnLoginV3", noticeGet)

	uat_wbslot_fd_mux.HandleFunc("/webservice/event/user", eventuser)
	uat_wbslot_fd_mux.HandleFunc("/me/fulljp/JPInfoProto", event_trigger)
	uat_wbslot_fd_mux.HandleFunc("/me/fulljp/JPInfoAllProto", event_trigger)

	uat_wbslot_fd_mux.HandleFunc("/", gameapi)
	//uat_wbslot_fd_mux.HandleFunc("/game-rtp-api/jili_21_ols/game/spinstat", gamertpapi)
	uat_wbslot_fd_mux.Handle("/game-rtp-api/", http.StripPrefix("/game-rtp-api/", http.HandlerFunc(gamertpapi)))

}

func eventuser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=604800")
	w.Write([]byte("{}"))
}

// https://uat-wbslot-fd-jlfafafa1.kafa010.com/webservice/event/user?

func event_trigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=604800")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.Write([]byte("{}"))
}

func event_error(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=604800")

	w.Write([]byte("{}"))
}

type (
	M = define.M
	D = define.D
)

// https://uat-wbslot-fd.jlfafafa1.com/subagentservice/MakeUserSubAgent?accountId=944261&apiId=555&gameId=2&siteId=114298358&agentId=rp_Online@api-555.game&linecode=0&site=jilid.rslotszs001.com&token=882f3dc924647733b90cb50eba762135a413b7ba
func makeUserSubAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=604800")
	w.Write([]byte(`"OK"`))
}

// https://uat-wbslot-fd.jlfafafa1.com/rankingservice/user/GetMailList?AccountId=944261&Lang=en-US&token=13bca6568a0348fda246362742870c192d62a0f0
func getMailList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Tye", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=604800")
	w.Write([]byte(`{"Error":"","Mails":[]}`))
}

func favoriteserviceOnLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Tye", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=604800")
	w.Write([]byte(`{"cmdType":4,"content":{"Enabled":true,"Favorites":null,"Promotions":null,"Expired":null,"DAU":null,"BigWined":false}}`))
}

func getRankingList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Tye", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=604800")
	w.Write([]byte(`{"Data":null,"Error":"ranking not exist or begin","NeedWebView":false,"ExtraWebView":0}`))
}

func vipGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Tye", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=604800")
	w.Write([]byte(`{"Data":null,"Error":"setting not exist"}`))
}

func noticeGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Tye", "application/x-protobuf")
	w.Header().Set("Cache-Control", "public, max-age=604800")
	// w.Write([]byte(`{"Data":null,"Error":"setting not exist"}`))
}
