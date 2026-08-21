package ppcomm

import (
	"cmp"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"serve/comm/plats/pp"

	"serve/comm/mux"
)

func FetchMgckey(game, user string, client *http.Client) (mgckey string, location *url.URL, err error) {
	var urlret struct {
		Url string
	}
	err = mux.HttpInvoke("https://plats.rpgamestest.com/plat/PP/LaunchGame", map[string]any{
		"Game": game,
		"UID":  user,
		"Lang": "en",
	}, &urlret)

	if err != nil {
		return
	}

	// client req the rest
	req, err := http.NewRequest(http.MethodGet, urlret.Url, nil)

	if err != nil {
		return
	}

	fmt.Println(urlret.Url)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	// resp, err := http.DefaultTransport.RoundTrip(req)
	trans := cmp.Or(client.Transport, http.DefaultTransport)
	resp, err := trans.RoundTrip(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		err = errors.New(resp.Status)
		return
	}

	location, err = url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return
	}

	query := location.Query()
	mgckey = query.Get("mgckey")
	location.RawQuery = ""

	return
}

func GetMgckey(game, user string, client *http.Client) (mgckey string, location *url.URL, err error) {
	// client req the rest
	uid, guid, cash := "ceshi1938", "4ae52ed2e1a8c353878ba65ed7791ac4", "5000000"
	uid = user
	guid = GuidMap[game]
	uurl, _ := pp.GetKey(uid, guid, cash)
	//req, err := http.NewRequest(http.MethodGet, uurl, nil)
	//
	//if err != nil {
	//	return
	//}
	//
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	//
	//trans := cmp.Or(client.Transport, http.DefaultTransport)
	//resp, err := trans.RoundTrip(req)
	//if err != nil {
	//	return
	//}
	//
	//defer resp.Body.Close()
	//if resp.StatusCode != http.StatusFound {
	//	err = errors.New(resp.Status)
	//	return
	//}
	//
	//location, err = url.Parse(resp.Header.Get("Location"))
	//if err != nil {
	//	return
	//}
	location, err = GetRedirect(uurl, client)
	if err != nil {
		err = errors.New("重定向失败")
		return
	}
	if location.String() == "" {
		err = errors.New("重定向失败")
		return
	}
	location, err = GetRedirect(location.String(), client)
	if err != nil {
		err = errors.New("重定向失败")
		return
	}
	query := location.Query()
	mgckey = query.Get("mgckey")
	location.RawQuery = ""

	return
}

// 重定向方法
func GetRedirect(uurl string, client *http.Client) (*url.URL, error) {
	var location *url.URL
	req, err := http.NewRequest(http.MethodGet, uurl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	trans := cmp.Or(client.Transport, http.DefaultTransport)
	resp, err := trans.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		err = errors.New(resp.Status)
		return nil, err
	}

	location, err = url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return nil, err
	}
	return location, nil
}
