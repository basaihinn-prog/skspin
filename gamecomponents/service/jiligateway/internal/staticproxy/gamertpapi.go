package staticproxy

import (
	"errors"
	"fmt"
	"game/comm/mq"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func gamertpapi(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	//if !(r.Method == http.MethodPost && (contentType == "application/x-protobuf" || contentType == "application/json") || strings.HasPrefix(r.URL.Path, "/mp/")) {
	//if !(r.Method == http.MethodPost && (contentType == "application/x-protobuf" || contentType == "application/json") || strings.HasPrefix(r.URL.Path, "/game-rtp-api")) {
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//}

	var (
		payload []byte
		err     error
	)

	if r.Method == http.MethodOptions {
		return
	}

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

	fmt.Println(string(payload))

	header := nats.Header(r.Header)
	header.Set("query", r.URL.RawQuery)
	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: r.URL.Path,
		Data:    payload,
		Header:  header,
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
