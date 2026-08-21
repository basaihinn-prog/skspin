// package main
//
// import (
//
//	"context"
//	"fmt"
//	"serve/servicejdb/jdbcomm"
//	"time"
//
//	"serve/comm/db"
//	"serve/comm/robot"
//	"serve/servicejdb/jdb_14083/internal"
//	"serve/servicejdb/jdb_14083/internal/gendata"
//
// )
//
//	func main() {
//		var (
//			game      = internal.GameID // 游戏ID
//			line      = internal.Line   // 游戏line, 参照游戏发送doSpin请求里面的‘l’值
//			insertC   = 0.05            // 插入数据库时的投注，参照游戏发送doSpin请求里面的‘c’值
//			normalCnt = 5000            // 需要拉取普通盘的次数
//			gameCnt   = 100             // 需要拉取游戏的次数
//			// 使用最低3个挡位
//			arrC = map[float64]int{ // 游戏下注挡位，k：下注值，v:权重
//				0.05: 50,
//				0.1:  10,
//				0.15: 10,
//				0.2:  5,
//				0.25: 2,
//			}
//			gameWeight = 10 // 如果需要更多机器人拉数据，可以把这个值改大点，建议不修改
//			GBuckets   = gendata.GBuckets
//		)
//		mongoaddr := "mongodb://myAdmin:myAdminPassword1@18.61.185.51:27017/?authSource=admin"
//		db.DialToMongo(mongoaddr, game)
//
//		coll := db.Collection2("robot", "game")
//		coll.InsertOne(context.TODO(), robot.GameInfo{
//			Id:        game,
//			Weight:    int64(gameWeight),
//			StartTime: time.Now().Unix(),
//			IsEnd:     false,
//		})
//
//		go jdbcomm.FetchData(&jdbcomm.FetchDataParam{
//			Game:     game,
//			ArrC:     arrC,
//			GBuckets: GBuckets,
//			InsertC:  insertC,
//			L:        line,
//			Ty:       jdbcomm.GameTypeNormal,
//		})
//		//go jdbcomm.FetchData(&jdbcomm.FetchDataParam{
//		//	Game:     game,
//		//	ArrC:     arrC,
//		//	GBuckets: GBuckets,
//		//	InsertC:  insertC,
//		//	L:        line,
//		//	Ty:       jdbcomm.GameTypeGame,
//		//})
//		needMap := map[int]int64{
//			jdbcomm.GameTypeNormal: int64(normalCnt),
//			jdbcomm.GameTypeGame:   int64(gameCnt),
//		}
//		go jdbcomm.CheckGameCnt(fmt.Sprintf("pp_%s", game), needMap)
//		select {}
//	}
package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main2() {
	// 设置中断信号捕获
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	tk, _, err := jdbcomm.GetTKLocation()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = jdbcomm.Action101(tk)
	if err != nil {
		fmt.Println(err)
		return
	}
	a6, err := jdbcomm.Action6(tk)
	if err != nil {
		fmt.Println(err)
		return
	}
	a19, err := jdbcomm.Action19(tk)
	_ = a19
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = jdbcomm.Action23(tk)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = jdbcomm.Action24(tk)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = jdbcomm.Action13(tk)
	if err != nil {
		fmt.Println(err)
		return
	}
	go jdbcomm.HeartBeat(tk)

	// 编码
	encoded := base64.StdEncoding.EncodeToString([]byte("data"))

	// 解码
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	fmt.Println(string(decoded))
	// 连接到 WebSocket 服务器
	wsurl := "wss://est03.js-mingyi.com/websocket"
	log.Printf("Connecting to %s", wsurl)
	header := http.Header{}
	header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	conn, _, err := websocket.DefaultDialer.Dial(wsurl, header)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// 用于接收服务器消息的 goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%x", message)
			nr, err := jdbcomm.NewFromBinaryData(message)
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			if p, ok := nr.GetSFSObject("p"); ok {
				ptk, ok := p.GetString("tk")
				if !ok {
					return
				}
				ct, ok := p.GetInt("ct")
				if !ok {
					return
				}
				jdbcomm.CompressionThreshold = int(ct)
				if !ok {
					return
				}
				ms, ok := p.GetInt("ms")
				if !ok {
					return
				}
				jdbcomm.MaxMessageSize = int(ms)
				//自己做一个
				so2 := jdbcomm.SFSObject{}
				so2.Init()
				so2.PutByte("c", 0)
				so2.PutShort("a", 1)

				api := jdbcomm.SFSObject{}
				api.Init()
				api.PutString("zn", "JDB_ZONE_GAME")
				api.PutString("un", a6.Data[0].UID)
				hash := md5.Sum([]byte(ptk + "a"))
				api.PutString("pw", hex.EncodeToString(hash[:]))
				so2.PutSFSObject("p", &api)
				fmt.Println(so2)
				ndr, err := so2.ToBinary()
				if err != nil {
					fmt.Println(err)
				}
				conn.WriteMessage(websocket.BinaryMessage, ndr)
			}
		}
	}()
	//// 解码
	//decoded2, err := base64.StdEncoding.DecodeString("gAWWEgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAAlnYW1lTG9naW4AAXIE/////wABcBIAFwADdWlkCAANZGVtbzAwMDM2NkBYWAAIZ2FtZVR5cGUEAAAADgALbWFjaGluZVR5cGUEAAA2ywAGYmFua0lkCAAAAAxzdGFydEJhbGFuY2UHAAAAAAAAAAAABWRlYnVnAQAAB2dhbWVVaWQIAA1kZW1vMDAwMzY2QFhYAAhnYW1lUGFzcwgAB2Y0MWE4Y2QACHVzZXJOYW1lCAANZGVtbzAwMDM2NkBYWAAKc2Vzc2lvbklEMAgAAAAKc2Vzc2lvbklEMQgAAAAKc2Vzc2lvbklEMggAAAAKc2Vzc2lvbklEMwgD0ENENDQxNEMwREIxQzc4MThEOTVCQzY0RDk5NzMzRjZCQkNBN0YxRkNBODI3NDg2NzhFNEQ4MjUzMEJEQUZFNEJFNkI0M0JBMDdBMzA1OENFNUY1RjVGRTM5NTI1RjUwQTY5MUNEQkY0RUM5M0VCMjVENzkzQzZGM0U4REYzQkFFM0QwN0NGQjYwQjIxQzAwMkU5MkY4MjBERTZCNDEwRUE3MEU5RTc5MjcyQTBCQzJDNDU1QjI5RTJDRTQ1OEUyRjE3QTYyMTJDMEJEM0VDMTFGOUQ3QUUzM0ZGOTNBRjJFMDgyQTNCNzIwRTJBODIxRTYxOEMxRDNFMkQ2OTQwMTg4Qzg5QzQxOEVFOUZFRjFEQTg4MkNBNjBEQzZDOEZFMDAyNEE0MTFGOTJCQUU0NTZENkJFRTU3OTA5MTEwRjk1MjFCMkY5NTJFOTlDNjRENTFFMTg5MzBGRkYwNUZFMzBCRTAyQUZGMUYzNkJDMUYwNDc0QThENTQ5NjlFQTFDODg2MjA0REVGQ0VGRDU3MEM1MTk5RkJGM0NGRTI1RTUxRThGOTQ2MTgyQTkyQTZBOEVDMzg0NzRFNjFBRjk5QzZGRUQ3NzZFQkU4MTlDNjZEQUEzMTgxMEM2MDYyRTAxQjlCRERBODdENkMwMDFEODAxRTFFQUZDNUEwNjNBMURFNUMzNkE1QkI2OTVFQ0UwNjZCODhBRTVDNDE4Rjc5QUE5QTQ2QUUyRDk2MUUwODEzNUE3MzgyQzQ0MDc1QjQ0NTZCRTY2MzgwNDI5MThCMTMxMjRCREI0MTY1Qjc2M0Y1NzhEQjIwRUQwNTJGNkEyOEIwRTYzRjYzN0QxM0M4QTJDMUY5NUNEMzY0REFGOEIxODhFQjBBODExOEM0RTlCQ0M1ODdBQzVBNjhCRkMxMkRCMjBFREE0MkFGRURCRjM0MTkwN0UwRDg0MjE4RDFFRjc4NDZGQUNGRkM5QzA0MDIyMjFBOTFEMjZFRDIwMTk1M0JFNTVGRUE4OEEzRUNCMDMxQkUzNDY1QjVDODA4NDM2Qjc3NUE1Q0EwQzM5NkNBMzUxNjcyQkEyQzZBNzRBQUI4NUFBRTM2NkRFOTU0OEU5NEIzMUZCQzREREE1MDREN0ZDNTkwQzlDNzJERTQzRjczMjIwN0ZCQjVBNkIyMDRBMjNBRkQ2MzdGRDYzQjZFODI2MEJGMkRCMDQ0QzQ0OEY2N0YzNzgxNEE2OURCNkM0NEZBNkUwMTc4RUZEREEyQzNBMjFCQTZBNzQ0QkQ4QjY0ODhFNjE4MUJEQkRFMzUACnNlc3Npb25JRDQIAAAABnVzZVNTTAEBAAhwYXNzd29yZAgAAWEACmNsaWVudFR5cGUIAANXZWIAAXQIAAAADWdhbWVMb2dpbk5hbWUIAAlnYW1lTG9naW4ABHpvbmUIAA1KREJfWk9ORV9HQU1FAARwb3J0BAAAAbsABGhvc3QIABNlc3QwMy5qcy1taW5neWkuY29tAAh6b25lTmFtZQgADUpEQl9aT05FX0dBTUU=")
	//if err != nil {
	//}
	////fmt.Println(string(decoded2))
	//fmt.Printf("%x\n", decoded2)
	//so, err := jdbcomm.NewFromBinaryData(decoded2)
	//if err != nil {
	//
	//}
	//b2, err := so.ToBinary()
	//fmt.Printf("%x\n", b2)
	//nodeServer(decoded2)

	//自己做一个
	so2 := jdbcomm.SFSObject{}
	so2.Init()
	so2.PutByte("c", 0)
	so2.PutShort("a", 0)

	api := jdbcomm.SFSObject{}
	api.Init()
	api.PutString("api", "1.7.15")
	api.PutString("cl", "JavaScript")
	so2.PutSFSObject("p", &api)
	fmt.Println(so2)
	ndr, err := so2.ToBinary()
	fmt.Printf("%x\n", ndr)
	//fmt.Printf("%x\n", decoded2)
	conn.WriteMessage(websocket.BinaryMessage, ndr)
	//decoded2, err = base64.URLEncoding.DecodeString("gABjEgADAAFjAgAAAWEDAAEAAXASAAMAAnpuCAANSkRCX1pPTkVfR0FNRQACdW4IAA1kZW1vMDAxNjgzQFhYAAJwdwgAIDM5ZTFhYjY4ZGVjNTcyYzY2OGNjZjhkMWU0ODFkMDlm")
	//if err != nil {
	//}
	//fmt.Printf("%x\n", decoded2)
	//nodeServer(decoded2)

	//var params ParamsData
	//params.C = 0
	//params.A = 1
	//params.P = make(map[string]interface{})
	//params.P["zn"] = "JDB_ZONE_GAME"
	//params.P["un"] = a6.Data[0].UID
	//params.P["pw"], _ = uuid.NewUUID()
	//buf, err := json.Marshal(params)
	//if err != nil {
	//}
	//nodeServerEncode(buf)
	//conn.WriteMessage(websocket.BinaryMessage, buf)
	//conn.WriteMessage(websocket.TextMessage, []byte("gABNEgADAAFwEgADAAJjdAR/////AAJtcwQAB6EgAAJ0awgAIDNhMjc2NDQ0MzQ4NjM5M2IyNDUzZjg3ZWNlY2IxYjhhAAFhAwAAAAFjAgA="))
	//conn.WriteMessage(websocket.TextMessage, []byte("gABjEgADAAFjAgAAAWEDAAEAAXASAAMAAnpuCAANSkRCX1pPTkVfR0FNRQACdW4IAA1kZW1vMDAxNjUyQFhYAAJwdwgAIGQ0MzgzMjk3NGZmNTQ2ZTYxZjJkYWQwYmRlNWVmMjI0"))
	//conn.WriteMessage(websocket.TextMessage, []byte("gAJ4EgADAAFwEgAGAAJycwMAAAACem4IAA1KREJfWk9ORV9HQU1FAAJ1bggADWRlbW8wMDE2NTJAWFgAAnBpAwAAAAJybBEADBEACwQAAAACCAAJU0xPVF9ST09NCAAHZGVmYXVsdAEBAQABAAMGmwMTiBEAAAMAAAMAABEACQQAAAADCAAMUFVTT1lTX0xPQkJZCAAHZGVmYXVsdAEAAQABAAMAEAMTiBEAABEACQQAAAAECAANVE9OR0lUU19MT0JCWQgAB2RlZmF1bHQBAAEAAQADAA8DE4gRAAARAAkEAAAABQgAC1JVTU1ZX0xPQkJZCAAHZGVmYXVsdAEAAQABAAMAAAMTiBEAABEACQQAAAAGCAAMUlVOTklOR19HQU1FCAAHZGVmYXVsdAEAAQABAAMAEQMTiBEAABEACQQAAADnCAAFMTgwMjAIAAdkZWZhdWx0AQABAAEAAwAAAxOIEQAAEQAJBAAAAOgIAAUxODAyMQgAB2RlZmF1bHQBAAEAAQADAAADE4gRAAARAAkEAAAA6QgAC1NJTkdMRV9TUElOCAAHZGVmYXVsdAEAAQABAAMACwMTiBEAABEACQQAAADqCAAFMTgwMjYIAAdkZWZhdWx0AQABAAEAAwAFAxOIEQAAEQAJBAAAAOsIAAVNSU5FUwgAB2RlZmF1bHQBAAEAAQADAFEDE4gRAAARAAkEAAAA7AgAC0NBU0lOT19ST09NCAAHZGVmYXVsdAEAAQABAAMAAAMTiBEAABEACQQAAADtCAAFMTgwMjIIAAdkZWZhdWx0AQABAAEAAwAAAxOIEQAAAAJpZAQAAI5TAAFhAwABAAFjAgA="))

	// 定时发送消息
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			//return
		case t := <-ticker.C:
			fmt.Println("ts" + t.String())
			//// 每秒发送一条消息
			//message := []byte("Hello from client at " + t.String())
			//err := conn.WriteMessage(websocket.TextMessage, message)
			//if err != nil {
			//	log.Println("Write error:", err)
			//	return
			//}
		case <-interrupt:
			// 捕获中断信号，优雅关闭连接
			log.Println("Interrupt received, closing connection...")
			//err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
