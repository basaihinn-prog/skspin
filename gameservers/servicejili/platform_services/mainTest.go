package main

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"serve/comm/db"
//	"serve/comm/define"
//	"serve/comm/mux"
//	"serve/comm/ut"
//	"serve/servicejili/jiliut"
//	"serve/servicejili/platform_services/internal"
//	"serve/servicejili/platform_services/internal/models"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"io"
//	"log"
//	"log/slog"
//	"net/url"
//	"strconv"
//
//	"net/http"
//)
//
//// ResponseData 定义一个结构体，用于接收JSON数据
//type ResponseData struct {
//	Code    int           `json:"code"`
//	Message string        `json:"message"`
//	Player  models.Player `json:"player"` //用户信息
//}
//
//type RespData struct {
//	Code  int    `json:"code"`
//	Error string `json:"error"`
//	Data  any    `json:"data"`
//}
//
//type LaunchGameRequest struct {
//	UID  string `json:"UID"`
//	Game string `json:"Game"`
//	Lang string `json:"Lang"`
//}
//
//type LaunchGameResp struct {
//	Code  int    `json:"code"`
//	Error string `json:"error"`
//	Data  struct {
//		Url string `json:"Url"`
//	} `json:"data"`
//}
//
//const RPApiUrl_test = "https://gamecenter.rpgamestest.com/"
//
//func main() {
//
//	// 设置路由
//	//http.HandleFunc("/post", handlePost)
//	http.HandleFunc("/testLogin", CORSMiddleware)
//
//	http.HandleFunc("/testPPGameList", CORSGetPPGameList)
//
//	http.HandleFunc("/testLaunchGamePP", CORSLaunchGamePP)
//
//	http.HandleFunc("/fundTransferInPP", CORSFundTransferIn)
//
//	http.HandleFunc("/fundTransferOutPP", CORSFundTransferOut)
//
//	http.HandleFunc("/getPlayerCold", CORSGetPlayerCold)
//
//	// 启动服务器
//	fmt.Println("Server is listening on http://localhost:8080")
//	if err := http.ListenAndServe(":8080", nil); err != nil {
//		fmt.Printf("Error starting server: %s\n", err)
//	}
//
//}
//
//func findFromPlayer(account string, password string) (playerResp *models.Player, err error) {
//	mongoaddr := jiliut.GetFetchMongoAddr()
//	db.DialToMongo(mongoaddr, internal.GameID)
//
//	one := &models.Player{}
//	err = db.Collection("players").FindOne(context.TODO(), bson.M{"account": account, "password": password}).Decode(one)
//	if err != nil {
//		return
//	}
//	playerResp = one
//	return
//}
//
//func updatePlayer(uid string, playergold float64) {
//	mongoaddr := jiliut.GetFetchMongoAddr()
//
//	// 设置MongoDB连接
//	clientOptions := options.Client().ApplyURI(mongoaddr)
//	client, err := mongo.Connect(context.TODO(), clientOptions)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 检查连接
//	err = client.Ping(context.TODO(), nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Connected to MongoDB!")
//
//	// 选择数据库和集合
//	collection := client.Database(internal.GameID).Collection("players")
//
//	// 准备更新条件（过滤器）和更新内容
//	filter := bson.D{{"uid", uid}}                                 // 查询条件
//	update := bson.D{{"$set", bson.D{{"playergold", playergold}}}} // 更新的字段和值
//
//	// 执行更新操作
//	result, err := collection.UpdateOne(context.TODO(), filter, update)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 输出更新结果
//	fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
//
//	// 关闭MongoDB连接
//	err = client.Disconnect(context.TODO())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Connection to MongoDB closed.")
//
//}
//
//func findPlayerByUId(uid string) (playerResp *models.Player, err error) {
//	mongoaddr := jiliut.GetFetchMongoAddr()
//	db.DialToMongo(mongoaddr, internal.GameID)
//
//	one := &models.Player{}
//	err = db.Collection("players").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(one)
//	if err != nil {
//		return
//	}
//	playerResp = one
//	return
//}
//
//func findPlayer(playerName string) (b bool) {
//	mongoaddr := jiliut.GetFetchMongoAddr()
//	db.DialToMongo(mongoaddr, internal.GameID)
//
//	one := &models.Player{}
//	db.Collection("players").FindOne(context.TODO(), bson.M{"account": playerName}).Decode(one)
//
//	fmt.Println(one)
//
//	if one.Account != "" {
//		return true
//	} else {
//		return false
//	}
//
//}
//
//func handlePost(w http.ResponseWriter, r *http.Request) {
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Error reading request body", http.StatusBadRequest)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println("Error:", err)
//		}
//	}(r.Body)
//
//	var playerInfo models.Player
//	// 解析JSON数据
//
//	err = json.Unmarshal(body, &playerInfo)
//	if err != nil {
//		return
//	}
//	fmt.Println("playerInfo:", playerInfo)
//
//	var playerResp *models.Player
//	playerResp, err = findFromPlayer(playerInfo.Account, playerInfo.Password)
//
//	var responseData ResponseData
//	if playerResp == nil {
//		fmt.Println("没有找到值")
//		//判断是否有账号
//		boolean := findPlayer(playerInfo.Account)
//		if boolean == true {
//			//密码错误
//			responseData.Message = "Password error！"
//			responseData.Code = -199
//			fmt.Println("密码错误")
//		} else {
//			//没有账号
//			fmt.Println("没有账号")
//			responseData.Message = "No account！"
//			responseData.Code = -200
//		}
//
//	} else {
//		fmt.Println("playerResp:", playerResp)
//		responseData.Player.PlayerId = playerResp.PlayerId
//		responseData.Player.Account = playerResp.Account
//		responseData.Player.UId = playerResp.UId
//		responseData.Player.UniqueId = playerResp.UniqueId
//		responseData.Player.PlayerGold = playerResp.PlayerGold
//		responseData.Player.UserName = playerResp.UserName
//		responseData.Code = 200
//		responseData.Message = "success"
//	}
//
//	//金币带出
//	goldOutResp := ppGoldOut(responseData.Player.UId)
//	if goldOutResp.Data.Amount > 0 {
//		//更新数据库
//		updatePlayer(responseData.Player.UId, goldOutResp.Data.Amount)
//
//		//更新金币数据
//		responseData.Player.PlayerGold = goldOutResp.Data.Amount
//	}
//
//	// 设置响应头
//	w.Header().Set("Content-Type", "application/json")
//	response, err := json.Marshal(responseData)
//	if err != nil {
//		http.Error(w, "Error creating response", http.StatusInternalServerError)
//		return
//	}
//	// 写入响应体
//	_, err = w.Write(response)
//	if err != nil {
//		return
//	}
//}
//
//func CORSGetPlayerCold(w http.ResponseWriter, r *http.Request) {
//	// 设置CORS响应头
//	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
//	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//
//	// 调用下一个处理器
//	getPlayerCold(w, r)
//}
//
//func getPlayerCold(w http.ResponseWriter, r *http.Request) {
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Error reading request body", http.StatusBadRequest)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println("Error:", err)
//		}
//	}(r.Body)
//
//	var playerInfo models.Player
//	// 解析JSON数据
//
//	err = json.Unmarshal(body, &playerInfo)
//	if err != nil {
//		return
//	}
//	fmt.Println("playerInfo:", playerInfo)
//
//	//金币带出
//	goldOutResp := ppGoldOut(playerInfo.UId)
//	if goldOutResp.Data.Amount > 0 {
//		//更新数据库
//		updatePlayer(playerInfo.UId, goldOutResp.Data.Amount)
//
//		//更新金币数据
//		playerInfo.PlayerGold = goldOutResp.Data.Amount
//	} else {
//		p, err := findPlayerByUId(playerInfo.UId)
//		if err != nil {
//			return
//		}
//		playerInfo.PlayerGold = p.PlayerGold
//	}
//
//	// 设置响应头
//	w.Header().Set("Content-Type", "application/json")
//	response, err := json.Marshal(playerInfo)
//	if err != nil {
//		http.Error(w, "Error creating response", http.StatusInternalServerError)
//		return
//	}
//	// 写入响应体
//	_, err = w.Write(response)
//	if err != nil {
//		return
//	}
//}
//
//func CORSMiddleware(w http.ResponseWriter, r *http.Request) {
//	// 设置CORS响应头
//	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
//	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//
//	// 调用下一个处理器
//	handlePost(w, r)
//}
//
//func CORSGetPPGameList(w http.ResponseWriter, r *http.Request) {
//	// 设置CORS响应头
//	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
//	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//
//	// 调用下一个处理器
//	handlePPGameList(w, r)
//}
//
//func GetPPGameList() string {
//	// 目标URL
//	postUrl := "https://plats.rpgamestest.com/plat/PP/GetGameList"
//
//	// 发起GET请求
//	resp, err := http.Get(postUrl)
//	if err != nil {
//		// 处理错误
//		fmt.Println("Error:", err)
//		return ""
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println("Error:", err)
//		}
//	}(resp.Body) // 确保在函数返回时关闭响应体
//
//	// 读取响应体
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		// 处理读取响应体时的错误
//		fmt.Println("Error reading response body:", err)
//		return ""
//	}
//
//	return string(body)
//}
//
//// 拉取PP的游戏列表
//func handlePPGameList(w http.ResponseWriter, r *http.Request) {
//	_, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Error reading request body", http.StatusBadRequest)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//		}
//	}(r.Body)
//
//	ppGameList := GetPPGameList()
//
//	// 设置响应头
//	w.Header().Set("Content-Type", "application/json")
//	response, err := json.Marshal(ppGameList)
//	if err != nil {
//		http.Error(w, "Error creating response", http.StatusInternalServerError)
//		return
//	}
//	// 写入响应体
//	_, err = w.Write(response)
//	if err != nil {
//		return
//	}
//}
//
//type CreatePP struct {
//	UID string `json:"UID"`
//}
//
//// pp创建用户账号
//func ppCreatePlayer(uid string) {
//	var createPP CreatePP
//	createPP.UID = uid
//	jsonData, err := json.Marshal(createPP)
//	jsonReader := bytes.NewReader(jsonData)
//	//发送post请求
//	req, err := http.NewRequest("POST", "https://plats.rpgamestest.com/plat/PP/Regist", jsonReader)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	//设置请求头
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "admin")
//	req.Header.Set("Pid", "100001")
//	req.Header.Set("Appid", "faketrans")
//	req.Header.Set("Appsecret", "b6337af9-a91a-4085-b1f2-466923470735")
//	//发起请求
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}(resp.Body)
//	// 读取响应体
//	var body []byte
//	body, err = io.ReadAll(resp.Body)
//	fmt.Println(string(body))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//}
//
//// ////////
//type ppGoldInStruct struct {
//	UID    string  `json:"UID"`
//	Amount float64 `json:"Amount"`
//}
//
//// pp金币带入
//func ppGoldIn(uid string, amount float64) {
//	var ppGoldInStruct ppGoldInStruct
//	ppGoldInStruct.UID = uid
//	ppGoldInStruct.Amount = amount
//	jsonData, err := json.Marshal(ppGoldInStruct)
//	jsonReader := bytes.NewReader(jsonData)
//	//发送post请求
//	req, err := http.NewRequest("POST", "https://plats.rpgamestest.com/plat/PP/FundTransferIn", jsonReader)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	//设置请求头
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "admin")
//	req.Header.Set("Pid", "100001")
//	//req.Header.Set("Appid", "faketrans")
//	//req.Header.Set("Appsecret", "b6337af9-a91a-4085-b1f2-466923470735")
//	//发起请求
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}(resp.Body)
//	// 读取响应体
//	var body []byte
//	body, err = io.ReadAll(resp.Body)
//	fmt.Println(string(body))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//}
//
//type PPGoldOutResp struct {
//	Code  int    `json:"code"`
//	Error string `json:"error"`
//	Data  struct {
//		Amount float64 `json:"Amount"`
//		Status string  `json:"Status"`
//	} `json:"data"`
//}
//
//// pp金币带出
//func ppGoldOut(uid string) *PPGoldOutResp {
//	var createPP CreatePP
//	createPP.UID = uid
//	jsonData, err := json.Marshal(createPP)
//	jsonReader := bytes.NewReader(jsonData)
//	//发送post请求
//	req, err := http.NewRequest("POST", "https://plats.rpgamestest.com/plat/PP/FundTransferOut", jsonReader)
//	if err != nil {
//		fmt.Println(err)
//	}
//	//设置请求头
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "admin")
//	req.Header.Set("Pid", "100001")
//	//req.Header.Set("Appid", "faketrans")
//	//req.Header.Set("Appsecret", "b6337af9-a91a-4085-b1f2-466923470735")
//	//发起请求
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}(resp.Body)
//	// 读取响应体
//	var body []byte
//	body, err = io.ReadAll(resp.Body)
//	fmt.Println(string(body))
//
//	var ppGoldOutResp PPGoldOutResp
//	err = json.Unmarshal(body, &ppGoldOutResp)
//
//	if err != nil {
//		fmt.Println(err)
//	}
//	return &ppGoldOutResp
//
//}
//
////
//
//func CORSLaunchGamePP(w http.ResponseWriter, r *http.Request) {
//	// 设置CORS响应头
//	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
//	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//
//	// 调用下一个处理器
//	launchGamePP(w, r)
//}
//
//func launchGamePP(w http.ResponseWriter, r *http.Request) {
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Error reading request body", http.StatusBadRequest)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println("launchGamePP Error:", err)
//		}
//	}(r.Body)
//
//	fmt.Println("sssss", string(body))
//
//	//请求转化成json
//	var launchGameRequest LaunchGameRequest
//	fmt.Println("body==", string(body))
//	err = json.Unmarshal(body, &launchGameRequest)
//	if err != nil {
//		fmt.Println("launchGameRequest Error:", err)
//		return
//	}
//	//查询用户剩余的金币
//	player, err := findPlayerByUId(launchGameRequest.UID)
//	if err != nil {
//		fmt.Println("error:", err)
//		return
//	}
//
//	//创建用户
//	ppCreatePlayer(player.UId)
//
//	//带入金币
//	ppGoldIn(player.UId, player.PlayerGold)
//
//	game, err := LaunchGame(player.UId, launchGameRequest.Game, launchGameRequest.Lang)
//	if err != nil {
//		return
//	}
//	fmt.Println("game=====", game)
//	var launchGameResp LaunchGameResp
//	launchGameResp.Code = 0
//	launchGameResp.Data.Url = game
//	jsonData, err := json.Marshal(launchGameResp)
//
//	// 设置响应头
//	w.Header().Set("Content-Type", "application/json")
//
//	// 写入响应体
//	_, err = w.Write(jsonData)
//	if err != nil {
//		return
//	}
//}
//
///////////
//
//func CORSFundTransferIn(w http.ResponseWriter, r *http.Request) {
//	// 设置CORS响应头
//	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
//	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//
//	// 调用下一个处理器
//	ppFundTransferIn(w, r)
//}
//
//// PP游戏带入金币
//func ppFundTransferIn(w http.ResponseWriter, r *http.Request) {
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Error reading request body", http.StatusBadRequest)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println("Error:", err)
//		}
//	}(r.Body)
//
//	// 发送post请求
//	req, err := http.NewRequest("POST", "https://plats.rpgamestest.com/plat/PP/FundTransferIn", bytes.NewBuffer(body))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// 设置请求头
//	req.Header.Set("Content-Type", "text/plain")
//	req.Header.Set("Authorization", "admin")
//	req.Header.Set("Pid", "100001")
//	req.Header.Set("Appid", "faketrans")
//	req.Header.Set("Appsecret", "b6337af9-a91a-4085-b1f2-466923470735")
//
//	// 发起请求
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}(resp.Body)
//
//	// 读取响应体
//	body, err = io.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// 打印响应体
//	fmt.Println("==============", string(body))
//
//	// 设置响应头
//	w.Header().Set("Content-Type", "application/json")
//
//	// 写入响应体
//	_, err = w.Write(body)
//	if err != nil {
//		return
//	}
//}
//
///////////PP带出
//
//func CORSFundTransferOut(w http.ResponseWriter, r *http.Request) {
//	// 设置CORS响应头
//	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
//	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//
//	// 调用下一个处理器
//	ppFundTransferOut(w, r)
//}
//
//// PP游戏带入金币
//func ppFundTransferOut(w http.ResponseWriter, r *http.Request) {
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Error reading request body", http.StatusBadRequest)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println("Error:", err)
//		}
//	}(r.Body)
//
//	// 发送post请求
//	req, err := http.NewRequest("POST", "https://plats.rpgamestest.com/plat/PP/FundTransferOut", bytes.NewBuffer(body))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// 设置请求头
//	req.Header.Set("Content-Type", "text/plain")
//	req.Header.Set("Authorization", "admin")
//	req.Header.Set("Pid", "100001")
//	req.Header.Set("Appid", "faketrans")
//	req.Header.Set("Appsecret", "b6337af9-a91a-4085-b1f2-466923470735")
//
//	// 发起请求
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}(resp.Body)
//
//	// 读取响应体
//	body, err = io.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// 打印响应体
//	fmt.Println("==============", string(body))
//
//	// 设置响应头
//	w.Header().Set("Content-Type", "application/json")
//
//	// 写入响应体
//	_, err = w.Write(body)
//	if err != nil {
//		return
//	}
//}
//
//func LaunchGame(username, game, lang string) (url string, err error) {
//	var ret struct {
//		GameURL string
//	}
//	err = invoke("/game/start/", map[string]string{
//		"externalPlayerId": username,
//		"gameId":           game,
//		"language":         lang,
//		"platform":         "MOBILE",
//	}, &ret)
//	url = ret.GameURL
//	return
//}
//
//type callResponse struct {
//	Error       string
//	Description string
//}
//
//type PPConfig struct {
//	ApiUrl      string
//	IconUrl     string
//	SecureLogin string
//	SecretKey   string
//}
//
//var ppConfig = PPConfig{
//	ApiUrl:      "https://api.prerelease-env.biz/IntegrationService/v3/http/CasinoGameAPI",
//	IconUrl:     "https://happycasino.prerelease-env.biz/game_pic/square/200/%s.png",
//	SecureLogin: "hllgd_hollygod",
//	SecretKey:   "C1Bd51884d654e10",
//}
//
//func invoke(method string, args map[string]string, result interface{}) (err error) {
//
//	cfg := ppConfig
//	var data = url.Values{}
//	data.Set("secureLogin", cfg.SecureLogin)
//	//THB
//
//	// args := platcomm.GetArgs(ps)
//	for k, v := range args {
//		data.Set(k, v)
//	}
//
//	content := data.Encode() + cfg.SecretKey
//
//	data.Set("hash", ut.ToMD5Str(content))
//
//	fields := slog.With(
//		"plat", "pp",
//		"url", cfg.ApiUrl+method,
//		"reqpayload", data.Encode(),
//	)
//	defer func() {
//		// fields["error"] = err
//		// logrus.WithFields(fields).Info("invoke!!")
//
//		fields.Info("invoke!!", "error", err)
//	}()
//
//	resp, err := http.PostForm(cfg.ApiUrl+method, data)
//	if err != nil {
//		return
//	}
//	fields = fields.With("httpstatus", resp.Status)
//
//	defer resp.Body.Close()
//
//	ret, err := io.ReadAll(resp.Body)
//	if err != nil {
//		fields = fields.With("httpstatus", resp.Status)
//		return
//	}
//	fields = fields.With("respbody", mux.TruncMsg(ret, 1024))
//
//	var callRet callResponse
//	err = json.Unmarshal(ret, &callRet)
//	if err != nil {
//		return err
//	}
//	if callRet.Error != "0" {
//		errcode, _ := strconv.Atoi(callRet.Error)
//		// return logic.CodeError{
//		// 	Code:    errcode,
//		// 	Message: callRet.Description,
//		// }
//		return define.NewErrCode(callRet.Description, errcode)
//	}
//
//	if result != nil {
//		err = json.Unmarshal(ret, result)
//	}
//	return
//}
