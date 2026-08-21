package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	RegMsgProc("/AdminInfo/ZIPBetLogAdd", "下注历史导出任务(添加)", "AdminInfo", addBetHistoryDownload, BetHistoryDownload{
		PageSize: 100000,
	})
	RegMsgProc("/AdminInfo/ZIPBetLogList", "下注历史导出任务(列表)", "AdminInfo", betHistoryDownloadList, listParam{})
	ensuredownloadPath()
}

const (
	waitingProcess processStatus = iota // 已生成下载任务
	inProcess                           // 下载任务处理中
	finished                            //下载文件已生成
)

var (
	downloadColl      *mongodb.Collection     = DB.Collection("BetHistoryDownload")
	downloadChan      chan primitive.ObjectID = make(chan primitive.ObjectID)
	BHdownloadPath    string                  = "./BHdownload"
	BHdownloadFilePre string                  = "下注历史"
	CSVSuf            string                  = ".csv"
	//ZIPSuf            string                  = ".zip"
)

type processStatus int

type BetHistoryDownload struct {
	ID            primitive.ObjectID `json:"ID" bson:"_id"`
	Pid           int64              `json:"Pid" bson:"Pid"`
	Username      string             `json:"Username" bson:"Username"`
	Lang          string             `json:"Lang" bson:"Lang"`
	OperatorId    int64              `json:"OperatorId" bson:"OperatorId"`
	OrderId       string             `json:"OrderId" bson:"OrderId"`
	UserID        string             `json:"UserID" bson:"UserID"`
	StartTime     int64              `json:"StartTime" bson:"StartTime"`
	EndTime       int64              `json:"EndTime" bson:"EndTime"`
	GameID        string             `json:"GameID" bson:"GameID"`
	Bet           int64              `json:"Bet" bson:"Bet"`
	Win           int64              `json:"Win" bson:"Win"`
	WinLose       int64              `json:"WinLose" bson:"WinLose"`
	Balance       int64              `json:"Balance" bson:"Balance"`
	PageSize      int64              `json:"-" bson:"PageSize"`
	FileName      string             `json:"FileName" bson:"FileName"`
	CreateAt      int64              `json:"CreateAt" bson:"CreateAt"`
	Creater       string             `json:"Creater" bson:"Creater"`
	ProcessStatus processStatus      `json:"ProcessStatus" bson:"ProcessStatus"`
}

type pageSetting struct {
	PageSize   int
	PageNumber int
	Count      int
}

type listParam struct {
	OperatorId int64
	StartTime  int64
	EndTime    int64
	pageSetting
}

type listDownload struct {
	List []BetHistoryDownload
	pageSetting
}

func addBetHistoryDownload(ctx *Context, ps getBetLogListNewParams, ret *comm.Empty) (err error) {
	if ps.UserID == "" {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	if ps.StartTime != 0 && ps.EndTime != 0 {
		if ps.EndTime-ps.StartTime > 7*24*60*60 {
			return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
		}
	} else if ps.StartTime == 0 && ps.EndTime != 0 {
		e := time.Unix(ps.EndTime, 0)
		ps.StartTime = e.AddDate(0, 0, -7).Unix()
	} else if ps.StartTime != 0 && ps.EndTime == 0 {
		e := time.Unix(ps.StartTime, 0)
		ps.EndTime = e.AddDate(0, 0, 7).Unix()
	} else {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	ps.PageSize = 100 * 1000
	appID, err := getAppID(ctx)
	if err != nil {
		return
	}
	tmp_id := primitive.NewObjectID()
	var bet, win, winLose, balance int64
	if ps.Bet != nil {
		bet = *ps.Bet
	}
	if ps.Win != nil {
		win = *ps.Win
	}
	if ps.WinLose != nil {
		winLose = *ps.WinLose
	}
	if ps.Balance != nil {
		balance = *ps.Balance
	}
	d := BetHistoryDownload{
		ID:            tmp_id,
		Pid:           ctx.PID,
		Username:      ctx.Username,
		Lang:          ctx.Lang,
		OperatorId:    ps.OperatorId,
		OrderId:       ps.OrderId,
		UserID:        ps.UserID,
		StartTime:     ps.StartTime,
		EndTime:       ps.EndTime,
		GameID:        ps.GameID,
		Bet:           bet,
		Win:           win,
		WinLose:       winLose,
		Balance:       balance,
		PageSize:      100000,
		CreateAt:      time.Now().Unix(),
		Creater:       appID,
		ProcessStatus: waitingProcess,
	}
	err = downloadColl.InsertOne(d)
	if err != nil {
		return err
	}
	go func() {
		downloadChan <- tmp_id
	}()
	return
}

func betHistoryDownloadList(ctx *Context, ps listParam, ret *listDownload) (err error) {
	// 取消7天时间限制
	// if ps.StartTime != 0 && ps.EndTime != 0 {
	// 	if ps.EndTime-ps.StartTime > 7*24*60*60 {
	// 		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	// 	}
	// } else if ps.StartTime == 0 && ps.EndTime != 0 {
	// 	e := time.Unix(ps.EndTime, 0)
	// 	ps.StartTime = e.AddDate(0, 0, -7).Unix()
	// 	if ps.StartTime <= 0 {
	// 		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	// 	}
	// } else if ps.StartTime != 0 && ps.EndTime == 0 {
	// 	e := time.Unix(ps.StartTime, 0)
	// 	ps.EndTime = e.AddDate(0, 0, 7).Unix()
	// 	if ps.EndTime <= 0 {
	// 		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	// 	}
	// } else {
	// 	return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	// }
	if ps.PageNumber == 0 {
		ps.PageNumber = 1
	}
	if ps.PageSize == 0 {
		ps.PageSize = 20
	}
	filter := bson.M{"OperatorId": ps.OperatorId}
	if ps.StartTime != 0 {
		filter["StartTime"] = bson.M{"$gte": ps.StartTime}
	}
	if ps.EndTime != 0 {
		filter["EndTime"] = bson.M{"$lte": ps.EndTime}
	}
	findOptions := options.Find()
	skip := (ps.PageNumber - 1) * ps.PageSize
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"CreateAt": -1})
	findOptions.SetLimit(int64(ps.PageSize))

	err = downloadColl.FindAllOpt(filter, &ret.List, findOptions)
	if err != nil {
		return err
	}
	for i := range ret.List {
		dback := &ret.List[i]
		dback.FileName = fmt.Sprintf("%s%s%s", lang.GetLang(ctx.Lang, BHdownloadFilePre), dback.ID.Hex(), CSVSuf)
	}
	ret.PageNumber = ps.PageNumber
	ret.PageSize = ps.PageSize
	c, _ := downloadColl.CountDocuments(filter)
	ret.Count = int(c)
	return
}

func removeBetLogDownload() {
	slog.Log(context.Background(), slog.LevelInfo, " >>>>>> removeBetLogDownload 9.30 remove old records start")
	delTime := time.Now().AddDate(0, 0, -7).Unix()
	filter := bson.M{}
	filter["CreateAt"] = bson.M{"$lte": delTime}
	ret := listDownload{}
	err := downloadColl.FindAll(filter, &ret.List)
	if err == nil {
		for _, v := range ret.List {
			fname := fmt.Sprintf("%s/%s", BHdownloadPath, v.FileName)
			if _, err := os.Stat(fname); os.IsNotExist(err) {
				fmt.Printf("%s not exist!", fname)
				continue
			}
			err = os.Remove(fname)
			if err != nil {
				fmt.Printf("remove file occurrd error :%s", err.Error())
				continue
			}
			downloadColl.DeleteId(v.ID)
		}
	}
	slog.Log(context.Background(), slog.LevelInfo, " <<<<<< removeBetLogDownload 9.30 remove old records end")
}

func doDownloadBetHistoryExcle(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	cker := &checker{}
	pid, _ := cker.GetPidAndUname(token)
	fmt.Printf("download login pid :%d", pid)
	if pid == 0 {
		return
	}
	type pID struct {
		ID string
	}
	param4ID := pID{}
	err := json.NewDecoder(r.Body).Decode(&param4ID)
	fmt.Printf("download login pid :%s", param4ID.ID)
	if err != nil {
		return
	}
	fileName := fmt.Sprintf("%s%s%s", BHdownloadFilePre, param4ID.ID, CSVSuf)
	filePath := fmt.Sprintf("%s/%s", BHdownloadPath, fileName)
	fmt.Printf("download full path :%s", filePath)
	http.ServeFile(w, r, filePath)

	// 二进制流
	// file, err := os.Open(filePath)
	// if err != nil {
	// 	return
	// }
	// defer file.Close()
	// reader := bufio.NewReader(file)
	// buf := make([]byte, 1024)
	// for {
	// 	n, err := reader.Read(buf)
	// 	if err != nil && err != io.EOF {
	// 		return
	// 	}
	// 	if n == 0 {
	// 		break
	// 	}
	// 	_, err = w.Write(buf[:n])
	// 	if err != nil {
	// 		return
	// 	}
	// }
	// w.Header().Set("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	// w.Header().Set("Content-Transfer-Encoding", "binary")
}

func downloadInProcess(ID primitive.ObjectID) (err error) {
	update := bson.M{
		"$set": bson.M{
			"ProcessStatus": inProcess,
		},
	}
	downloadColl.UpdateId(ID, update)
	return
}

func downloadFinished(ID primitive.ObjectID, FileName string) (err error) {
	update := bson.M{
		"$set": bson.M{
			"ProcessStatus": finished,
			"FileName":      FileName,
		},
	}
	downloadColl.UpdateId(ID, update)
	return
}

func DownBetHistory() {

	go prepareBetHistoryDown()

	go doDownBetHistory()
}

func prepareBetHistoryDown() {
	for {
		filter := bson.M{"ProcessStatus": waitingProcess}
		d := BetHistoryDownload{}
		err := downloadColl.FindOne(filter, &d)
		if err == nil {
			downloadChan <- d.ID
		} else {
			break
		}
	}
}

func doDownBetHistory() {
	for {
		id := <-downloadChan
		for {
			d := BetHistoryDownload{}
			err := downloadColl.FindId(id, &d)
			if err == nil {
				if d.ProcessStatus == waitingProcess {
					downloadInProcess(id)
					c := &Context{
						PID:      d.Pid,
						Username: d.Username,
						Lang:     d.Lang,
					}
					p := getBetLogListNewParams{
						OperatorId: d.OperatorId,
						OrderId:    d.OrderId,
						// Pid:        d.Pid,
						UserID:    d.UserID,
						StartTime: d.StartTime,
						EndTime:   d.EndTime,
						GameID:    d.GameID,
						Bet:       &(d.Bet),
						Win:       &(d.Win),
						WinLose:   &(d.WinLose),
						Balance:   &(d.Balance),
						PageSize:  d.PageSize,
					}
					downloadBetLogListNew(c, p, id.Hex())
					fileName := fmt.Sprintf("%s%s%s", BHdownloadFilePre, id.Hex(), CSVSuf)
					downloadFinished(id, fileName)
				}
				break
			} else {
				time.Sleep(time.Second)
				continue
			}
		}
	}
}

func ensuredownloadPath() {
	_, err := os.Stat(BHdownloadPath)
	if os.IsNotExist(err) {
		os.MkdirAll(BHdownloadPath, 0755)
	}
}

func getAppID(ctx *Context) (string, error) {
	user := comm.User{}
	err := CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return "", err
	}
	appID := user.AppID
	return appID, nil
}
