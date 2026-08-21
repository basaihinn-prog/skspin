package jili

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"serve/comm/define"
	"serve/comm/ut"

	"github.com/samber/lo"
	"golang.org/x/net/proxy"
)

var (
	loc = time.FixedZone("UTC-4", -4*60*60)
)

func keynow() string {
	now := time.Now().In(loc)
	return now.Format("06012")

	// return "19011"
}

func toMD5Str(s string) string {
	data := md5.Sum([]byte(s))
	return hex.EncodeToString(data[:])
}

func keyG(cfg JiliConfig) string {
	ans := toMD5Str(keynow() + cfg.AgentId + cfg.AgentKey)
	return ans
	// return "886cc5601e01c7b342b7ed7ec87f202c"
}

type CallResponse struct {
	ErrorCode int
	Message   string
	Data      interface{}
}

var _httpclient *http.Client

func getHttpClient() *http.Client {
	if _httpclient == nil {
		// ssh -Nf  -D 127.0.0.1:1080 doudou-test
		dialer := lo.Must(proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct))

		_httpclient = &http.Client{
			Transport: &http.Transport{
				Dial: dialer.Dial,
			},
		}
	}

	return _httpclient
}

func invoke(method string, args OrderedArgs, result interface{}) (err error) {
	cfg := jiliConfig

	args = args.Append("AgentId", cfg.AgentId)

	// Key = {6 个任意字符} + MD5(所有请求参数串 + KeyG) + {6 个任意字符}
	key := "000000" + toMD5Str(args.Encode()+keyG(cfg)) + "000000"
	args = args.Append("Key", key)

	u := lo.Must(url.JoinPath(cfg.ApiUrl, method))
	// u := cfg.ApiUrl + method

	// return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	querystring := args.Encode()

	d := slog.With(
		"plat", "jili",
		"url", u,
		"payload", querystring,
	)

	defer func() {
		d.With("err", err).Info("invoke!!")
	}()

	resp, err := /*getHttpClient()*/ http.Post(u, "application/x-www-form-urlencoded", strings.NewReader(querystring))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ret, err := io.ReadAll(resp.Body)
	// fmt.Printf("jili body: [%s]\n", ret)
	if err != nil {
		return
	}

	d = d.With("resp", ut.TruncMsg(ret, 1024))
	// log.Printf("Jili>>>请求[%s]: %s\nJili>>>返回:%s\n", method, querystring, comm.TruncMsg(ret, 512))

	var callRet CallResponse
	callRet.Data = result
	err = json.Unmarshal(ret, &callRet)
	if err != nil {
		return err
	}

	if callRet.ErrorCode != 0 {
		// return logic.CodeError{
		// 	Code:    callRet.ErrorCode,
		// 	Message: callRet.Message,
		// }

		return define.NewErrCode(callRet.Message, callRet.ErrorCode)
	}

	return
}
