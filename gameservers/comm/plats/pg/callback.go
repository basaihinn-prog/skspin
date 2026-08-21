package pg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/pg/verifySession", verifySession)
}

func verifySession(w http.ResponseWriter, r *http.Request) {
	var resp callResponse
	if err := r.ParseForm(); err != nil {
		resp.Error = &callError{
			Code:    "1200",
			Message: err.Error(),
		}
	} else {
		username := r.PostForm.Get("operator_player_session")
		fmt.Println(r.PostForm)
		resp.Data = map[string]string{
			"player_name": username,
			"nickname":    username,
			"currency":    GetConfig().Currency,
		}
	}

	buf, _ := json.Marshal(resp)

	w.Write(buf)
}
