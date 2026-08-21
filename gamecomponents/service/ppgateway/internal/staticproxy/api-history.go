package staticproxy

import (
	"fmt"
	"game/comm/mq"
	"net/http"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// /api/history/v2/settings/general
// func

// https://5g6kpi7kjf.uapuqhki.net/gs2c/api/history/v2/play-session/last-items?token=97c1960d32e717aa911478261beccc19a6dd8c9d1a11a326b5bbf64fdf8f0bf6&symbol=vs20olympx

// https://5g6kpi7kjf.uapuqhki.net/gs2c/api/history/v3/action/children?id=39295534355091&token=e9f224bec8a561da3ec0f1e0444c4e8d3c0e5b3ed2c69d4465a46c1e68a0d9b6&symbol=vs20olympx

// https://backoffice-sg53.ppgames.net/admin/api/history/v3/action/children?id=39297197956091&token=abdfef6dde5532244529e4c9577a185af22294eb75d820ea31fe5711758235c7

// https://ppgames.rpgamestest.com/gs2c/api/history/v3/action/children?id=39297197956091&token=

// https://ppgames.rpgamestest.com/gs2c/api/history/v2/play-session/by-round?id=66ebbb292235a2b8b233d713&token=&symbol=vs20olympx

func history_forward(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	subj := fmt.Sprintf("pp_%s.%s", query.Get("symbol"), r.URL.Path)

	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: subj,
		Data:    []byte(r.URL.RawQuery),
	}, time.Second*60)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errstr := resp.Header.Get("error"); errstr != "" {
		http.Error(w, errstr, http.StatusInternalServerError)
		return
	}
	// 设置 CORS 相关响应头
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 设置缓存控制
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	// 设置内容编码
	//w.Header().Set("Content-Encoding", "gzip")
	// 设置内容类型
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	w.Write(resp.Data)
}

func history_by_round(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	pair := strings.Split(query.Get("id"), "@")
	symbol, id := pair[0], pair[1]

	// query.Add("symbol", symbol)
	// query.Add("bet", round)

	subj := fmt.Sprintf("%s.%s", symbol, r.URL.Path)

	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: subj,
		Data:    []byte(id),
	}, time.Second*60)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json;charset=UTF-8")

	if errstr := resp.Header.Get("error"); errstr != "" {
		http.Error(w, errstr, http.StatusInternalServerError)
		return
	}

	w.Write(resp.Data)
}
