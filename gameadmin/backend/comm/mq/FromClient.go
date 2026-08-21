package mq

import (
	"encoding/json"
	"game/duck/logger"
	"game/duck/rpc1"
	"reflect"
	"strings"
)

type Forward struct {
	Pids []int64
	Data []byte
	L    string
}

type Base struct {
	Route string      `json:"route"`
	CBID  int         `json:"cbId"`
	Args  interface{} `json:"args"`
	Err   string      `json:"err,omitempty"`
}

type BareBase struct {
	Route string `json:"route"`
	CBID  int    `json:"cbId"`
}

type Request struct {
	cbid         int
	skipAutoSend bool
}

func (r *Request) TakeCBID() int {
	r.skipAutoSend = true
	return r.cbid
}

// hander 需要这样定义
// 目前限制第二个参数必须是 结构体指针。

// 返回和错误分开的写法
// func (*MSG) Enter(player *Player, req *SlotsPlayReq) (*Resp, error)

// 返回值可以是错误，也可以是其他。会自动识别
// func (*MSG) Enter(player *Player, req *SlotsPlayReq) any

// 第三个参数可用于接受CBID等其他数据。
// 并且设置 SkipAutoSend = true, 可以跳过自动发送结果
// func (*Player, userStruct, *Request?)

func MakeMsg(route, args any) []byte {
	mm := Base{}

	switch args.(type) {
	case error:
		mm.Err = rpc1.GetErrorMessage(args.(error))
	default:
		mm.Args = args
	}

	switch route.(type) {
	case int:
		mm.CBID = route.(int)
		if mm.CBID == 0 {
			logger.Err("cbid is 0")
		}
	case string:
		mm.Route = route.(string)
	default:
		logger.Errf("error route type [%T]\n", route)
	}

	msg, _ := json.Marshal(mm)
	return msg
}

// Unmarshal decode json args
func Unmarshal(message []byte, p interface{}) *Base {
	var base = &Base{
		Args: p,
	}

	err := json.Unmarshal(message, base)
	if err != nil {
		logger.Err("Unmarshal:", string(message))
		logger.Err(err)
	}

	return base
}

// Unmarshal decode json args
func UnmarshalBare(message []byte) *BareBase {
	var base = &BareBase{}

	err := json.Unmarshal(message, base)
	if err != nil {
		logger.Err("Unmarshal:", string(message))
		logger.Err(err)
	}

	return base
}

func chIsLower(ch byte) bool {
	return 'a' <= ch && ch <= 'z'
}

func upperFirstCharOfWord(word string) string {
	if len(word) == 0 {
		return word
	}

	if !chIsLower(word[0]) {
		return word
	}

	buf := []byte(word)
	buf[0] -= 'a' - 'A'

	return string(buf)
}

func PlayerMsgHandler[T any](subject string, hdr any, playerGetFn func(pid int64) T) {
	playerMsgHandlerP(subject, hdr, func(pid int64) any { return playerGetFn(pid) })
}

func playerMsgHandlerP(subject string, hdr any, playerGetFn func(pid int64) any) {
	typ := reflect.TypeOf(hdr)
	val := reflect.ValueOf(hdr)

	Subscribe(subject, func(fwd *Forward) {
		logger.Info(">>Recv", fwd.Pids, string(fwd.Data))
		base := UnmarshalBare(fwd.Data)

		arr := strings.Split(base.Route, ".")
		if len(arr) != 2 {
			logger.Info("bad method", base.Route)
			return
		}

		method := upperFirstCharOfWord(arr[1])
		mtype, ok := typ.MethodByName(method)
		if !ok {
			logger.Info("调用了不允许的函数", method)
			return
		}

		pid := fwd.Pids[0]
		plr := playerGetFn(pid)
		if plr == nil {
			logger.Info("获取玩家失败", pid)
			return
		}

		reqType := mtype.Type.In(2)
		req := reflect.New(reqType.Elem())

		Unmarshal(fwd.Data, req.Interface())

		reqInfo := &Request{cbid: base.CBID}

		vs := make([]reflect.Value, 2)
		vs[0] = reflect.ValueOf(plr)
		vs[1] = req

		// (0)(1,2,3) 接收器也占一个in
		if mtype.Type.NumIn() == 4 {
			vs = append(vs, reflect.ValueOf(reqInfo))
		}

		fn := val.MethodByName(method)
		ret := fn.Call(vs)

		if reqInfo.skipAutoSend || base.CBID == 0 {
			return
		}

		var send any
		if len(ret) == 1 {
			send = ret[0].Interface()

		} else {
			send = ret[0].Interface()
			if !ret[1].IsNil() {
				send = ret[1].Interface()
			}
		}

		fwd.Data = MakeMsg(base.CBID, send)

		forward(fwd)
		// Publish("gateway", fwd)
	})
}

func forward(f *Forward) {
	logger.Info("<<Send", f.Pids, string(f.Data))

	Publish("gateway", f)
}

func Send2ps(ps []int64, route any, data any) {
	if len(ps) == 0 {
		return
	}

	forward(&Forward{
		Pids: ps,
		Data: MakeMsg(route, data),
	})
}
