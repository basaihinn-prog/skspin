package pp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"serve/comm/define"
	"serve/comm/mux"
	"serve/comm/ut"
)

type callResponse struct {
	Error       string
	Description string
}

func invoke(method string, args map[string]string, result interface{}) (err error) {
	cfg := ppConfig

	var data = url.Values{}
	data.Set("secureLogin", cfg.SecureLogin)
	//THB

	// args := platcomm.GetArgs(ps)
	for k, v := range args {
		data.Set(k, v)
	}

	content := data.Encode() + cfg.SecretKey

	data.Set("hash", ut.ToMD5Str(content))

	fields := slog.With(
		"plat", "pp",
		"url", cfg.ApiUrl+method,
		"reqpayload", data.Encode(),
	)
	defer func() {
		// fields["error"] = err
		// logrus.WithFields(fields).Info("invoke!!")

		fields.Info("invoke!!", "error", err)
	}()

	resp, err := http.PostForm(cfg.ApiUrl+method, data)
	if err != nil {
		return
	}
	fields = fields.With("httpstatus", resp.Status)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fields = fields.With("respbody", mux.TruncMsg(ret, 1024))

	var callRet callResponse
	err = json.Unmarshal(ret, &callRet)
	if err != nil {
		return err
	}
	if callRet.Error != "0" {
		errcode, _ := strconv.Atoi(callRet.Error)
		// return logic.CodeError{
		// 	Code:    errcode,
		// 	Message: callRet.Description,
		// }
		return define.NewErrCode(callRet.Description, errcode)
	}

	if result != nil {
		err = json.Unmarshal(ret, result)
	}
	return
}

func invokeFRB(method string, query map[string]string, payload any, result interface{}) (err error) {
	cfg := ppConfig

	var data = url.Values{}
	data.Set("secureLogin", cfg.SecureLogin)
	//THB

	// args := platcomm.GetArgs(ps)
	for k, v := range query {
		data.Set(k, v)
	}

	content := data.Encode() + cfg.SecretKey

	data.Set("hash", ut.ToMD5Str(content))

	fields := slog.With(
		"plat", "pp",
		// "url", cfg.ApiUrl+method,
		// "reqpayload", data.Encode(),
	)
	defer func() {
		fields.Info("invoke!!", "error", err)
	}()

	ul, err := url.Parse(cfg.ApiUrl)
	if err != nil {
		return
	}

	ul = ul.JoinPath("../FreeRoundsBonusAPI", method)

	ul.RawQuery = data.Encode()

	fields = fields.With("url", ul.String())

	pload, err := json.Marshal(payload)
	if err != nil {
		return
	}
	fields = fields.With("payload", string(pload))

	req, _ := http.NewRequest(http.MethodPost, ul.String(), bytes.NewReader(pload))
	req.Header.Set("Content-Type", "application/json")

	respbody, _, err := ut.DoHttpReq(http.DefaultClient, req)

	if err != nil {
		return
	}
	fields = fields.With("respbody", mux.TruncMsg(respbody, 1024))

	var callRet callResponse
	err = json.Unmarshal(respbody, &callRet)
	if err != nil {
		return err
	}
	if callRet.Error != "0" {
		errcode, _ := strconv.Atoi(callRet.Error)
		return define.NewErrCode(callRet.Description, errcode)
	}

	if result != nil {
		err = json.Unmarshal(respbody, result)
	}
	return
}
