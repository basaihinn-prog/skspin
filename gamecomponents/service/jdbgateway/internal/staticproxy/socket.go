package staticproxy

import (
	"encoding/json"
	"fmt"
	"game/comm/mq"
	"game/comm/slotsmongo"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/nats-io/nats.go"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	GEN_HEARTBEAT = "GEN_HEARTBEAT"
	GameLogin     = "gameLogin"
	H5Init        = "h5.init"
	H5Spin        = "h5.spin"
)

// ConnInfo 存储单个 WebSocket 连接的信息
type ConnInfo struct {
	Conn        *websocket.Conn // WebSocket 连接
	Pid         int64           // 玩家 ID
	MachineType int32           // 游戏id
	CreatedAt   time.Time       // 连接创建时间
}

// ConnectionManager WebSocket 连接管理器
type ConnectionManager struct {
	conns  map[*websocket.Conn]*ConnInfo // 连接映射
	mutex  sync.RWMutex                  // 读写锁，确保线程安全
	pidMap map[int64]*websocket.Conn     // pid 到连接的映射，用于快速查找
}

// NewConnectionManager 创建一个新的连接管理器
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns:  make(map[*websocket.Conn]*ConnInfo),
		pidMap: make(map[int64]*websocket.Conn),
	}
}

// AddConn 添加一个新的 WebSocket 连接
func (cm *ConnectionManager) AddConn(conn *websocket.Conn, pid int64) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	connInfo := &ConnInfo{
		Conn:      conn,
		Pid:       pid,
		CreatedAt: time.Now(),
	}
	cm.conns[conn] = connInfo
	if pid > 0 { // 如果有有效的 pid，则添加到 pidMap
		cm.pidMap[pid] = conn
	}
	slog.Info("Added new connection", "pid", pid, "remoteAddr", conn.RemoteAddr().String())
}

// StartHeartbeat 开始心跳检测，定时向所有连接发送心跳包
func (cm *ConnectionManager) StartHeartbeat() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		<-ticker.C
		cm.mutex.RLock()
		for _, connInfo := range cm.conns {
			so := slotsmongo.NewSFSObject()
			p := slotsmongo.NewSFSObject()
			p.PutInt("r", 235)
			p.PutShort("uc", int16(len(cm.conns)))
			so.PutSFSObject("p", p)
			so.PutShort("a", 1001)
			so.PutByte("c", 0)
			bytes, err := so.ToBinary()
			if err != nil {
				slog.Error("ToBinary error", "err", err)
				return
			}
			if err := connInfo.Conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
				slog.Error("Failed to send heartbeat", "pid", connInfo.Pid, "error", err)
			}
		}
		cm.mutex.RUnlock()
	}
}
func (cm *ConnectionManager) UpdateConnMType(conn *websocket.Conn, pid int64, mType int32) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	connInfo := &ConnInfo{
		Conn:        conn,
		Pid:         pid,
		CreatedAt:   time.Now(),
		MachineType: mType,
	}
	cm.conns[conn] = connInfo
	if pid > 0 { // 如果有有效的 pid，则添加到 pidMap
		cm.pidMap[pid] = conn
	}
	slog.Info("Added new connection", "pid", pid, "remoteAddr", conn.RemoteAddr().String())
}

// RemoveConn 移除一个 WebSocket 连接
func (cm *ConnectionManager) RemoveConn(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if connInfo, exists := cm.conns[conn]; exists {
		delete(cm.conns, conn)
		if connInfo.Pid > 0 { // 如果有 pid，则从 pidMap 中移除
			delete(cm.pidMap, connInfo.Pid)
		}
		slog.Info("Removed connection", "pid", connInfo.Pid, "remoteAddr", conn.RemoteAddr().String())
	}
}

// GetConnByPid 根据 pid 获取对应的 WebSocket 连接
func (cm *ConnectionManager) GetConnByPid(pid int64) (*websocket.Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	conn, exists := cm.pidMap[pid]
	return conn, exists
}

// GetConnInfo 获取某个连接的详细信息
func (cm *ConnectionManager) GetConnInfo(conn *websocket.Conn) (*ConnInfo, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connInfo, exists := cm.conns[conn]
	return connInfo, exists
}

// Broadcast 向所有连接广播消息
func (cm *ConnectionManager) Broadcast(messageType websocket.MessageType, data []byte) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	for conn := range cm.conns {
		conn.WriteMessage(messageType, data)
	}
}

// SendToPid 向特定 pid 的连接发送消息
func (cm *ConnectionManager) SendToPid(pid int64, messageType websocket.MessageType, data []byte) error {
	conn, exists := cm.GetConnByPid(pid)
	if !exists {
		return fmt.Errorf("no connection found for pid: %d", pid)
	}
	return conn.WriteMessage(messageType, data)
}

// GlobalConnManager 全局连接管理器实例
var GlobalConnManager = NewConnectionManager()

func OnWebsocket(w http.ResponseWriter, r *http.Request) {
	slog.Info("成功收到wss")
	ws := newUpgrader()
	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	conn.SetSession(r.URL.Query())
	log.Println("OnOpen:", conn.RemoteAddr().String())
}

// newUpgrader 创建 WebSocket Upgrader
func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	u.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		defer func() {
			c.SetReadDeadline(time.Now().Add(nbhttp.DefaultKeepaliveTime))
		}()

		reqSo, err := slotsmongo.NewFromBinaryData(data)
		if err != nil {
			slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
			return
		}

		// 心跳
		p, _ := reqSo.GetSFSObject("p")
		pc, _ := p.GetString("c")
		if pc == GEN_HEARTBEAT {
			return
		}

		action, ok := reqSo.GetShort("a")
		if ok {
			// version + javascript
			if action == 0 {
				so := slotsmongo.C0A0Pctmstk()
				bytes, err := so.ToBinary()
				if err != nil {
					slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
				}
				c.WriteMessage(messageType, bytes)
			}
			if action == 1 {
				userName, _ := p.GetString("un")
				splits := strings.Split(userName, "@")

				pid, err := strconv.Atoi(splits[0])
				if err != nil {
					slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
					return
				}
				so := slotsmongo.C0A1P(pid, userName)
				bytes, err := so.ToBinary()
				if err != nil {
					slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
					return
				}
				c.WriteMessage(messageType, bytes)
				// 在登录成功后，将连接添加到管理器中
				// 假设登录成功后可以从返回数据中获取 pid，这里使用硬编码的示例值
				GlobalConnManager.AddConn(c, int64(pid))
			}
			if action == 13 {
				if pc == GameLogin {
					pp, _ := p.GetSFSObject("p")
					info, ok := GlobalConnManager.GetConnInfo(c)
					if !ok {
						slog.Error("GameLogin Err", "globalConnManager.GetConnInfo err")
						return
					}
					//余额
					so := slotsmongo.C1A13gameLoginReturn(info.Pid)
					machineType, _ := pp.GetInt("machineType")
					GlobalConnManager.UpdateConnMType(c, info.Pid, machineType)
					fmt.Println(so)
					bytes, err := so.ToBinary()
					if err != nil {
						slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
						return
					}
					c.WriteMessage(messageType, bytes)
				}
				if pc == H5Init {
					sendData := slotsmongo.GatewaySend{}
					sendData.Data = data
					info, ok := GlobalConnManager.GetConnInfo(c)
					if !ok {
						slog.Error("H5Spin Err", "globalConnManager.GetConnInfo err")
						return
					}
					sendData.Pid = info.Pid
					payload, err := json.Marshal(sendData)
					if err != nil {
						slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
						return
					}
					subj := fmt.Sprintf("jdb_%d.%s", info.MachineType, H5Init)
					resp, err := mq.NC().RequestMsg(&nats.Msg{
						Subject: subj,
						Data:    payload,
					}, time.Second*60)
					if err != nil {
						slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
						return
					}
					so := slotsmongo.C1A13h5initResponse(resp.Data)
					fmt.Println(so)
					bytes, err := so.ToBinary()
					if err != nil {
						slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
						return
					}
					c.WriteMessage(messageType, bytes)
				}
				if pc == H5Spin {
					sendData := slotsmongo.GatewaySend{}
					sendData.Data = data
					info, ok := GlobalConnManager.GetConnInfo(c)
					if !ok {
						slog.Error("H5Spin Err", "globalConnManager.GetConnInfo err")
						return
					}
					sendData.Pid = info.Pid
					payload, err := json.Marshal(sendData)
					if err != nil {
						slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
						return
					}
					subj := fmt.Sprintf("jdb_%d.%s", info.MachineType, H5Spin)
					resp, err := mq.NC().RequestMsg(&nats.Msg{
						Subject: subj,
						Data:    payload,
					}, time.Second*60)
					if err != nil {
						slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
						return
					}
					so := slotsmongo.C1A13h5spinResponse(resp.Data)
					bytes, err := so.ToBinary()
					if err != nil {
						slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
						return
					}
					c.WriteMessage(messageType, bytes)
				}
			}
		}
		fmt.Println(reqSo)
	})

	u.OnOpen(func(c *websocket.Conn) {
		slog.Info("OnOpen")
		// 连接打开时可以初始化一些信息，但 pid 可能需要在登录后才能获取
	})

	u.OnClose(func(c *websocket.Conn, err error) {
		slog.Info("OnClose")
		// 连接关闭时从管理器中移除
		GlobalConnManager.RemoveConn(c)
	})

	return u
}
