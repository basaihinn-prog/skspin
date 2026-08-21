package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"serve/fish_comm/http"
	"serve/fish_comm/mysql"
	auth_proto "serve/service_fish/domain/auth/proto"
	"strconv"
	"strings"
	"sync"
	"time"

	"serve/fish_comm/common"
	common_proto "serve/fish_comm/common/proto"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"

	"github.com/gogo/protobuf/proto"
)

const (
	EMSGID_eLoginCall          = "EMSGID_eLoginCall"
	EMSGID_eLoginRecall        = "EMSGID_eLoginRecall"
	EMSGID_eBroadcastPlayerIn  = "EMSGID_eBroadcastPlayerIn"
	EMSGID_eBroadcastPlayerOut = "EMSGID_eBroadcastPlayerOut"
	ActionUserDisconnect       = "ActionUserDisconnect"
	ActionLoginSuccess         = "ActionLoginSuccess"
	ActionLoginFailed          = "ActionLoginFailed"
	ActionCheckToken           = "ActionCheckToken"
	ActionUserHostExtIdSave    = "ActionUserHostExtIdSave"
	ActionAuthUserInfoAsk      = "ActionAuthUserInfoAsk"
)

var Service = &service{
	Id:       "AuthService",
	poolSize: os.Getenv("POOL_SIZE"),
	in:       make(chan *flux.Action, common.Service.ChanSize),
	mutex:    sync.Mutex{},
}

type service struct {
	Id       string
	poolSize string
	in       chan *flux.Action
	mutex    sync.Mutex
}

type ApiInfo struct {
	BaseUrl       string `json:"base_url"`
	Auth          string `json:"auth"`
	Logout        string `json:"logout"`
	Bet           string `json:"bet"`
	Result        string `json:"dbResult"`
	RefundBet     string `json:"refundbet"`
	BonusAward    string `json:"bonusaward"`
	ResultEx      string `json:"resultex"`
	GetBalance    string `json:"getbalance"`
	EnableHttpLog bool   `json:"enable_http_log"`
	RetryTimes    int    `json:"retry_times"`
	SignIn        string `json:"signin"`
}

// db model
type hostId struct {
	HostId              string
	WalletType          int
	ApiInfo             string
	GuestEnable         int
	GuestBalance        uint64
	Status              int
	FishDepositMultiple uint64
	FishShareType       int
	BillCurrency        string
}

type MemberInfo struct {
	HostId          string
	WalletType      int
	GuestEnable     int
	StatusCode      int    `json:"status_code"`
	MemberId        string `json:"member_id"`
	MemberName      string `json:"member_name"`
	Balance         uint64 `json:"balance"`
	AccountType     int    `json:"account_type"`
	RealToken       string
	DepositMultiple uint64
	FishShareType   int
	BillCurrency    string
}

type MemberStatus struct {
	Alive bool
}

func init() {
	Service.run()
	logger.Service.Zap.Infow("Service created.",
		"Service", Service.Id,
		"Chan", fmt.Sprintf("%p", Service.in),
	)
}

func (s *service) run() {
	flux.Register(s.Id, s.in)
	authUsers := make(map[string]*MemberInfo)
	authTokens := make(map[string]string)
	authHostExtIds := make(map[string]string)

	poolSize := 10

	if s.poolSize != "" {
		if poolSize, _ = strconv.Atoi(s.poolSize); poolSize == 0 {
			poolSize = 10
		}
	}

	for i := 0; i < poolSize*10; i++ {
		go func() {
			for action := range s.in {
				s.handleAction(action, authUsers, authTokens, authHostExtIds)
			}
		}()
	}
}

func (s *service) handleAction(
	action *flux.Action,
	authUsers map[string]*MemberInfo,
	authTokens map[string]string,
	authHostExtIds map[string]string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch action.Key().Name() {
	case ActionUserDisconnect:
		secWebSocketKey := action.Payload()[0].(string)
		delete(authUsers, secWebSocketKey)
		delete(authTokens, secWebSocketKey)
		delete(authHostExtIds, secWebSocketKey)

	case EMSGID_eLoginCall:
		s.eLoginCall(action, authUsers, authTokens, authHostExtIds)

	case ActionCheckToken:
		isTokenValid := action.Payload()[0].(chan bool)
		token := action.Payload()[1].(string)

		if _, ok := authTokens[token]; ok {
			isTokenValid <- ok
		} else {
			isTokenValid <- ok
		}

		close(isTokenValid)
		isTokenValid = nil

	case EMSGID_eBroadcastPlayerOut:
		secWebSocketKey := action.Key().From()
		gameRoomUuid := action.Payload()[0].(string)
		seatId := action.Payload()[1].(uint32)
		balance := action.Payload()[2].(uint64)

		if v, ok := authUsers[secWebSocketKey]; ok {
			flux.Send(EMSGID_eBroadcastPlayerOut, s.Id, gameRoomUuid,
				s.BroadcastPlayerOut(gameRoomUuid, v.MemberId, v.MemberName, seatId, balance),
			)
		}

	case ActionUserHostExtIdSave:
		secWebSocketKey := action.Key().From()
		hostExtId := action.Payload()[0].(string)
		isDone := action.Payload()[1].(chan bool)

		authHostExtIds[secWebSocketKey] = hostExtId

		isDone <- true
		close(isDone)

	case ActionAuthUserInfoAsk:
		secWebSocketKey := action.Payload()[0].(string)
		balance := action.Payload()[1].(uint64)

		if v, ok := authUsers[secWebSocketKey]; ok {
			data := &auth_proto.UserInfoAsk{
				Msgid: common_proto.EMSGID_eMemberInfoAsk,
			}

			data.Player = append(data.Player, &auth_proto.UserInfo{
				PlayerId:    v.MemberId,
				PlayerName:  v.MemberName,
				PlayerCent:  balance,
				SeatId:      0,
				BetRateLine: nil,
				Visible:     nil,
			})

			dataByte, _ := proto.Marshal(data)

			flux.Send(ActionAuthUserInfoAsk, s.Id, secWebSocketKey, dataByte)
		}

	case EMSGID_eLoginRecall:
		secWebSocketKey := action.Key().From()
		isSuccess := action.Payload()[0].(bool)

		isFind := false

		for k, v := range authTokens {
			if v == secWebSocketKey {
				isFind = true
				s.eLoginRecall(secWebSocketKey, k, isSuccess)
				break
			}
		}

		if !isFind {
			logger.Service.Zap.Errorw(Auth_TOKEN_NOT_FOUND,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, Auth_TOKEN_NOT_FOUND)
		}
	}
}

func (s *service) parseToken(secWebSocketKey, token string, authHostExtIds map[string]string) (realToken, hostExtId string, isGuestToken bool, err error) {
	ss := strings.Split(token, "(*--)")

	if len(ss) < 2 {
		return "", "", false, errors.New(fmt.Sprintf("Token %s format error.", token))
	}

	if v, ok := authHostExtIds[secWebSocketKey]; ok {
		if v != ss[1] {
			return ss[0], ss[1], false, errors.New(fmt.Sprintf("Unknown host_ext_id %s.", token))
		}

		if ss[0] == "" {
			return ss[0], ss[1], ss[0] == "" && v == ss[1], nil
		}

		return ss[0], ss[1], false, nil
	} else {
		return ss[0], ss[1], false, errors.New(fmt.Sprintf("host_ext_id %s not found.", token))
	}
}

func (s *service) BroadcastPlayerIn(roomUuid string) *auth_proto.BroadcastPlayerIn {
	return &auth_proto.BroadcastPlayerIn{
		Msgid:    common_proto.EMSGID_eBroadcastPlayerIn,
		RoomUuid: roomUuid,
	}
}

func (s *service) BroadcastPlayerOut(gameRoomUuid, memberId, memberName string, seatId uint32, balance uint64) []byte {
	data := &auth_proto.BroadcastPlayerOut{
		Msgid:    common_proto.EMSGID_eBroadcastPlayerOut,
		RoomUuid: gameRoomUuid,
	}

	data.Player = &auth_proto.UserInfo{
		PlayerId:    memberId,
		PlayerCent:  balance,
		PlayerName:  memberName,
		SeatId:      seatId,
		BetRateLine: nil,
		Visible:     nil,
	}

	dataByte, _ := proto.Marshal(data)
	return dataByte
}

func (s *service) guestMode(secWebSocketKey, token string, dbResult *hostId) (memberInfo *MemberInfo, ok bool) {
	if token == "" {
		logger.Service.Zap.Errorw(Auth_TOKEN_EMPTY,
			"GameUser", secWebSocketKey,
			"HostId", dbResult.HostId,
			"GuestBalance", dbResult.GuestBalance,
		)
		errorcode.Service.Fatal(secWebSocketKey, Auth_TOKEN_EMPTY)
		return nil, false
	}

	memberInfo = &MemberInfo{
		HostId:       dbResult.HostId,
		WalletType:   -1,
		GuestEnable:  1,
		MemberId:     "Guest",
		MemberName:   "Guest",
		Balance:      dbResult.GuestBalance,
		BillCurrency: dbResult.BillCurrency,
	}
	return memberInfo, true
}

func (s *service) loginSuccess(secWebSocketKey, token string, m *MemberInfo) {
	flux.Send(ActionLoginSuccess, s.Id, secWebSocketKey, m)
	logger.Service.Zap.Infow("Auth Login Success.",
		"GameUser", secWebSocketKey,
		"HostId", m.HostId,
		"MemberId", m.MemberId,
		"MemberName", m.MemberName,
		"BillCurrency", m.BillCurrency,
	)
}

func (s *service) loginFailed(secWebSocketKey, token string) {
	flux.Send(ActionLoginFailed, s.Id, secWebSocketKey, secWebSocketKey)
	time.Sleep(time.Second)
	s.eLoginRecall(secWebSocketKey, token, false)
	logger.Service.Zap.Infow("Auth Login Failed.", "GameUser", secWebSocketKey)
}

func (s *service) eLoginRecall(secWebSocketKey, token string, isSuccess bool) {
	data := &auth_proto.LoginRecall{
		Msgid: common_proto.EMSGID_eLoginRecall,
		Token: token,
	}

	if isSuccess {
		data.StatusCode = common_proto.Status_kSuccess
	} else {
		data.StatusCode = common_proto.Status_kInvalid
	}

	dataByte, _ := proto.Marshal(data)
	flux.Send(EMSGID_eLoginRecall, s.Id, secWebSocketKey, dataByte)
}

func (s *service) ColumnVisible(status uint32) *auth_proto.UserInfo_Column {
	// ID, NAME, BALANCE, AVATAR
	// (1 ASCII(49):show, 0 ASCII(48):hidden)
	// All show:15 (1111)

	data := &auth_proto.UserInfo_Column{}

	for k, v := range []rune(strconv.FormatUint(uint64(status), 2)) {
		switch k {
		case 0: // ID
			data.Id = v == 49

		case 1: // NAME
			data.Name = v == 49

		case 2: // BALANCE
			data.Balance = v == 49

		case 3: // AVATAR
			data.Avatar = v == 49
		}
	}

	return data
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}

func (s *service) eLoginCall(action *flux.Action, authUsers map[string]*MemberInfo, authTokens, authHostExtIds map[string]string) {
	secWebSocketKey := action.Key().From()
	data := &auth_proto.LoginCall{}

	if err := proto.Unmarshal(action.Payload()[0].([]byte), data); err != nil {
		logger.Service.Zap.Errorw(Auth_LOGIN_CALL_PROTO_INVALID,
			"GameUser", secWebSocketKey,
			"MemberId", data.MemberId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, Auth_LOGIN_CALL_PROTO_INVALID)
		return
	}

	logger.Service.Zap.Infow("Auth Login Start.",
		"GameUser", secWebSocketKey,
		"MemberId", data.MemberId,
	)

	if data.Token == "" {
		logger.Service.Zap.Errorw(Auth_TOKEN_EMPTY,
			"GameUser", secWebSocketKey,
			"MemberId", data.MemberId,
		)
		errorcode.Service.Fatal(secWebSocketKey, Auth_TOKEN_EMPTY)
		return
	}

	realToken, hostExtId, isGuestToken, err := s.parseToken(secWebSocketKey, data.Token, authHostExtIds)

	if err != nil {
		logger.Service.Zap.Errorw(Auth_TOKEN_INVALID,
			"GameUser", secWebSocketKey,
			"MemberId", data.MemberId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, Auth_TOKEN_INVALID)
		return
	}

	dbResult := &hostId{}

	if db, err := mysql.Repository.GameDB(hostExtId); err == nil {
		if ok := db.
			Table("host_id").
			Select("host_id, wallet_type, api_info, guest_enable, guest_balance, status, fish_deposit_multiple, fish_share_type, host_id.bill_currency").
			Where("host_ext_id = ?", hostExtId).
			Scan(dbResult).
			RowsAffected; ok != 1 {
			logger.Service.Zap.Errorw(Auth_HOST_ID_NOT_FOUND,
				"GameUser", secWebSocketKey,
				"MemberId", data.MemberId,
			)
			errorcode.Service.Fatal(secWebSocketKey, Auth_HOST_ID_NOT_FOUND)
			return
		}
	} else {
		logger.Service.Zap.Errorw(Auth_GAME_DB_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"HostExtId", hostExtId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, Auth_GAME_DB_NOT_FOUND)
		return
	}

	if dbResult.GuestEnable == 1 && isGuestToken {
		if memberInfo, ok := s.guestMode(secWebSocketKey, data.Token, dbResult); ok {
			authUsers[secWebSocketKey] = memberInfo
			authTokens[data.Token] = secWebSocketKey
			s.loginSuccess(secWebSocketKey, data.Token, memberInfo)
		}
		return
	}

	apiInfo := &ApiInfo{}

	if err := json.Unmarshal([]byte(dbResult.ApiInfo), apiInfo); err != nil {
		logger.Service.Zap.Errorw(Auth_JSON_API_INFO_INVALID,
			"GameUser", secWebSocketKey,
			"MemberId", data.MemberId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, Auth_JSON_API_INFO_INVALID)
		return
	}

	parameter := fmt.Sprintf("access_token=%s&step=1", realToken)
	url := fmt.Sprintf("%s/%s?%s", apiInfo.BaseUrl, apiInfo.Auth, parameter)

	memberInfo := &MemberInfo{
		HostId:          dbResult.HostId,
		WalletType:      dbResult.WalletType,
		GuestEnable:     0,
		RealToken:       realToken,
		DepositMultiple: dbResult.FishDepositMultiple,
		FishShareType:   dbResult.FishShareType,
		BillCurrency:    dbResult.BillCurrency,
	}

	var body []byte

	if body, err = http.Service.Get(url); err != nil {
		logger.Service.Zap.Errorw(Auth_API_URL_GET_FAILED,
			"GameUser", secWebSocketKey,
			"MemberId", data.MemberId,
			"URL", url,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, Auth_API_URL_GET_FAILED)
		return
	}

	if err := json.Unmarshal(body, memberInfo); err != nil {
		logger.Service.Zap.Errorw(Auth_JSON_MEMBER_INFO_INVALID,
			"GameUser", secWebSocketKey,
			"MemberId", data.MemberId,
			"URL", url,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, Auth_JSON_MEMBER_INFO_INVALID)
		return
	}

	logger.Service.Zap.Debugw("Get GameUser Profile",
		"GameUser", secWebSocketKey,
		"URL", url,
		"Profile", memberInfo,
		"Body", body,
	)

	// Process Old Share Wallet
	if memberInfo.WalletType == 0 && memberInfo.FishShareType == 1 {
		memberInfo.WalletType = 2
	}

	if memberInfo.MemberName == "" {
		memberInfo.MemberName = memberInfo.MemberId
	}

	if memberInfo.StatusCode == 0 && memberInfo.MemberId != "" {
		authUsers[secWebSocketKey] = memberInfo
		authTokens[data.Token] = secWebSocketKey
		s.loginSuccess(secWebSocketKey, data.Token, memberInfo)
		return
	}

	s.loginFailed(secWebSocketKey, data.Token)
}
