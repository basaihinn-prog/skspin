package staticproxy

import (
	"errors"
	"fmt"
	"game/comm/define"
	"game/comm/mq"
	"game/duck/ut2/jwtutil"
	"game/service/jiligateway/internal/gamedata"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func checkIp(ip string) (allow bool, err error) {
	var netAllow struct {
		Allow bool
	}
	err = mq.Invoke("/AdminInfo/Interior/whiteList_exists", map[string]any{
		"IP": ip,
	}, &netAllow)
	if err != nil {
		return
	}
	allow = netAllow.Allow
	return
}

func gameapi(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !(r.Method == http.MethodPost && (contentType == "application/x-protobuf" || contentType == "application/json") || strings.HasPrefix(r.URL.Path, "/mp/")) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var (
		payload []byte
		err     error
	)

	switch r.Method {
	case http.MethodPost:
		payload, err = io.ReadAll(r.Body)
	case http.MethodGet:
		// payload = []byte(r.URL.RawQuery)
	default:
		err = errors.New("wrong method! " + r.Method)
	}

	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte(err.Error()))

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//ip校验
	_, ip, loc := gamedata.IsBlockLoc(r)
	fmt.Printf("userid=%d,ip=%s,loc=%s\n", "", ip, loc)
	netAllow, err := checkIp(ip)
	if err != nil {
		err = define.NewErrCode("check ip error", 1302)
		return
	}
	if !netAllow {
		err = define.NewErrCode("Your ip exinclude in net whitelist", 1302)
		return
	}

	header := nats.Header(r.Header)
	header.Set("query", r.URL.RawQuery)
	plat := `jili:`
	token := header.Get("Token")
	if token != "" {
		_, game, err := jwtutil.ParseTokenData(token)
		if err != nil {
			return
		}
		if strings.HasPrefix(game, "tada") {
			plat = `tada:`
		}
	}

	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: plat + r.URL.Path,
		Data:    payload,
		Header:  nats.Header(header),
	}, time.Second*60)

	// resp, err := mq.NC().Request("jili:"+r.URL.Path, payload, time.Second*60)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if errstr := resp.Header.Get("error"); errstr != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errstr))
		return
	}

	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	if strings.HasPrefix(r.URL.Path, "/mp/") {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	w.Write(resp.Data)
}
