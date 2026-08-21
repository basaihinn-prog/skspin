package staticproxy

import (
	"fmt"
	"game/comm/mq"
	"github.com/nats-io/nats.go"
	"io"
	"net/http"
	"time"
)

func RTP(w http.ResponseWriter, r *http.Request) {
	payload, _ := io.ReadAll(r.Body)
	header := nats.Header(r.Header)
	header.Set("query", r.URL.RawQuery)

	url := r.URL.Path + ".rtp"
	fmt.Println(url)
	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: url,
		Data:    payload,
		Header:  header,
	}, time.Second*60)

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

	w.Write(resp.Data)
}
