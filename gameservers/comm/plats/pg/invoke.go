package pg

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"serve/comm/define"
	"serve/comm/ut"

	"github.com/google/uuid"
)

type callError struct {
	Code    string
	Message string
}

type callResponse struct {
	Data  interface{}
	Error *callError
}

func invoke(method string, ps define.M, result interface{}) (err error) {

	cfg := GetConfig()
	// if err != nil {
	// 	return
	// }

	domain := cfg.PgSoftAPIDomain
	if method == "/Bet/v4/GetHistory" {
		domain = cfg.DataGrabAPIDomain
	}

	ul, err := url.Parse(domain)
	if err != nil {
		return
	}
	traceId := uuid.NewString()
	ul.RawQuery = "trace_id=" + traceId

	if strings.HasPrefix(method, "/external") {
		ul.Path = method
	} else {
		ul.Path = path.Join(ul.Path, method)
	}

	var values = define.M2Values(ps)
	values.Set("operator_token", cfg.OperatorToken)
	values.Set("secret_key", cfg.SecretKey)

	d := slog.With(
		"plat", "pg",
		"url", ul.String(),
		"payload", values.Encode(),
	)

	defer func() {
		d.With("err", err).Info("invoke!!")
	}()

	resp, err := http.PostForm(ul.String(), values)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	d = d.With("status", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	d = d.With("resp", ut.TruncMsg(body, 1024))

	if resp.Header.Get("Content-Type") == "text/html" {
		p := result.(*[]byte)
		*p = body
		return
	}

	var callRet callResponse
	callRet.Data = result
	err = json.Unmarshal(body, &callRet)
	if err != nil {
		return err
	}

	if callRet.Error != nil {
		code, _ := strconv.Atoi(callRet.Error.Code)
		err = define.NewErrCode(callRet.Error.Message, code)

		return
	}

	return
}
