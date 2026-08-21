//go:build staging || prod
// +build staging prod

package auth

import (
	"encoding/json"
	"fmt"
	auth_proto "serve/service_fish/domain/auth/proto"

	"serve/fish_comm/flux"
	errorcode "serve/fish_comm/flux/error-code"
	"serve/fish_comm/flux/http"
	"serve/fish_comm/flux/logger"
	"serve/fish_comm/flux/mysql"

	"github.com/gogo/protobuf/proto"
)

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
