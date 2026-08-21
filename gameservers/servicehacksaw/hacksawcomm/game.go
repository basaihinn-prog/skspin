package hacksawcomm

type Game struct {
	ID        string `json:"id"`
	DBName    string `json:"dbname"`
	MType     string `json:"mType"`
	GameType  string `json:"gameType"`
	GName     string `json:"gName"`
	UserName  string `json:"userName"`
	UniqueKey string `json:"uniqueKey"`
	Base      string `json:"base"`
}

//func (s *Game) Run(spinFun func() SFSObject) {
//	// lo.Must(mq.ConnectServerMust("127.0.0.1:11002"))
//	// 设置中断信号捕获
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt)
//
//	tk, location, err := GetHuiDuNewmgcKey(&http.Client{}, s.ID)
//	if err != nil {
//		panic(err)
//	}
//	var GBuckets = GBuckets
//	//mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
//	//db.DialToMongo("mongodb://myAdmin:myAdminPassword1@18.61.185.51:27017/?authSource=admin", lazy.ServiceName)
//	//db.DialToMongo("mongodb://myAdmin:myAdminPassword1@172.31.32.77:27017/?authSource=admin", lazy.ServiceName)
//	db.DialToMongo("mongodb://myAdmin:myAdminPassword1@18.61.185.51:27017/?authSource=admin", lazy.ServiceName)
//	coll := db.Collection2(s.DBName, "simulate")
//	fmt.Println(tk, location)
//	//需要的预处理
//	url2 := "https://eweb03.js-mingyi.com/frontendAPI.do"
//	method := "POST"
//
//	data := url.Values{}
//	data.Set("action", "101")
//	data.Set("x", tk)
//
//	// 编码为字符串
//	payload := strings.NewReader(data.Encode())
//	println(payload)
//	client := &http.Client{}
//	req, err := http.NewRequest(method, url2, payload)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//
//	req.Header.Add("Accept", "application/json, text/plain, */*")
//	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
//	req.Header.Add("Connection", "keep-alive")
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Origin", "https://jifjie.h9vo10aqz.com")
//	req.Header.Add("Referer", "https://jifjie.h9vo10aqz.com/")
//	req.Header.Add("Sec-Fetch-Dest", "empty")
//	req.Header.Add("Sec-Fetch-Mode", "cors")
//	req.Header.Add("Sec-Fetch-Site", "cross-site")
//	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
//	req.Header.Add("jots", "")
//	req.Header.Add("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
//	req.Header.Add("sec-ch-ua-mobile", "?0")
//	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
//
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	defer res.Body.Close()
//
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	fmt.Println(string(body))
//	var resp ResponseAction101
//	err = json.Unmarshal(body, &resp)
//	if err != nil {
//		fmt.Println("解析错误:", err)
//		panic(err)
//		return
//	}
//	if resp.Status != "0000" {
//		fmt.Println("服务器返回错误:", resp.Status)
//	}
//
//	data.Set("action", "6")
//	// 编码为字符串
//	payload = strings.NewReader(data.Encode())
//	println(payload)
//	req, err = http.NewRequest(method, url2, payload)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	req.Header.Add("Accept", "application/json, text/plain, */*")
//	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
//	req.Header.Add("Connection", "keep-alive")
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Origin", "https://jifjie.h9vo10aqz.com")
//	req.Header.Add("Referer", "https://jifjie.h9vo10aqz.com/")
//	req.Header.Add("Sec-Fetch-Dest", "empty")
//	req.Header.Add("Sec-Fetch-Mode", "cors")
//	req.Header.Add("Sec-Fetch-Site", "cross-site")
//	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
//	req.Header.Add("jots", "892e6599-c1ff-497d-aed9-59ab5e2cd956")
//	req.Header.Add("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
//	req.Header.Add("sec-ch-ua-mobile", "?0")
//	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
//
//	res, err = client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	defer res.Body.Close()
//
//	body, err = io.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	fmt.Println(string(body))
//	var a6 ResponseAction6
//	err = json.Unmarshal(body, &a6)
//	if err != nil {
//		fmt.Println("解析错误:", err)
//		panic(err)
//		return
//	}
//	if resp.Status != "0000" {
//		fmt.Println("服务器返回错误:", resp.Status)
//	}
//	fmt.Println(a6)
//
//	data.Set("action", "19")
//	data.Set("gameType", s.GameType)
//	data.Set("mType", s.MType)
//	data.Set("gName", s.GName)
//	data.Set("clientType", "web")
//	data.Set("gameLine", "est03.js-mingyi.com_443_0")
//	// 编码为字符串
//	payload = strings.NewReader(data.Encode())
//	println(payload)
//	req, err = http.NewRequest(method, url2, payload)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	req.Header.Add("Accept", "application/json, text/plain, */*")
//	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
//	req.Header.Add("Connection", "keep-alive")
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Origin", "https://jifjie.h9vo10aqz.com")
//	req.Header.Add("Referer", "https://jifjie.h9vo10aqz.com/")
//	req.Header.Add("Sec-Fetch-Dest", "empty")
//	req.Header.Add("Sec-Fetch-Mode", "cors")
//	req.Header.Add("Sec-Fetch-Site", "cross-site")
//	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
//	req.Header.Add("jots", "892e6599-c1ff-497d-aed9-59ab5e2cd956")
//	req.Header.Add("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
//	req.Header.Add("sec-ch-ua-mobile", "?0")
//	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
//
//	res, err = client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	defer res.Body.Close()
//
//	body, err = io.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	var a19 ResponseAction19
//	err = json.Unmarshal(body, &a19)
//	if err != nil {
//		fmt.Println("解析错误:", err)
//		panic(err)
//		return
//	}
//	if resp.Status != "0000" {
//		fmt.Println("服务器返回错误:", resp.Status)
//	}
//	fmt.Println(a19)
//
//	data.Set("action", "23")
//	data.Set("mType", s.MType)
//	// 编码为字符串
//	payload = strings.NewReader(data.Encode())
//	println(payload)
//	req, err = http.NewRequest(method, url2, payload)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	req.Header.Add("Accept", "application/json, text/plain, */*")
//	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
//	req.Header.Add("Connection", "keep-alive")
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Origin", "https://jifjie.h9vo10aqz.com")
//	req.Header.Add("Referer", "https://jifjie.h9vo10aqz.com/")
//	req.Header.Add("Sec-Fetch-Dest", "empty")
//	req.Header.Add("Sec-Fetch-Mode", "cors")
//	req.Header.Add("Sec-Fetch-Site", "cross-site")
//	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
//	req.Header.Add("jots", "892e6599-c1ff-497d-aed9-59ab5e2cd956")
//	req.Header.Add("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
//	req.Header.Add("sec-ch-ua-mobile", "?0")
//	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
//
//	res, err = client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	defer res.Body.Close()
//
//	body, err = io.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	fmt.Println(string(body))
//
//	data.Set("action", "24")
//	data.Set("mType", s.MType)
//	// 编码为字符串
//	payload = strings.NewReader(data.Encode())
//	println(payload)
//	req, err = http.NewRequest(method, url2, payload)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	req.Header.Add("Accept", "application/json, text/plain, */*")
//	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
//	req.Header.Add("Connection", "keep-alive")
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Origin", "https://jifjie.h9vo10aqz.com")
//	req.Header.Add("Referer", "https://jifjie.h9vo10aqz.com/")
//	req.Header.Add("Sec-Fetch-Dest", "empty")
//	req.Header.Add("Sec-Fetch-Mode", "cors")
//	req.Header.Add("Sec-Fetch-Site", "cross-site")
//	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
//	req.Header.Add("jots", "892e6599-c1ff-497d-aed9-59ab5e2cd956")
//	req.Header.Add("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
//	req.Header.Add("sec-ch-ua-mobile", "?0")
//	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
//
//	res, err = client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	defer res.Body.Close()
//
//	body, err = io.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	fmt.Println(string(body))
//
//	str := `{"level":"GAME","event":"EGRET_READY","message":{"message":"GAME IS READY."},"accessToken":"%s","apiServer":"https://eweb03.js-mingyi.com","gName":"%s",
//				"gameType":"14","mType":"%s","ui_version":"3.150.0","uniqueKey":"%s",
//				"userAgent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36","userName":"%s"}`
//
//	data.Set("action", "13")
//	data.Set("data", fmt.Sprintf(str, tk, s.GName, s.MType, s.UniqueKey, s.UserName))
//	// 编码为字符串
//	payload = strings.NewReader(data.Encode())
//	println(payload)
//	req, err = http.NewRequest(method, url2, payload)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	req.Header.Add("Accept", "application/json, text/plain, */*")
//	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
//	req.Header.Add("Connection", "keep-alive")
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Origin", "https://jifjie.h9vo10aqz.com")
//	req.Header.Add("Referer", "https://jifjie.h9vo10aqz.com/")
//	req.Header.Add("Sec-Fetch-Dest", "empty")
//	req.Header.Add("Sec-Fetch-Mode", "cors")
//	req.Header.Add("Sec-Fetch-Site", "cross-site")
//	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
//	req.Header.Add("jots", "892e6599-c1ff-497d-aed9-59ab5e2cd956")
//	req.Header.Add("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
//	req.Header.Add("sec-ch-ua-mobile", "?0")
//	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
//
//	res, err = client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	defer res.Body.Close()
//
//	body, err = io.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//		return
//	}
//	fmt.Println(string(body))
//
//	hb := func() {
//		data.Set("action", "24")
//		data.Set("mType", "14082")
//		// 编码为字符串
//		payload = strings.NewReader(data.Encode())
//		println(payload)
//		req, err = http.NewRequest(method, url2, payload)
//		if err != nil {
//			fmt.Println(err)
//			panic(err)
//			return
//		}
//		req.Header.Add("Accept", "application/json, text/plain, */*")
//		req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
//		req.Header.Add("Connection", "keep-alive")
//		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//		req.Header.Add("Origin", "https://jifjie.h9vo10aqz.com")
//		req.Header.Add("Referer", "https://jifjie.h9vo10aqz.com/")
//		req.Header.Add("Sec-Fetch-Dest", "empty")
//		req.Header.Add("Sec-Fetch-Mode", "cors")
//		req.Header.Add("Sec-Fetch-Site", "cross-site")
//		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
//		req.Header.Add("jots", "892e6599-c1ff-497d-aed9-59ab5e2cd956")
//		req.Header.Add("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
//		req.Header.Add("sec-ch-ua-mobile", "?0")
//		req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
//
//		res, err = client.Do(req)
//		if err != nil {
//			fmt.Println(err)
//			panic(err)
//			return
//		}
//		defer res.Body.Close()
//
//		body, err = io.ReadAll(res.Body)
//		if err != nil {
//			fmt.Println(err)
//			panic(err)
//			return
//		}
//		fmt.Println(string(body))
//	}
//	go func() {
//		t := time.NewTicker(time.Millisecond * 15000)
//		for {
//			<-t.C
//			hb()
//		}
//	}()
//
//	// 连接到 WebSocket 服务器
//	wsurl := "wss://est03.js-mingyi.com/websocket"
//	log.Printf("Connecting to %s", wsurl)
//	header := http.Header{}
//	header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
//	conn, _, err := websocket.DefaultDialer.Dial(wsurl, header)
//	if err != nil {
//		log.Fatal("Dial error:", err)
//		panic(err)
//	}
//	defer conn.Close()
//	apiFunc := func(conn *websocket.Conn) {
//		//自己做一个
//		so2 := SFSObject{}
//		so2.Init()
//		so2.PutByte("c", 0)
//		so2.PutShort("a", 0)
//
//		api := SFSObject{}
//		api.Init()
//		api.PutString("api", "1.7.15")
//		api.PutString("cl", "JavaScript")
//		so2.PutSFSObject("p", &api)
//		fmt.Printf("so2:%v\n", so2)
//		ndr, _ := so2.ToBinary()
//
//		fmt.Printf("hdr: %x\n", ndr)
//		//fmt.Printf("%x\n", decoded2)
//		conn.WriteMessage(websocket.BinaryMessage, ndr)
//	}
//	// 用于接收服务器消息的 goroutine
//	done := make(chan struct{})
//	spinChannel := make(chan bool)
//
//	go func() {
//		defer close(done)
//		for {
//			_, message, err := conn.ReadMessage()
//			if err != nil {
//				panic(err)
//				return
//			}
//			//fmt.Printf("%x", message)
//			nr, err := NewFromBinaryData(message)
//			if err != nil {
//				log.Println("Read error:", err)
//				panic(err)
//				return
//			}
//			a, ok := nr.GetShort("a")
//			//心跳
//			if ok && a == 1001 {
//				continue
//			}
//			if p, ok := nr.GetSFSObject("p"); ok {
//				action, ok := nr.GetShort("a")
//				//存在0
//				if ok {
//					//version + javascript
//					if action == 0 {
//						ptk, ok := p.GetString("tk")
//						if !ok {
//							return
//						}
//						ct, ok := p.GetInt("ct")
//						if !ok {
//							return
//						}
//						CompressionThreshold = int(ct)
//						if !ok {
//							return
//						}
//						ms, ok := p.GetInt("ms")
//						if !ok {
//							return
//						}
//						MaxMessageSize = int(ms)
//						//自己做一个
//						so2 := SFSObject{}
//						so2.Init()
//						so2.PutByte("c", 0)
//						so2.PutShort("a", 1)
//
//						api := SFSObject{}
//						api.Init()
//						api.PutString("zn", "JDB_ZONE_GAME")
//						api.PutString("un", a6.Data[0].UID)
//						hash := md5.Sum([]byte(ptk + "a"))
//						api.PutString("pw", hex.EncodeToString(hash[:]))
//						so2.PutSFSObject("p", &api)
//						fmt.Println(so2)
//						ndr, err := so2.ToBinary()
//						if err != nil {
//							fmt.Println(err)
//							panic(err)
//						}
//						conn.WriteMessage(websocket.BinaryMessage, ndr)
//					}
//					//login
//					if action == 1 {
//						//{c 1,a 13,p{c "gameLogin",r -1,p{}}}
//						//so := hacksawcomm.CreateSFSObject(1, 13)
//						so := SFSObject{}
//						so.Init()
//						so.PutByte("c", 1)
//						so.PutShort("a", 13)
//						p := SFSObject{}
//						p.Init()
//						p.PutString("c", "gameLogin")
//						p.PutInt("r", -1)
//
//						pp := SFSObject{}
//						pp.Init()
//						pp.PutString("uid", a19.Data.Result6.UID)
//						pp.PutInt("gameType", int32(a19.Data.Result10.GameType))
//						pp.PutInt("machineType", int32(a19.Data.Result10.MachineType))
//						pp.PutString("bankId", "")
//						pp.PutDouble("startBalance", 0.0)
//						pp.PutBool("debug", false)
//						pp.PutString("gameUid", a19.Data.Result10.GameUid)
//						pp.PutString("gamePass", a19.Data.Result10.GamePass)
//						pp.PutString("userName", a19.Data.Result6.UID)
//						pp.PutString("sessionID0", a19.Data.Result10.S0)
//						pp.PutString("sessionID1", a19.Data.Result10.S1)
//						pp.PutString("sessionID2", a19.Data.Result10.S2)
//						pp.PutString("sessionID3", a19.Data.Result10.S3)
//						pp.PutString("sessionID4", a19.Data.Result10.S4)
//						pp.PutBool("useSSL", a19.Data.Result10.UseSSL)
//						pp.PutString("password", "a")
//						pp.PutString("clientType", "Web")
//						pp.PutString("t", "")
//						pp.PutString("gameLoginName", "gameLogin")
//						pp.PutString("zone", a19.Data.Result10.Zone)
//						pp.PutInt("port", 443)
//						pp.PutString("host", "est03.js-mingyi.com")
//						pp.PutString("zoneName", a19.Data.Result10.Zone)
//						p.PutSFSObject("p", &pp)
//						so.PutSFSObject("p", &p)
//						bytes, err := so.ToBinary()
//						if err != nil {
//							fmt.Println(err)
//							panic(err)
//						}
//						//返回加入房间
//						conn.WriteMessage(websocket.BinaryMessage, bytes)
//					}
//					//加入房间后
//					if action == 13 {
//						c, _ := p.GetString("c")
//						if c == "gameLoginReturn" {
//							so := CreateH5Init()
//							bytes, err := so.ToBinary()
//							if err != nil {
//								fmt.Println(err)
//								panic(err)
//							}
//							conn.WriteMessage(websocket.BinaryMessage, bytes)
//							//go func() {
//							//	for {
//							//		so2 := hacksawcomm.CreateHeartBeat()
//							//		bytes2, err := so2.ToBinary()
//							//		if err != nil {
//							//			fmt.Println(err)
//							//		}
//							//		conn.WriteMessage(websocket.BinaryMessage, bytes2)
//							//		time.Sleep(time.Millisecond * 20000)
//							//	}
//							//
//							//}()
//							spinChannel <- true
//						}
//						//
//						if c == "h5.spinResponse" {
//							pp, _ := p.GetSFSObject("p")
//							m := Variables{}
//							bytes, _ := pp.GetByteArray("entity")
//							err = json.Unmarshal(bytes, &m)
//							if err != nil {
//								fmt.Println(err)
//								panic(err)
//								return
//							}
//							fmt.Printf("%v", m)
//							totalWin, ok := FindKeyInNestedMap(m, "totalWin")
//							_, hasGame := FindKeyInNestedMap(m, "specialHitPattern")
//							displayBet, _ := FindKeyInNestedMap(m, "displayBet")
//							times := 0.0
//							if totalWin != 0 && ok {
//								times = ut.Round6(totalWin.(float64) / displayBet.(float64))
//							}
//							fmt.Println(m)
//							objectId := primitive.NewObjectID()
//							doc := SimulateData{
//								Id:                   objectId,
//								DropPan:              m,
//								HasGame:              hasGame,
//								Times:                times,
//								BucketId:             GBuckets.GetBucket(times, hasGame, 0),
//								Type:                 0,
//								Selected:             true,
//								BucketHeartBeat:      1,
//								BucketWave:           1,
//								BucketGov:            1,
//								BucketMix:            1,
//								BucketStable:         1,
//								BucketHighAward:      1,
//								BucketSuperHighAward: 1,
//							}
//							coll.InsertOne(context.TODO(), doc)
//							fmt.Println("success -------->", s.DBName, objectId.Hex(), "type", 0)
//							time.Sleep(time.Second * 2)
//							spinChannel <- true
//						}
//
//					}
//				}
//			}
//		}
//	}()
//	// 解码
//	/*	decoded2, err := base64.StdEncoding.DecodeString(s.Base)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Printf("%x\n", decoded2)
//		so, err := NewFromBinaryData(decoded2)
//		if err != nil {
//			panic(err)
//		}
//		b2, err := so.ToBinary()
//		fmt.Printf("%x\n", b2)*/
//	//nodeServer(decoded2)
//
//	apiFunc(conn)
//
//	// 定时发送消息
//	ticker := time.NewTicker(time.Second)
//	defer ticker.Stop()
//
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Println("Recovered in f", r)
//		}
//	}()
//	for {
//		select {
//		case spin := <-spinChannel:
//			if spin {
//				so := spinFun()
//				bytes, err := so.ToBinary()
//				if err != nil {
//					fmt.Println(err)
//					panic(err)
//				}
//				err = conn.WriteMessage(websocket.BinaryMessage, bytes)
//				if err != nil {
//					conn, _, err = websocket.DefaultDialer.Dial(wsurl, header)
//					apiFunc(conn)
//					panic(err)
//				}
//			}
//		case <-done:
//			panic("socket break")
//		case <-interrupt:
//			// 捕获中断信号，优雅关闭连接
//			log.Println("Interrupt received, closing connection...")
//			//err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
//			if err != nil {
//				log.Println("Write close error:", err)
//				panic(err)
//				return
//			}
//			select {
//			case <-done:
//			case <-time.After(time.Second):
//			}
//			return
//		}
//	}
//}
