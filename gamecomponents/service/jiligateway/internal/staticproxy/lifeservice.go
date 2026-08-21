package staticproxy

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func init() {

	uat_wbslot_fd_mux.HandleFunc("/lifeservice/ws2", lifeservice)
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// 'wss://uat-wbslot-fd-jlfafafa1.kafa010.com/lifeservice/ws2'
func lifeservice(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	_, _, err = c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	c.WriteJSON(M{"error": 0})

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s, type: %v", message, mt)
		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}
