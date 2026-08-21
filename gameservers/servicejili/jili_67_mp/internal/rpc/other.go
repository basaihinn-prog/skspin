package rpc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"serve/servicejili/jili_67_mp/internal"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/unionjp/jpinfo", internal.GameShortName), jpinfo)
	jiliut.RegRpc(fmt.Sprintf("/%s/unionjp/jpinfoall", internal.GameShortName), jpinfoall)
	jiliut.RegRpc(fmt.Sprintf("/%s/fulljp/jpinfo", internal.GameShortName), fjpinfo)
	jiliut.RegRpc(fmt.Sprintf("/%s/fulljp/jpinfoall", internal.GameShortName), fjpinfoall)
	jiliut.RegRpc(fmt.Sprintf("/%s/item/itemlist", internal.GameShortName), itemlist)
	jiliut.RegRpc(fmt.Sprintf("/%s/item/getmail", internal.GameShortName), getmail)
	jiliut.RegRpc(fmt.Sprintf("/%s/item/allgameitemlist", internal.GameShortName), allgameitemlist)
	jiliut.RegRpc(fmt.Sprintf("/%s/vipsignup/signinfo", internal.GameShortName), signinfo)
	jiliut.RegRpc(fmt.Sprintf("/%s/mission/getdailymission", internal.GameShortName), getdailymission)
	jiliut.RegRpc(fmt.Sprintf("/%s/account/heart", internal.GameShortName), heart)

}

func jpinfo(ps *nats.Msg) (ret []byte, err error) {
	retData := map[string]interface{}{
		"info": nil,
		"ret":  235,
		"type": 61,
	}
	ret, _ = json.Marshal(retData)
	return
}

func jpinfoall(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")
	retData := map[string]interface{}{
		"info": map[string]interface{}{
			"glist": nil,
			"plist": []map[string]interface{}{
				{
					"type":  1,
					"value": 0,
				},
				{
					"type":  2,
					"value": 0,
				},
				{
					"type":  3,
					"value": 0,
				},
			},
		},
		"ret":   0,
		"type":  64,
		"token": token,
	}
	ret, _ = json.Marshal(retData)
	return
}

func fjpinfo(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")
	retData := map[string]interface{}{
		"info": map[string]interface{}{
			"full":   0,
			"minbet": 0,
			"minvip": 0,
			"value":  0,
		},
		"token": token,
		"ret":   0,
		"type":  66,
	}
	ret, _ = json.Marshal(retData)
	return
}

func fjpinfoall(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")
	retData := map[string]interface{}{
		"info": map[string]interface{}{
			"acktype": 0,
			"list":    []interface{}{},
		},
		"token": token,
		"ret":   0,
		"type":  68,
	}
	ret, _ = json.Marshal(retData)
	return
}

func itemlist(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")

	retData := map[string]interface{}{
		"info": map[string]interface{}{
			"playeritemdata": nil,
			"result":         0,
		},
		"token": token,
		"ret":   0,
		"type":  21,
	}
	ret, _ = json.Marshal(retData)
	return
}

func getmail(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")

	retData := map[string]interface{}{
		"info": map[string]interface{}{
			"isread":          1,
			"result":          0,
			"usermailackdata": nil,
		},
		"token": token,
		"ret":   0,
		"type":  24,
	}
	ret, _ = json.Marshal(retData)
	return
}

func allgameitemlist(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")

	retData := map[string]interface{}{
		"info": map[string]interface{}{
			"result":     0,
			"canuseitem": nil,
		},
		"token": token,
		"ret":   0,
		"type":  26,
	}
	ret, _ = json.Marshal(retData)
	return
}

func signinfo(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")

	retData := map[string]interface{}{
		"ack": map[string]interface{}{
			"Info": map[string]interface{}{
				"Error":        1,
				"ExpiredTime":  "1970-01-01 00:00:00.000",
				"TreasureList": nil,
			},
			"Result": 0,
		},
		"token": token,
		"ret":   0,
		"type":  43,
	}
	ret, _ = json.Marshal(retData)
	return
}

func getdailymission(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")

	retData := map[string]interface{}{
		"info":  map[string]interface{}{},
		"token": token,
		"ret":   0,
		"type":  31,
	}
	ret, _ = json.Marshal(retData)
	return
}

func heart(ps *nats.Msg) (ret []byte, err error) {
	retData := map[string]interface{}{
		"info": map[string]interface{}{
			"message": "test01",
			"state":   0,
		},
		"ret":  0,
		"type": 98,
	}
	ret, _ = json.Marshal(retData)
	return
}
