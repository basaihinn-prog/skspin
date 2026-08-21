package main

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
)

type __msg_send_ struct {
	t   int
	buf []byte
}

type IMsgHandler interface {
	OnRecvFrom(conn *Connection, message []byte)
	OnClose(conn *Connection)
}

type Connection struct {
	ws       *websocket.Conn
	sendChan chan __msg_send_

	waitD time.Duration

	done chan struct{}

	handler          IMsgHandler
	ip               string
	gameid, language string
	pid              int64
}

func NewConnection(ws *websocket.Conn, ip string, handler IMsgHandler) *Connection {
	conn := new(Connection)
	conn.sendChan = make(chan __msg_send_, 1024)
	conn.ws = ws
	conn.waitD = time.Second //time.Second
	// conn.waitD = time.Minute //time.Second

	conn.done = make(chan struct{})
	conn.handler = handler
	conn.ip = ip

	return conn
}

var pong = []byte("pong")
var ping = []byte("ping")

func (this *Connection) readPump() {
	ws := this.ws

	defer func() {
		ws.Close()
		close(this.done)
	}()

	for {
		// 7秒没数据读入就会 触发 timeout错误
		ws.SetReadDeadline(time.Now().Add(this.waitD * 7))
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			if !bytes.Equal(message, pong) {
				this.handler.OnRecvFrom(this, message)
			}
		}
	}
}

func (this *Connection) writePump() {
	ticker := time.NewTicker(2 * time.Second)
	ws := this.ws
	sendChan := this.sendChan

	defer func() {
		ticker.Stop()
		ws.Close()
		// Info("write pump exit")
		this.handler.OnClose(this)
	}()

	for {
		select {
		case <-this.done:
			return
		case message, ok := <-sendChan:
			if !ok {
				return
			}

			ws.SetWriteDeadline(time.Now().Add(this.waitD * 7))
			err := ws.WriteMessage(message.t, message.buf)
			if err != nil {
				return
			}
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(this.waitD * 7))
			err := ws.WriteMessage(websocket.TextMessage, ping)
			if err != nil {
				return
			}
		}
	}
}

func (this *Connection) Close() {
	this.ws.Close()
}

func (this *Connection) ServeWs() {
	go this.writePump()
	this.readPump()
}

func (this *Connection) SendMsg(msg []byte) {
	select {
	case this.sendChan <- __msg_send_{
		t:   websocket.TextMessage,
		buf: msg,
	}:
	default:
	}
}

func (this *Connection) IP() string {
	return this.ip
}
