package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"game/comm"
	"game/duck/lang"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/saveOperator4GameUpdateMsg", "配置游戏更新消息的运营商", "AdminInfo", saveOperator4GameUpdateMsg, paramTelegramDest{
		DestOperatorID: []int64{},
	})

	RegMsgProc("/AdminInfo/saveMsg4GameUpdateMsg", "配置游戏更新的消息内容", "AdminInfo", saveMsg4GameUpdateMsg, paramTelegramMsg{
		Msg: "",
	})

	RegMsgProc("/AdminInfo/getConfig4GameUpdateMsg", "获取游戏更新消息的配置", "AdminInfo", getConfig4GameUpdateMsg, comm.Empty{})

	RegMsgProc("/AdminInfo/sendMsg4GameUpdate", "发送游戏更新消息", "AdminInfo", sendMsg4GameUpdate, paramGameUpdate4Telegram{
		paramTelegramMsg:  paramTelegramMsg{Msg: ""},
		paramTelegramDest: paramTelegramDest{DestOperatorID: []int64{}},
	})

	RegMsgProc("/AdminInfo/Result4GameUpdateMsg", "获取游戏更新消息发送结果", "AdminInfo", Result4GameUpdateMsg, comm.Empty{})

	RegMsgProc("/AdminInfo/Status4GameUpdateMsg", "获取游戏更新消息是否发送完成", "AdminInfo", Status4GameUpdateMsg, comm.Empty{})
}

const (
	GameUpdate4Telegram_Msg      GameUpdate4TelegramID = iota // 消息
	GameUpdate4Telegram_Operator                              // 发送的渠道
	GameUpdate4Telegram_Result                                // 发送结果
)

const telegram_url = "https://api.telegram.org/bot%s/sendMessage"

type sendTelegramStatus struct {
	Status int64
}

var sendingStatus = sendTelegramStatus{Status: 0}

type TelegramBackStatus struct {
	Ok bool `json:"ok"`
}
type GameUpdate4TelegramID int
type GameUpdateMsg4Telegram struct {
	ID                GameUpdate4TelegramID `json:"ID" bson:"ID"`
	Completed         *bool                 `json:"Completed" bson:"Completed"`
	Msg               string                `json:"Msg" bson:"Msg"`
	DestOperatorID    []int64               `json:"DestOperatorID" bson:"DestOperatorID"`       //目标值
	FailOperatorID    []int64               `json:"FailOperatorID" bson:"FailOperatorID"`       //失败的
	SuccessOperatorID []int64               `json:"SuccessOperatorID" bson:"SuccessOperatorID"` //成功的
}
type paramGameUpdate4Telegram struct {
	paramTelegramMsg
	paramTelegramDest
}

type paramTelegramMsg struct {
	Msg string
}

type paramTelegramDest struct {
	DestOperatorID []int64
}

func saveMsg4GameUpdateMsg(ctx *Context, ps paramTelegramMsg, ret *comm.Empty) (err error) {
	if len(ps.Msg) == 0 {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	col := DB.Collection("GameUpdateMsg4Telegram")
	err = col.UpsertOne(bson.M{"ID": GameUpdate4Telegram_Msg},
		bson.M{"$set": bson.M{
			"Msg": ps.Msg,
		}})
	if err != nil {
		return
	}
	return nil
}

func saveOperator4GameUpdateMsg(ctx *Context, ps paramTelegramDest, ret *comm.Empty) (err error) {
	if len(ps.DestOperatorID) == 0 {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	col := DB.Collection("GameUpdateMsg4Telegram")
	err = col.UpsertOne(bson.M{"ID": GameUpdate4Telegram_Operator},
		bson.M{"$set": bson.M{
			"DestOperatorID": ps.DestOperatorID,
		}})
	if err != nil {
		return
	}
	return nil
}

func getConfig4GameUpdateMsg(ctx *Context, ps comm.Empty, ret *paramGameUpdate4Telegram) (err error) {
	tmp := &GameUpdateMsg4Telegram{}
	col := DB.Collection("GameUpdateMsg4Telegram")
	err = col.FindOne(bson.M{"ID": GameUpdate4Telegram_Operator}, tmp)
	if err != nil {
		return
	}
	ret.DestOperatorID = tmp.DestOperatorID
	err = col.FindOne(bson.M{"ID": GameUpdate4Telegram_Msg}, tmp)
	if err != nil {
		return
	}
	ret.Msg = tmp.Msg
	return nil
}

func Status4GameUpdateMsg(ctx *Context, ps comm.Empty, ret *sendTelegramStatus) (err error) {
	*ret = sendingStatus
	if ret.Status == 0 {
		col := DB.Collection("GameUpdateMsg4Telegram")
		tmp := &GameUpdateMsg4Telegram{}
		err = col.FindOne(bson.M{"ID": GameUpdate4Telegram_Result}, tmp)
		if err == nil && tmp.Completed != nil && !*tmp.Completed { // 发送消息中，异常中断 completed == nil
			col.UpsertOne(bson.M{"ID": GameUpdate4Telegram_Result},
				bson.M{"$set": bson.M{
					"Completed": nil,
				}})
		}
	}
	return nil
}

func Result4GameUpdateMsg(ctx *Context, ps comm.Empty, ret *GameUpdateMsg4Telegram) (err error) {
	col := DB.Collection("GameUpdateMsg4Telegram")
	col.FindOne(bson.M{"ID": GameUpdate4Telegram_Result}, ret)
	return nil
}

func sendMsg4GameUpdate(ctx *Context, ps paramGameUpdate4Telegram, ret *comm.Empty) (err error) {
	notSending := atomic.CompareAndSwapInt64(&sendingStatus.Status, 0, 1)
	if !notSending {
		return errors.New(lang.GetLang(ctx.Lang, "消息发送中"))
	}
	if len(ps.Msg) == 0 {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	col := DB.Collection("GameUpdateMsg4Telegram")
	err = col.UpsertOne(bson.M{"ID": GameUpdate4Telegram_Msg},
		bson.M{"$set": bson.M{
			"Msg": ps.Msg,
		}})
	if err != nil {
		return
	}
	err = col.UpsertOne(bson.M{"ID": GameUpdate4Telegram_Result},
		bson.M{"$set": bson.M{
			"Completed":         false,
			"DestOperatorID":    ps.DestOperatorID,
			"FailOperatorID":    []int64{},
			"SuccessOperatorID": []int64{},
		}})
	if err != nil {
		return
	}
	go doSend(ps)
	return nil
}

func doSend(ps paramGameUpdate4Telegram) {
	var sendGroup sync.WaitGroup
	filter, update := bson.M{"ID": GameUpdate4Telegram_Result}, bson.M{}
	col := DB.Collection("GameUpdateMsg4Telegram")
	for _, v := range ps.DestOperatorID {
		sendGroup.Add(1)
		go func() {
			defer sendGroup.Done()
			o := &comm.Operator{}
			CollAdminOperator.FindId(v, o)
			if len(o.Robot) == 0 || len(o.ChatID) == 0 {
				update = bson.M{"$addToSet": bson.M{"FailOperatorID": v}}
			} else {
				baseURL := fmt.Sprintf(telegram_url, o.Robot)
				params := url.Values{}
				params.Add("chat_id", o.ChatID)
				params.Add("text", ps.Msg)
				fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
				r, e := http.Get(fullURL)
				if e != nil {
					update = bson.M{"$addToSet": bson.M{"FailOperatorID": v}}
					fmt.Println(e.Error())
				} else {
					// 读取响应体
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					defer r.Body.Close()
					responseString := string(body)
					var status TelegramBackStatus
					json.Unmarshal([]byte(responseString), &status)
					if status.Ok {
						update = bson.M{"$addToSet": bson.M{"SuccessOperatorID": v}}
					} else {
						update = bson.M{"$addToSet": bson.M{"FailOperatorID": v}}
					}
				}
			}
			err := col.UpdateOne(filter, update)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	sendGroup.Wait()
	col.UpsertOne(bson.M{"ID": GameUpdate4Telegram_Result},
		bson.M{"$set": bson.M{
			"Completed": true,
		}})
	sendingStatus.Status = 0
}
