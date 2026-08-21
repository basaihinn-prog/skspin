package staticproxy

import (
	_ "embed"
	"fmt"
	"game/duck/ut2/jwtutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	jsonStr := `{
    "status": "0000",
    "data": "%s"
}`
	r.ParseForm()
	//pid@THB
	pid, err := strconv.ParseInt(strings.Split(r.FormValue("uidA"), "@")[0], 10, 64)
	if err != nil {
		w.Write([]byte("{\"status\":\"8888\",\"data\":[]}"))
		return
	}
	token, err := jwtutil.NewToken(pid, time.Now().Add(time.Hour))
	if err != nil {
		w.Write([]byte("{\"status\":\"8888\",\"data\":[]}"))
		return
	}
	jsonStr = fmt.Sprintf(jsonStr, token)
	if err != nil {
		w.Write([]byte("{\"status\":\"8888\",\"data\":[]}"))
		return
	}
	w.Write([]byte(jsonStr))
}
