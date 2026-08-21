//go:build dev
// +build dev

package auth

import (
	"serve/fish_comm/flux"
	errorcode "serve/fish_comm/flux/error-code"
	"serve/service_fish/domain/auth/proto"
	auth_proto "serve/service_fish/domain/auth/proto"

	"github.com/gogo/protobuf/proto"
)

func (s *service) eLoginCall(action *flux.Action, authUsers map[string]*MemberInfo, authTokens, authHostExtIds map[string]string) {
	secWebSocketKey := action.Key().From()
	data := &auth_proto.LoginCall{}

	if err := proto.Unmarshal(action.Payload()[0].([]byte), data); err != nil {
		errorcode.Service.Fatal(secWebSocketKey, Auth_LOGIN_CALL_PROTO_INVALID)
		return
	}

	if data.Token == "" {
		errorcode.Service.Fatal(secWebSocketKey, Auth_TOKEN_EMPTY)
		return
	}

	memberInfo := &MemberInfo{
		HostId:       "DEV",
		WalletType:   -1,
		GuestEnable:  1,
		MemberId:     "DEV",
		MemberName:   "DEV",
		Balance:      99999,
		BillCurrency: "CNY",
	}

	authUsers[secWebSocketKey] = memberInfo
	authTokens[data.Token] = secWebSocketKey
	s.loginSuccess(secWebSocketKey, data.Token, memberInfo)
}
