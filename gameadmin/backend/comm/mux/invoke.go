package mux

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

var (
	cli = &http.Client{Timeout: 10 * time.Second}
)

func HttpInvoke(url string, args, reply any) (err error) {
	lg := slog.With("url", url)
	defer func() {
		lg.Info("HttpInvoke", "error", err)
	}()

	replyType := reflect.TypeOf(reply)
	if reply != nil && replyType.Kind() != reflect.Ptr {
		log.Panic("reply type should be a pointer")
	}

	payload, err := json.Marshal(args)
	if err != nil {
		return
	}
	lg = lg.With("payload", string(payload))

	resp, err := cli.Post(url, ContentTypeJson, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var ret Response
	ret.Data = reply

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	lg = lg.With("respbody", string(body))

	err = json.Unmarshal(body, &ret)
	// err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return err
	}

	if ret.Error != "" {
		return errors.New(ret.Error)
	}
	return nil
}
