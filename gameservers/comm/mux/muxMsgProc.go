package mux

/*
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"

	"serve/comm/db"
	"serve/comm/mq"
	"github.com/nats-io/nats.go"
)

const (
	MsgProcParamUID = iota
	MsgProcParamArgs
	MsgProcParamReply
	MsgProcParamNum
)

type RawBase struct {
	Svr   string          `json:"svr"`
	Route string          `json:"route"`
	CBID  int             `json:"cbId"`
	Args  json.RawMessage `json:"args"`
}

// type WSReqContext struct {
// 	Pid      int64
// 	Language string
// }

// var (
// 	WSReqContextType = reflect.TypeOf(&WSReqContext{})
// )

type WSReqContext = db.WSReqContext

var WSReqContextType = db.WSReqContextType

// func (m *Mux) OnRecv(pids []int64, message []byte) (cbid int, pid int64, result any) {

// func (m *Mux) OnRecv(pids []int, message []byte) {
// 	go m.onRecv(pids, message)
// }

func validMsgPorc(handler interface{}) {
	t := reflect.TypeOf(handler)

	Assert(t.NumIn() == MsgProcParamNum)

	pstypeReply := t.In(MsgProcParamReply)
	Assert(pstypeReply.Kind() == reflect.Ptr)

	// Assert(t.In(MsgProcParamUID).Kind() == reflect.Int64)
	Assert(t.In(MsgProcParamUID) == WSReqContextType)

	Assert(t.NumOut() == 1)
	u := reflect.TypeOf((*error)(nil)).Elem()
	Assert(t.Out(0).Implements(u))
}

func MsgProcGetArg0(r *http.Request) (ans any, err error) {
	// Authorization: pid:12345
	// token := r.Header.Get("Authorization")
	// if !strings.HasPrefix(token, "pid:") {
	// 	err = fmt.Errorf("Authorization[%s] format Error!", token)
	// 	return
	// }
	// username := strings.TrimPrefix(token, "pid:")
	username := r.Header.Get("pid")
	language := r.Header.Get("language")
	if language == "" {
		language = "th"
	}

	pid, err := strconv.Atoi(username)
	if err != nil {
		return
	}

	ctx := &WSReqContext{
		Pid:      int64(pid),
		Language: language,
	}

	ans = ctx
	return
}

func (m *Mux) RegMsgProc(path string, desc string, kind string, handler interface{}, ps interface{}) *PHandler {
	validMsgPorc(handler)

	data := &PHandler{
		Path:         path,
		Handler:      handler,
		Desc:         desc,
		Kind:         kind,
		ParamsSample: ps,
		Class:        ClassMsgProc,
		GetArg0:      MsgProcGetArg0,
	}

	return m.Add(data)
}
func RegMsgProc(path string, desc string, kind string, handler interface{}, ps interface{}) *PHandler {
	return DefaultRpcMux.RegMsgProc(path, desc, kind, handler, ps)
}

func (m *Mux) SubscribeWSMsg(subject string, gobNC *nats.EncodedConn) {
	_, err := gobNC.Subscribe(subject, func(fwd *mq.Forward) {
		// cbid, pid, result := m.OnRecv(fwd.Pids, fwd.Data)
		cbid, pid, result := m.OnRecv(fwd)
		if cbid != 0 {
			mq.Send2ps([]int64{pid}, cbid, result)
		}
	})

	if err != nil {
		log.Panicf("SubscribeWSMsg Fail!, subject=%v, error=%v", subject, err.Error())
	}
}

func SubscribeWSMsg(subject string, gobNC *nats.EncodedConn) {
	DefaultRpcMux.SubscribeWSMsg(subject, gobNC)
}

func (m *Mux) OnRecv(fwd *mq.Forward) (cbid int, pid int64, result any) {
	pids := fwd.Pids
	message := fwd.Data

	logUnit := NewReqLogUnit()

	defer func() {
		logUnit.Print()
	}()

	var top RawBase
	err := json.Unmarshal(message, &top)
	if err != nil {
		// logger.Err(err)
		logUnit.Err = err.Error()
		return
	}

	pid = pids[0]
	logUnit.UID = strconv.Itoa(int(pid))
	logUnit.URI = top.Route
	cbid = top.CBID
	// var result interface{}
	// if top.CBID != 0 {
	// 	defer func() {
	// 		comm.DefaultMQDailer.Send2p(pid, top.CBID, result)
	// 	}()
	// }

	hand := m.handlers[top.Route]
	if hand == nil || hand.Class != "msgproc" {
		result = fmt.Errorf("handler[%s] not found", top.Route)
		logUnit.Err = fmt.Sprintf("handler[%s] not found", top.Route)
		return
	}

	handler := hand.Handler

	t := reflect.TypeOf(handler)
	paramsType := t.In(MsgProcParamArgs)
	paramsValue, err := bindwrap(paramsType, func(i interface{}) error {
		return json.Unmarshal(top.Args, i)
	})

	if err != nil {
		result = err
		logUnit.Err = err.Error()
		return
	}
	logUnit.Params = paramsValue.Interface()

	ctx := &WSReqContext{
		Pid:      pid,
		Language: fwd.L,
	}
	var in [MsgProcParamNum]reflect.Value
	in[MsgProcParamUID] = reflect.ValueOf(ctx)
	in[MsgProcParamArgs] = paramsValue

	pstypeReply := t.In(MsgProcParamReply)
	paramsReply := reflect.New(pstypeReply.Elem())
	in[MsgProcParamReply] = paramsReply

	fn := reflect.ValueOf(handler)

	defer func() {
		if x := recover(); x != nil {
			os.Stdout.Write(debug.Stack())
			log.Println(x)

			result = fmt.Errorf("panic %v", x)
			logUnit.Err = fmt.Sprintf("panic %v", x)
		}
	}()

	out := fn.Call(in[:])

	if !out[0].IsNil() {
		err = out[0].Interface().(error)
		result = err
		logUnit.Err = err.Error()
		return
	}
	reply := paramsReply.Interface()
	result = reply

	logUnit.Result = result

	return
}
*/
