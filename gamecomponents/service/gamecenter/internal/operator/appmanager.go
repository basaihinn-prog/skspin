package operator

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/comm/define"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/logger"
	"game/duck/mongodb"
	"game/duck/ut2"
	"game/service/gamecenter/internal/gcdb"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppManager struct {
	// appid -> app
	apps *ut2.SyncMap[string, *MemApp]

	// pid -> plr
	players *ut2.SyncMap[int64, *MemPlr]
}

var (
	AppMgr *AppManager
)

func NewAppManagerAndInit() {
	// hour := 1 * time.Hour

	ret := &AppManager{
		apps: ut2.NewSyncMap[string, *MemApp](),
		// players: ut2.NewActiveMap[int64, *DocPlayer](hour, hour),
		players: ut2.NewSyncMap[int64, *MemPlr](),
	}

	ret.loadApps()
	go ret.Update()

	AppMgr = ret
}

func (am *AppManager) Update() {
	for {
		time.Sleep(10 * time.Second)
		AppMgr.loadApps()
	}
}

func (am *AppManager) loadApps() {
	// logger.Info("=== appManager init ===")
	all := []*comm.Operator{}
	err := gcdb.NewOtherDB("GameAdmin").Collection("AdminOperator").FindAll(bson.M{}, &all)
	if err != nil {
		logger.Err("AppManager loadApps", err)
		return
	}
	for _, op := range all {
		if memapp, _ := am.apps.Load(op.AppID); memapp == nil {
			app := &MemApp{
				AppID:               op.AppID,
				ParentAppID:         op.Name,
				AppSecret:           op.AppSecret,
				Address:             op.Address,
				WalletMode:          op.WalletMode,
				PublishHistory:      op.PublishHistory,
				Status:              op.Status,
				CurrencyKey:         op.CurrencyKey,
				OperatorType:        op.OperatorType,
				players:             ut2.NewSyncMap[string, *MemPlr](),
				ReviewStatus:        op.ReviewStatus,
				CurrencyCtrlStatus:  op.CurrencyCtrlStatus,
				RTPOff:              op.RTPOff,
				PlayerRTPSettingOff: op.PlayerRTPSettingOff,
				HighRTPOff:          op.HighRTPOff,
			}
			am.apps.Store(app.AppID, app)
		} else {
			if memapp.AppSecret != op.AppSecret {
				memapp.AppSecret = op.AppSecret
			}

			if memapp.WalletMode != op.WalletMode {
				memapp.WalletMode = op.WalletMode
			}

			if memapp.Address != op.Address {
				memapp.Address = op.Address
			}

			if memapp.PublishHistory != op.PublishHistory {
				memapp.PublishHistory = op.PublishHistory
			}
			if memapp.Status != op.Status {
				memapp.Status = op.Status
			}
			if memapp.CurrencyKey != op.CurrencyKey {
				memapp.CurrencyKey = op.CurrencyKey
			}
			if memapp.ReviewStatus != op.ReviewStatus {
				memapp.ReviewStatus = op.ReviewStatus
			}
			if memapp.CurrencyCtrlStatus != op.CurrencyCtrlStatus {
				memapp.CurrencyCtrlStatus = op.CurrencyCtrlStatus
			}
			if memapp.RTPOff != op.RTPOff {
				memapp.RTPOff = op.RTPOff
			}
			if memapp.PlayerRTPSettingOff != op.PlayerRTPSettingOff {
				memapp.PlayerRTPSettingOff = op.PlayerRTPSettingOff
			}
			if memapp.HighRTPOff != op.HighRTPOff {
				memapp.HighRTPOff = op.HighRTPOff
			}
		}
	}
}

func (am *AppManager) GetApp(AppID string) *MemApp {
	app, ok := am.apps.Load(AppID)
	if !ok {
		return nil
	}

	return app
}

func (am *AppManager) UpdatePlrLoginTime(pid int64) (err error) {
	coll := gcdb.CollPlayers.Coll()

	var (
		plr DocPlayer
		now = time.Now()
	)
	err = coll.FindOneAndUpdate(context.TODO(), db.ID(pid), db.D("$set", db.D("LoginAt", mongodb.NewTimeStamp(now))), options.FindOneAndUpdate().SetReturnDocument(options.Before).SetProjection(db.D("LoginAt", 1))).Decode(&plr)

	if err != nil {
		return
	}

	if !ut.IsSameDate(plr.LoginAt.AsTime(), now) {
		slotsmongo.IncLoginCount(plr.AppID)
	}

	return nil
}

func (am *AppManager) GetPlr(pid int64) (memplr *MemPlr, err error) {
	memplr, ok := am.players.Load(pid)
	if !ok {
		memplr, err = loadPlrByPid(pid)
		if err != nil {
			return
		}

		lo.Must0(pid == memplr.Pid)
		am.players.Store(memplr.Pid, memplr)
		if app, ok := am.apps.Load(memplr.AppID); ok {
			app.players.Store(memplr.Uid, memplr)
		}
	}

	return
}

func (am *AppManager) GetPlr2(app *MemApp, uid string) (memplr *MemPlr, err error) {
	memplr, ok := app.players.Load(uid)
	if ok {
		return
	}

	memplr, err = loadPlrByUid(uid, app.AppID)
	if err != nil {
		lo.Must0(memplr == nil)
		return
	}

	lo.Must0(uid == memplr.Uid)
	lo.Must0(app.AppID == memplr.AppID)
	app.players.Store(memplr.Uid, memplr)
	am.players.Store(memplr.Pid, memplr)
	return
}

func (am *AppManager) EnsureUserExists(app *MemApp, userid string) (plr *MemPlr, err error) {
	app.regLock.Lock()
	defer app.regLock.Unlock()

	plr, _ = app.players.Load(userid)
	if plr != nil {
		if plr.Status == 0 {
			err = define.NewErrCode("Player status is close", 2002)
			return nil, err
		}
		return
	}

	coll := gcdb.CollPlayers.Coll()

	plr, err = loadPlrByUid(userid, app.AppID)
	if err == nil {
		if plr.Status == 0 {
			err = define.NewErrCode("Player status is close", 2002)
		}
		lo.Must0(plr != nil)
		app.players.Store(plr.Uid, plr)
		am.players.Store(plr.Pid, plr)
		return
	}

	pid, err := gcdb.NextID(gcdb.CollPlayers, 100000)
	if err != nil {
		return
	}

	one := DocPlayer{
		Pid: pid,
		// Language: lang,
		AppID:    app.AppID,
		Uid:      userid,
		Status:   1,
		CreateAt: mongodb.NowTimeStamp(),
	}
	_, err = coll.InsertOne(context.TODO(), one)
	if err != nil {
		return
	}
	plr = &MemPlr{
		AppID:  one.AppID,
		Pid:    one.Pid,
		Uid:    one.Uid,
		Status: one.Status,
	}

	app.players.Store(plr.Uid, plr)
	am.players.Store(plr.Pid, plr)
	slotsmongo.IncRegistCount(app.AppID)
	return
}

type PlayerStatus struct {
	UserID string
	Status int64
	AppID  string
}

func (am *AppManager) UpdateCacheUserStatus(app *MemApp, player PlayerStatus) (err error) {
	app.regLock.Lock()
	defer app.regLock.Unlock()

	plr, _ := app.players.Load(player.UserID)
	if plr != nil {
		plr.Status = player.Status
		app.players.Store(plr.Uid, plr)
		am.players.Store(plr.Pid, plr)
	}

	return
}

/*
func (am *AppManager) GetBalance(ctx context.Context, req *gamepb.GetBalanceReq) (*gamepb.GetBalanceResp, error) {
	memplr, err := am.GetPlr(req.Pid)
	if err != nil {
		return nil, err
	}

	app, ok := am.apps.Load(memplr.AppID)
	if !ok {
		return nil, errors.New("app not found.")
	}

	// resp, err := app.api.GetBalance(&gamepb.GetBalanceUidReq{
	// 	UserID: memplr.Uid,
	// })
	balance, err := app.api.GetBalance(memplr)
	if err != nil {
		return nil, err
	}

	return &gamepb.GetBalanceResp{
		Balance: balance,
		Uid:     memplr.Uid,
	}, nil
}

func (am *AppManager) ModifyGold(ctx context.Context, req *gamepb.ModifyGoldReq) (*gamepb.GetBalanceResp, error) {
	plr, err := am.GetPlr(req.Pid)
	if err != nil {
		return nil, err
	}

	app, ok := am.apps.Load(plr.AppID)
	if !ok {
		return nil, errors.New("app not found.")
	}

	before := time.Now()

	balance, err := app.api.ModifyGold(plr, req.Change)

	if err != nil {
		insertModifyLog(plr, req, before, 0, false)
	} else {
		insertModifyLog(plr, req, before, balance, true)
	}

	if err != nil {
		return nil, err
	}

	return &gamepb.GetBalanceResp{Balance: balance, Uid: plr.Uid}, nil
}
*/
