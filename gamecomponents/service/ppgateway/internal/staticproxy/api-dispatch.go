package staticproxy

import (
	"bytes"
	"fmt"
	"game/comm/mq"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func robotsCenter(w http.ResponseWriter, r *http.Request) {
	// action=doInit&symbol=vs20olympx&cver=230048&index=1&counter=1&repeat=0&mgckey=eyJQIjoxMDA3MDUsIkUiOjE3MjI0MzM4OTgsIlMiOjEwMDMsIkQiOiJwcF92czIwb2x5bXB4In0.GQtZWIrU-lZXjWNHiQhodLFd_EsAPUr31QyO9LoOgI4
	// action=doInit&symbol=vs20olympx&cver=237859&index=1&counter=1&repeat=0&mgckey=AUTHTOKEN@17f8076abf3edb3c651be065210634dcd0b4ba1231c206e4405c0fa4455b9ec0~stylename@hllgd_hollygod~SESSION@5df8d3d1-dfa3-4cff-983f-6bd0853c162f~SN@273d3dc7

	var (
		payload []byte
		err     error
	)

	payload, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(payload))
	r.Body = io.NopCloser(bytes.NewReader(payload))

	header := nats.Header(r.Header)
	header.Set("query", r.URL.RawQuery)

	subj := strings.Replace(r.URL.Path[1:], "/", ".", -1)

	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: subj,
		Data:    payload,
		Header:  header,
	}, time.Second*60)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errstr := resp.Header.Get("error"); errstr != "" {
		http.Error(w, errstr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源，或指定特定源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.Write(resp.Data)
}
