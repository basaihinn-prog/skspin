package pg

import (
	"encoding/json"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"strings"

	"serve/comm/define"
	"serve/comm/ut"

	"github.com/samber/lo"
)

func genTraceId() string {
	// return "VWDLLC27"
	var buf [6]byte
	for i := 0; i < len(buf); i++ {
		c := rand.IntN(26) + 'A'
		buf[i] = byte(c)
	}

	return string(buf[:]) + "27"
}

var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

func InvokePGService(m string, ps url.Values, pret any) (err error) {
	var apiurl string
	if strings.HasPrefix(m, "http") {
		apiurl = m
	} else {
		apiurl = lo.Must(url.JoinPath(pgcfg.ClientApiURL, m))
	}

	apiurl += "?traceId=" + genTraceId()

	lg := logger.With("apiurl", apiurl, "ps", ps.Encode())
	defer func() {

		lg.Info("InvokePGService", "error", err)
	}()
	resp, err := http.PostForm(apiurl, ps)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	lg = lg.With("body", ut.TruncMsg(body, 1024))

	var pgret define.PGRetWrapper
	err = json.Unmarshal(body, &pgret)
	if err != nil {
		return
	}
	if pgret.Err != nil && pgret.Err.Msg != "" {
		err = pgret.Err
		return
	}

	if pret != nil && pgret.Dt != nil {
		err = json.Unmarshal(*pgret.Dt, pret)
	}
	return
}

type GetPlayerTokenRet struct {
	Geu string
	Lau string
	Bau string
	Tk  string

	Sid  string
	PSid string
}

// type getPlayerToken
func GetPlayerToken(uid, gameid string) (ret GetPlayerTokenRet, err error) {
	pg := PG()
	err = pg.Regist(uid)
	if err != nil {
		return
	}

	balance, err := pg.GetBalance(uid)
	if err != nil {
		return
	}
	if balance < 10000000 {
		pg.FundTransferIn(uid, 10000000)
	}

	os, err := GenToken(uid)
	if err != nil {
		return
	}

	var ps = url.Values{}
	ps.Add("btt", "1")
	ps.Add("vc", "2")
	ps.Add("pf", "1")
	ps.Add("l", "en")
	ps.Add("gi", gameid)
	ps.Add("os", os)
	ps.Add("otk", pgcfg.OperatorToken)

	// ul := lo.Must(url.JoinPath(pgcfg.ClientApiURL, "/web-api/auth/session/v2/verifyOperatorPlayerSession"))

	err = InvokePGService("/web-api/auth/session/v2/verifyOperatorPlayerSession", ps, &ret)
	if err != nil {
		return
	}

	// https://api.pg-demo.com/game-api/piggy-gold/v2/GameInfo/Get?traceId=ICQUNP27
	// btt=1&atk=710BC381-46A8-470D-AD3C-D1C93A9D8269&pf=1

	var gameinfo struct {
		Ls struct {
			Si struct {
				Sid  string
				PSid string
			}
		}
	}
	ul, _ := url.JoinPath(ret.Geu, "/v2/GameInfo/Get")
	err = InvokePGService(ul, url.Values{
		"btt": []string{"1"},
		"pf":  []string{"1"},
		"atk": []string{ret.Tk},
	}, &gameinfo)
	if err != nil {
		return
	}

	ret.PSid = gameinfo.Ls.Si.PSid
	ret.Sid = gameinfo.Ls.Si.Sid

	return
}
