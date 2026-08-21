package staticproxy

import (
	"context"
	"encoding/json"
	"game/comm"
	"game/comm/db"
	"game/comm/mq"
	"game/comm/ut"
	"game/duck/ut2/jwtutil"
	"game/service/jiligateway/internal/gamedata"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

func init() {
	uat_history_api_mux.HandleFunc("/", static_assets)
	uat_history_api_mux.HandleFunc("/token", history_token)
	uat_history_api_mux.HandleFunc("/game-entrance/", game_entrance)
	// uat_history_api_mux.HandleFunc("/game-entrance/", game_entrance)

	// https://uat-history-api.kafa010.com/history/csh/get-history-record?EndRowIndex=10&LangId=en-US&LogIndexAsRoundIndex=false&Minutes=60&StartRowIndex=1
	uat_history_api_mux.HandleFunc("GET /history/{game}/get-history-record", wrapCallHistory("/history/{game}/get-history-record"))

	// https://uat-history-api.jlfafafa3.com/history/csh/get-single-round-log-summary/en-US/1718692563274366002
	uat_history_api_mux.HandleFunc("GET /history/{game}/get-single-round-log-summary/", wrapCallHistory("/history/{game}/get-single-round-log-summary/"))

	// https://uat-history-api.jlfafafa3.com/history/csh/get-log-plate-info/1718692563274366002/1718692563274376002
	uat_history_api_mux.HandleFunc("GET /history/{game}/get-log-plate-info/", wrapCallHistory("/history/{game}/get-log-plate-info/"))

	// https://uat-history-api.jlfafafa3.com/game-setting/is-enable/NickNameEnable
	uat_history_api_mux.HandleFunc("GET /game-setting/is-enable/NickNameEnable", nickNameEnable)

}

func nickNameEnable(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"Code":0,"Message":"成功","Data":{"IsEnable":false}}`))
}

// https://uat-history-api.jlfafafa3.com/token
func history_token(w http.ResponseWriter, r *http.Request) {
	// {"Token":"eyJQIjoxMDA4NTksIkUiOjE3MTg0NjAyMTQsIlMiOjEwMDEsIkQiOiJqaWxpXzJfY3NoIn0.49gDSLz8y3E_Q9dAhoyWmKSXbsFqS76K6frxZa9wQBY"}

	var token struct {
		Token string
	}

	err := ut.HttpRequestJson(r, &token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pid, err := jwtutil.ParseToken(token.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// {
	// "Code": 0,
	// "Message": "成功",
	// "Data": {
	//     "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjkxNTkxIiwianRpIjoiMGU5MGQxYWVhODU5NDdjMzk3YjJjNDI4MDFmOTQ1YWIiLCJBSUQiOiIxMjkxNTkxIiwiQVQiOiJ0ZXN0dXNlcjEyIiwiQVBJSUQiOiIyMTcxIiwiU0lURUlEIjoiMCIsIkxJTkVDT0RFIjoiMCIsIkxJQ0VOU0UiOiIwIiwiQVBJVFlQRSI6IjAiLCJXQUxMRVRUWVBFIjoiMSIsIlRIRU1FSUQiOiIwIiwiQ05PIjoiMyIsIkVYUE9ORU5UIjoiNCIsIlJBVElPIjoiMSIsIkNOQU1FIjoiIiwiQ1NZTUJPTCI6IiIsIlVOSVQiOiIxIiwiUkFURSI6IjEiLCJuYmYiOjE3MTg0MTU5NjUsImV4cCI6MTcxODQxOTU2NSwiaWF0IjoxNzE4NDE1OTY1LCJpc3MiOiJEcmFnb25GaXNoIiwiYXVkIjoiR2FtZUhpc3RvcnkifQ.WzhXNoSDE2SCxNAbcZa5tOcW8In9ASGTZp7MeQDK4-0"
	// }
	// '{"sub":"1234914","jti":"1a78dfc67e8a4d82b0c2761d6d41f25e","AID":"1234914","AT":"testuser1","APIID":"2171","SITEID":"0","LINECODE":"0","LICENSE":"0","APITYPE":"0","WALLETTYPE":"1","THEMEID":"0","CNO":"3","EXPONENT":"4","RATIO":"1","CNAME":"","CSYMBOL":"","UNIT":"1","RATE":"1","nbf":1718847811,"exp":1718851411,"iat":1718847811,"iss":"DragonFish","aud":"GameHistory"}'

	now := time.Now()
	nbf := now.Unix()

	aid := strconv.Itoa(int(pid))

	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": aid, "jti": "1a78dfc67e8a4d82b0c2761d6d41f25e", "AID": aid, "AT": aid, "APIID": "2171", "SITEID": "0", "LINECODE": "0", "LICENSE": "0", "APITYPE": "0", "WALLETTYPE": "1", "THEMEID": "0", "CNO": "3", "EXPONENT": "4", "RATIO": "1", "CNAME": "", "CSYMBOL": "", "UNIT": "1", "RATE": "1", "nbf": nbf, "exp": nbf + 3600, "iat": nbf, "iss": "DragonFish", "aud": "GameHistory",
	}).SignedString(jwtutil.SignatureKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ut.HttpReturnJson(w, map[string]any{
		"Code":    0,
		"Message": "成功",
		"Data": map[string]any{
			// "Token": token.Token,
			"Token": s,
		}})
}

// {"Code":0,"Message":"成功", "Data": ""}
type GameInfo struct {
	No       int    `json:"No"`
	Id       string `json:"Id"`
	Show     bool   `json:"Show"`
	IconType int    `json:"IconType"`
	Multiple int    `json:"Multiple"`
	Jackpot  bool   `json:"Jackpot"`
}

// var entranceStr []byte
// var entranceStrExpir int64
// var entranceMutex sync.Mutex

func game_entrance(w http.ResponseWriter, r *http.Request) {
	// entranceMutex.Lock()
	// var data []byte
	// defer func() {
	// 	io.WriteString(w, fmt.Sprintf("{\"Code\":0,\"Message\":\"成功\", \"Data\": \"%v\"}", string(data)))
	// 	entranceMutex.Unlock()
	// }()
	// if entranceStrExpir > time.Now().Unix() {
	// 	data = entranceStr
	// 	return
	// }
	coll := db.Collection2("game", "Games")
	var gameList []*comm.Game
	cur, err := coll.Find(context.TODO(), bson.M{"_id": bson.M{"$regex": "^jili_"}})
	if err != nil {
		return
	}
	err = cur.All(context.TODO(), &gameList)
	if err != nil {
		return
	}

	gameMap := gamedata.GameMap()
	ret := make([]GameInfo, 0, len(gameList))
	for j := range gameList {
		args := strings.Split(gameList[j].ID, "_")
		if len(args) != 3 {
			continue
		}
		no, _ := strconv.Atoi(args[1])

		game := gameMap[no]
		if game == nil {
			continue
		}

		ret = append(ret, GameInfo{
			No:       no,
			Id:       game.Id,
			Show:     true,
			IconType: 3,
			Multiple: 0,
			Jackpot:  false,
		})

	}

	// data, _ = json.Marshal(ret)

	// ans, _ := json.Marshal(map[string]any{
	// 	"Code":    0,
	// 	"Message": "成功",
	// 	"Data":    ret,
	// })

	w.Header().Set("Cache-Control", "max-age=3600")
	w.Write(MarshalJsonReturn(ret))

	// entranceStrExpir = time.Now().Add(10 * time.Second).Unix()
	// entranceStr = data
}

// https://uat-history-api.jlfafafa3.com/history/csh/get-log-plate-info/1718692563274366002/1718692563274376002

func wrapCallHistory(subject string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		game := r.PathValue("game")
		subject := strings.Replace(subject, "{game}", game, 1)
		payload := r.URL.String()
		// payload := r.URL.RawQuery
		ps := r.URL.Query()
		token := ps.Get("Token")
		_, game, err := jwtutil.ParseTokenData(token)
		if err != nil {
			return
		}
		plat := `jili:`
		if strings.HasPrefix(game, "tada") {
			plat = `tada:`
		}
		resp, err := mq.NC().RequestMsg(&nats.Msg{
			Subject: plat + subject,
			Data:    []byte(payload),
			Header:  nats.Header(r.Header),
		}, time.Second*60)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if errstr := resp.Header.Get("error"); errstr != "" {
			http.Error(w, errstr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		w.Write(resp.Data)
	}
}

type ReturnObj struct {
	Code    int
	Message string
	Data    any
}

func MarshalJsonReturn(data any) []byte {
	return lo.Must(json.Marshal(ReturnObj{
		Code:    0,
		Message: "成功",
		Data:    data,
	}))
}
