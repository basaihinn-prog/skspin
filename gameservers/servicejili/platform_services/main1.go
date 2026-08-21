package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/url"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/mux"
	"serve/comm/plats/platcomm"
	"serve/comm/ut"
	"serve/servicejili/jiliut"
	"serve/servicejili/platform_services/internal"
	"serve/servicejili/platform_services/internal/models"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"net/http"
)

// ResponseData 定义一个结构体，用于接收JSON数据
type ResponseData struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Player  models.Player `json:"player"` //用户信息
}

type GoldIn struct {
	UID    string  `json:"UID"`
	Amount float64 `json:"Amount"`
	Type   string  `json:"Type"`
}

type GoldOut struct {
	UID  string `json:"UID"`
	Type string `json:"Type"`
}

type GetGameListRequest struct {
	Type string `json:"Type"`
}

type LaunchGameRequest struct {
	UID  string `json:"UID"`
	Game string `json:"Game"`
	Lang string `json:"Lang"`
	Type string `json:"Type"`
}

type LaunchGameResp struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Data  struct {
		Url string `json:"Url"`
	} `json:"data"`
}

type GameInfo struct {
	Plat string `json:"Plat"`
	ID   string `json:"ID"`
	Name string `json:"Name"`
	Type int32  `json:"Type"`
	Icon string `json:"Icon"`
}

type GameListRequest struct {
	Code  int32                 `json:"code"`
	Error string                `json:"error"`
	Data  map[string][]GameInfo `json:"data"`
}

const RPApiUrl_test = "https://gamecenter.rpgamestest.com/"

func main() {

	// 设置路由
	//http.HandleFunc("/post", handlePost)
	http.HandleFunc("/login", CORSMiddleware)

	http.HandleFunc("/getGameList", CORSGetGameList)

	http.HandleFunc("/launchGame", CORSLaunchGame)

	//http.HandleFunc("/fundTransferIn", CORSFundTransferIn)

	//http.HandleFunc("/fundTransferOut", CORSFundTransferOut)

	http.HandleFunc("/getPlayerCold", CORSGetPlayerCold)

	// 启动服务器
	fmt.Println("Server is listening on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

}

func findFromPlayer(account string, password string) (playerResp *models.Player, err error) {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)

	one := &models.Player{}
	err = db.Collection("players").FindOne(context.TODO(), bson.M{"account": account, "password": password}).Decode(one)
	if err != nil {
		return
	}
	playerResp = one
	return
}

func updatePlayer(uid string, playergold float64) {
	mongoaddr := jiliut.GetFetchMongoAddr()

	// 设置MongoDB连接
	clientOptions := options.Client().ApplyURI(mongoaddr)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// 选择数据库和集合
	collection := client.Database(internal.GameID).Collection("players")

	// 准备更新条件（过滤器）和更新内容
	filter := bson.D{{"uid", uid}}                                 // 查询条件
	update := bson.D{{"$set", bson.D{{"playergold", playergold}}}} // 更新的字段和值

	// 执行更新操作
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	// 输出更新结果
	fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)

	// 关闭MongoDB连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}

func findPlayerByUId(uid string) (playerResp *models.Player, err error) {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)

	one := &models.Player{}
	err = db.Collection("players").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(one)
	if err != nil {
		return
	}
	playerResp = one
	return
}

func findPlayer(playerName string) (b bool) {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)

	one := &models.Player{}
	db.Collection("players").FindOne(context.TODO(), bson.M{"account": playerName}).Decode(one)

	fmt.Println(one)

	if one.Account != "" {
		return true
	} else {
		return false
	}

}

func handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(r.Body)

	var playerInfo models.Player
	// 解析JSON数据

	err = json.Unmarshal(body, &playerInfo)
	if err != nil {
		return
	}
	fmt.Println("playerAccount:", playerInfo.Account)

	var playerResp *models.Player
	playerResp, err = findFromPlayer(playerInfo.Account, playerInfo.Password)

	var responseData ResponseData
	if playerResp == nil {
		fmt.Println("没有找到值")
		//判断是否有账号
		boolean := findPlayer(playerInfo.Account)
		if boolean == true {
			//密码错误
			responseData.Message = "Password error！"
			responseData.Code = -199
			fmt.Println("密码错误")
		} else {
			//没有账号
			fmt.Println("没有账号")
			responseData.Message = "No account！"
			responseData.Code = -200
		}

	} else {
		//意外退出
		var goldOutResp PPGoldOutResp

		goldOutResp.Data.Amount, goldOutResp.Data.Status = FundTransferOut(playerResp.UId)

		if goldOutResp.Data.Amount == 0 {
			goldOutResp = goldOut(playerResp.UId, "PG")
		}

		if goldOutResp.Data.Amount > 0 {

			fmt.Printf("登录 playerInfo.UId:%s ,goldOutResp.Data.Amount:%f", playerInfo.UId, goldOutResp.Data.Amount)

			//更新数据库
			updatePlayer(playerResp.UId, goldOutResp.Data.Amount)
			//更新金币数据
			playerResp.PlayerGold = goldOutResp.Data.Amount
		}

		responseData.Player.PlayerId = playerResp.PlayerId
		responseData.Player.Account = playerResp.Account
		responseData.Player.UId = playerResp.UId
		responseData.Player.UniqueId = playerResp.UniqueId
		responseData.Player.PlayerGold = playerResp.PlayerGold
		responseData.Player.UserName = playerResp.UserName
		responseData.Code = 200
		responseData.Message = "success"
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}
	// 写入响应体
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

func CORSGetPlayerCold(w http.ResponseWriter, r *http.Request) {
	// 设置CORS响应头
	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// 调用下一个处理器
	getPlayerCold(w, r)
}

func getPlayerCold(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(r.Body)

	var playerInfo models.Player
	// 解析JSON数据

	err = json.Unmarshal(body, &playerInfo)
	if err != nil {
		return
	}
	//金币带出

	var goldOutResp PPGoldOutResp

	if playerInfo.Type == "PP" {
		goldOutResp.Data.Amount, goldOutResp.Data.Status = FundTransferOut(playerInfo.UId)
	} else {
		goldOutResp = goldOut(playerInfo.UId, playerInfo.Type)
	}

	if goldOutResp.Data.Amount >= 0 {
		//更新数据库

		fmt.Printf("playerInfo.UId:%s ,goldOutResp.Data.Amount:%f", playerInfo.UId, goldOutResp.Data.Amount)

		updatePlayer(playerInfo.UId, goldOutResp.Data.Amount)

		//更新金币数据
		playerInfo.PlayerGold = goldOutResp.Data.Amount
	} else {
		p, err := findPlayerByUId(playerInfo.UId)
		if err != nil {
			return
		}
		playerInfo.PlayerGold = p.PlayerGold
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(playerInfo)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}
	// 写入响应体
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

func CORSMiddleware(w http.ResponseWriter, r *http.Request) {
	// 设置CORS响应头
	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// 调用下一个处理器
	handlePost(w, r)
}

func CORSGetGameList(w http.ResponseWriter, r *http.Request) {
	// 设置CORS响应头
	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// 调用下一个处理器
	handlePPGameList(w, r)
}

func GetGameList(s string) string {
	// 目标URL
	postUrl := fmt.Sprintf("https://plats.rpgamestest.com/plat/%s/GetGameList", s)

	// 发起GET请求
	resp, err := http.Get(postUrl)
	if err != nil {
		// 处理错误
		fmt.Println("Error:", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(resp.Body) // 确保在函数返回时关闭响应体

	// 读取响应体
	body, err := io.ReadAll(resp.Body)

	var gameList GameListRequest
	err = json.Unmarshal(body, &gameList)

	oldSubstring := "dl"
	newSubstring := "ftp"

	if s == "PG" {
		for i := 0; i < len(gameList.Data["List"]); i++ {
			substring := strings.Replace(gameList.Data["List"][i].Icon, oldSubstring, newSubstring, 1)
			gameList.Data["List"][i].Icon = substring
		}
	}

	if err != nil {
		// 处理读取响应体时的错误
		fmt.Println("Error reading response body:", err)
	}
	body, err = json.Marshal(gameList)

	return string(body)
}

// 拉取PP的游戏列表
func handlePPGameList(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	var getGameListRequest GetGameListRequest
	err = json.Unmarshal(body, &getGameListRequest)
	if err != nil {
		return
	}

	ppGameList := GetGameList(getGameListRequest.Type)

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(ppGameList)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}
	// 写入响应体
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

type CreatePP struct {
	UID string `json:"UID"`
}

// 创建用户账号
func ppCreatePlayer(uid string, t string) {
	var createPP CreatePP
	createPP.UID = uid
	jsonData, err := json.Marshal(createPP)
	jsonReader := bytes.NewReader(jsonData)

	URL := fmt.Sprintf("https://plats.rpgamestest.com/plat/%s/Regist", t)

	//发送post请求
	req, err := http.NewRequest("POST", URL, jsonReader)
	if err != nil {
		fmt.Println(err)
		return
	}
	//设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "admin")
	req.Header.Set("Pid", "100001")
	req.Header.Set("Appid", "faketrans")
	req.Header.Set("Appsecret", "b6337af9-a91a-4085-b1f2-466923470735")
	//发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resp.Body)
	// 读取响应体
	var body []byte
	body, err = io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// ////////
type ppGoldInStruct struct {
	UID    string  `json:"UID"`
	Amount float64 `json:"Amount"`
}

// pp金币带入
func goldIn(uid string, amount float64, t string) {
	var ppGoldInStruct ppGoldInStruct
	ppGoldInStruct.UID = uid
	ppGoldInStruct.Amount = amount
	jsonData, err := json.Marshal(ppGoldInStruct)
	jsonReader := bytes.NewReader(jsonData)

	URL := fmt.Sprintf("https://plats.rpgamestest.com/plat/%s/FundTransferIn", t)

	//发送post请求
	req, err := http.NewRequest("POST", URL, jsonReader)
	if err != nil {
		fmt.Println(err)
		return
	}
	//设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "admin")
	req.Header.Set("Pid", "100001")
	//发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		//更新用户金币
		updatePlayer(uid, 0)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resp.Body)

}

type PPGoldOutResp struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Data  struct {
		Amount float64 `json:"Amount"`
		Status string  `json:"Status"`
	} `json:"data"`
}

// pp金币带出
func goldOut(uid string, t string) PPGoldOutResp {
	var createPP CreatePP
	createPP.UID = uid
	jsonData, err := json.Marshal(createPP)
	jsonReader := bytes.NewReader(jsonData)

	URL := fmt.Sprintf("https://plats.rpgamestest.com/plat/%s/FundTransferOut", t)

	//发送post请求
	req, err := http.NewRequest("POST", URL, jsonReader)
	if err != nil {
		fmt.Println(err)
	}
	//设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "admin")
	req.Header.Set("Pid", "100001")
	//发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resp.Body)
	// 读取响应体
	var body []byte
	body, err = io.ReadAll(resp.Body)

	var ppGoldOutResp PPGoldOutResp
	err = json.Unmarshal(body, &ppGoldOutResp)

	if err != nil {
		fmt.Println(err)
	}
	return ppGoldOutResp

}

//

func CORSLaunchGame(w http.ResponseWriter, r *http.Request) {
	// 设置CORS响应头
	w.Header().Set("Access-Control-Allow-Origin", "*") // 或者 "*"
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// 调用下一个处理器
	launchGamePP(w, r)
}

func launchGamePP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("launchGamePP Error:", err)
		}
	}(r.Body)

	//请求转化成json
	var launchGameRequest LaunchGameRequest

	err = json.Unmarshal(body, &launchGameRequest)
	if err != nil {
		fmt.Println("launchGameRequest Error:", err)
		return
	}

	//查询用户剩余的金币
	player, err := findPlayerByUId(launchGameRequest.UID)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	//创建用户
	ppCreatePlayer(player.UId, launchGameRequest.Type)

	//带出金币
	//goldOut(player.UId, launchGameRequest.Type)

	var launchGameResp LaunchGameResp
	launchGameResp.Code = 0

	if launchGameRequest.Type == "PP" {
		game, err := LaunchGame(player.UId, launchGameRequest.Game, launchGameRequest.Lang, player.PlayerGold)
		if err != nil {
			return
		} else {
			//更新用户金币
			updatePlayer(player.UId, 0)
		}
		fmt.Println("game=====", game)
		launchGameResp.Data.Url = game

		jsonData, err := json.Marshal(launchGameResp)

		// 设置响应头
		w.Header().Set("Content-Type", "application/json")

		// 写入响应体
		_, err = w.Write(jsonData)
		if err != nil {
			return
		}
	} else if launchGameRequest.Type == "PG" {

		//带入金币
		goldIn(player.UId, player.PlayerGold, launchGameRequest.Type)

		URL := fmt.Sprintf("https://plats.rpgamestest.com/plat/%s/LaunchGame", launchGameRequest.Type)
		// 发送post请求
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
		if err != nil {
			fmt.Println(err)
			return
		}
		// 设置请求头
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", "admin")
		req.Header.Set("Pid", "100001")

		// 发起请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}(resp.Body)

		// 读取响应体
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 打印响应体
		fmt.Println("PG==============", string(body))

		// 设置响应头
		w.Header().Set("Content-Type", "application/json")

		// 写入响应体
		_, err = w.Write(body)
		if err != nil {
			return
		}
	}
}

func LaunchGame(username, game, lang string, gold float64) (url string, err error) {
	err = Regist(username)
	if err != nil {
		return
	}

	FundTransferIn(username, gold)

	var ret struct {
		GameURL string
	}
	err = invoke("/game/start/", map[string]string{
		"externalPlayerId": username,
		"gameId":           game,
		"language":         lang,
		"platform":         "MOBILE",
	}, &ret)
	url = ret.GameURL
	return
}

type callResponse struct {
	Error       string
	Description string
}

type PPConfig struct {
	ApiUrl      string
	IconUrl     string
	SecureLogin string
	SecretKey   string
}

var ppConfig = PPConfig{
	ApiUrl:      "https://api.prerelease-env.biz/IntegrationService/v3/http/CasinoGameAPI",
	IconUrl:     "https://happycasino.prerelease-env.biz/game_pic/square/200/%s.png",
	SecureLogin: "hllgd_hollygod",
	SecretKey:   "C1Bd51884d654e10",
}

func invoke(method string, args map[string]string, result interface{}) (err error) {

	cfg := ppConfig
	var data = url.Values{}
	data.Set("secureLogin", cfg.SecureLogin)
	//THB

	// args := platcomm.GetArgs(ps)
	for k, v := range args {
		data.Set(k, v)
	}

	content := data.Encode() + cfg.SecretKey

	data.Set("hash", ut.ToMD5Str(content))

	fields := slog.With(
		"plat", "pp",
		"url", cfg.ApiUrl+method,
		"reqpayload", data.Encode(),
	)
	defer func() {
		// fields["error"] = err
		// logrus.WithFields(fields).Info("invoke!!")

		fields.Info("invoke!!", "error", err)
	}()

	resp, err := http.PostForm(cfg.ApiUrl+method, data)
	if err != nil {
		return
	}
	fields = fields.With("httpstatus", resp.Status)

	defer resp.Body.Close()

	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		fields = fields.With("httpstatus", resp.Status)
		return
	}
	fields = fields.With("respbody", mux.TruncMsg(ret, 1024))

	var callRet callResponse
	err = json.Unmarshal(ret, &callRet)
	if err != nil {
		return err
	}
	if callRet.Error != "0" {
		errcode, _ := strconv.Atoi(callRet.Error)
		// return logic.CodeError{
		// 	Code:    errcode,
		// 	Message: callRet.Description,
		// }
		return define.NewErrCode(callRet.Description, errcode)
	}

	if result != nil {
		err = json.Unmarshal(ret, result)
	}
	return
}

func Regist(username string) (err error) {
	ps := map[string]string{
		"externalPlayerId": username,
		"currency":         "THB",
		// "currency": "IDR",
	}
	var result struct {
		PlayerID int
	}
	err = invoke("/player/account/create/", ps, &result)
	return
}

func FundTransferIn(username string, amount float64) (status string) {
	orderID := uuid.NewString()
	_, err := trans(username, orderID, amount)
	status = platcomm.GetTransStatus(err)
	return

}

func FundTransferOut(username string) (amount float64, status string) {
	//invoke("/game/session/terminate/", map[string]string{
	//	// "externalTransactionId": orderID,
	//	"externalPlayerId": username,
	//}, nil)

	amount, err := GetBalance(username)
	if err != nil {
		status = "ERROR:" + err.Error()
		return
	}

	orderID := uuid.NewString()
	_, err = trans(username, orderID, -amount)
	status = platcomm.GetTransStatus(err)
	return
}

func GetBalance(username string) (balance float64, err error) {
	var ret transferRet
	err = invoke("/balance/current/", map[string]string{
		"externalPlayerId": username,
	}, &ret)

	balance = ret.Balance
	return
}

type transferRet struct {
	Balance float64
}

func trans(username string, orderID string, amount float64) (balance float64, err error) {
	var ret transferRet
	err = invoke("/balance/transfer/", map[string]string{
		"externalPlayerId":      username,
		"externalTransactionId": orderID,
		"amount":                strconv.FormatFloat(amount, 'f', -1, 64),
	}, &ret)

	balance = ret.Balance
	return
}
