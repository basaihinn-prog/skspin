package comm

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"
)

// var ProxyUrl string
var (
	// 跳过证书验证
	//client := &http.Client{Transport: tr}
	cli = &http.Client{Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: runtime.NumCPU(),
		}}
)

func getHttpClient() *http.Client {
	return cli
	// httpTransport := &http.Transport{
	// 	DisableKeepAlives: true,
	// }

	// if ProxyUrl != "" {
	// 	dialer, err := proxy.SOCKS5("tcp", ProxyUrl, nil, proxy.Direct)
	// 	if err != nil {
	// 		logger.Err("创建socks5错误")
	// 		return nil
	// 	}
	// 	httpTransport.Dial = dialer.Dial
	// }

	// return &http.Client{Transport: httpTransport, Timeout: 10 * time.Second}
}

const (
	charsetUTF8                    = "charset=UTF-8"
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
)

func PostJson(ctx context.Context, url string, body any, dest any, fn func(req *http.Request)) error {
	respBody, err := PostJsonRaw(ctx, url, body, fn)
	if err != nil {
		return err
	}

	return json.Unmarshal(respBody, dest)
}

func PostJsonCode(ctx context.Context, url string, body any, dest any, fn func(req *http.Request)) error {
	respBody, err := PostJsonRaw(ctx, url, body, fn)
	if err != nil {
		return err
	}

	var codeErr struct {
		Error string `json:"error"`
		Data  any    `json:"data"`
	}
	fmt.Println(string(respBody))
	codeErr.Data = dest

	err = json.Unmarshal(respBody, &codeErr)
	if err != nil {
		return err
	}

	if codeErr.Error != "" {
		// return fmt.Errorf("%v", codeErr.Error)
		return errors.New(codeErr.Error)
	}

	return nil
}

func PostJsonRaw(ctx context.Context, url string, body any, fn1 func(req *http.Request)) ([]byte, error) {
	cli := getHttpClient()

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// id := ut2.RandomString(8)

	// logger.Info(id, "POST请求", url)
	// logger.Info(id, "BODY", string(bodyBytes))

	bodyReader := io.NopCloser(bytes.NewReader(bodyBytes))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		// logger.Info(id, "请求错误", err)
		return nil, err
	}
	req.Header.Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
	if fn1 != nil {
		fn1(req)
	}

	resp, err := cli.Do(req)
	if err != nil {
		// logger.Info(id, "请求错误", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// logger.Info(id, "io.ReadAll", err)
		return nil, err
	}

	// logger.Info(id, "请求结束", string(respBody))

	return respBody, nil
}

// func NewEcho() *echo.Echo {
// 	e := echo.New()

// 	e.HideBanner = true
// 	e.JSONSerializer = JSONSerializer{}
// 	e.HTTPErrorHandler = func(err error, c echo.Context) {
// 		if c.Response().Committed {
// 			return
// 		}

// 		he, ok := err.(*echo.HTTPError)
// 		if ok {
// 			if he.Internal != nil {
// 				if herr, ok := he.Internal.(*echo.HTTPError); ok {
// 					he = herr
// 				}
// 			}
// 			err = fmt.Errorf("%v", he.Message)
// 		}

// 		c.JSON(http.StatusOK, err)
// 	}

// 	return e
// }

// type JSONSerializer struct {
// 	echo.DefaultJSONSerializer
// }

// func (d JSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
// 	enc := json.NewEncoder(c.Response())
// 	if indent != "" {
// 		enc.SetIndent("", indent)
// 	}
// 	return enc.Encode(lazy.WrapMsg(i))
// }
