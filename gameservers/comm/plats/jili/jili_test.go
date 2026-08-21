package jili

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameList(t *testing.T) {
	var p jili

	p.Regist("123456")
	p.FundTransferIn("123456", 1000)
	fmt.Println(p.GetBalance("123456"))

	games, err := p.GetGameList()
	assert.Nil(t, err)
	fmt.Println(games)

	ul, err := p.LaunchGame("123456", "2", "en", false)
	assert.Nil(t, err)
	fmt.Println(ul)
}

func TestLaunch(t *testing.T) {
	var p jili

	ul, err := p.LaunchGame("123456", "2", "en", false)
	assert.Nil(t, err)
	ul = strings.Replace(ul, "https://uat-wbgame.jlfafafa3.com", "https://jilid-rslotszs001.kafa010.com", 1)

	fmt.Println(ul)
}

// url: https://wb-api-2.qu5feb8n.com/api1/CreateMember
// args: "Account=123456&AgentId=ZFRT679_XUU_THB&Key=00000039d608b46c60d9794de3c53a0e5810f1000000"
// resp: "{\"ErrorCode\":10,\"Message\":\"Server Maintaining ... (2171)\",\"Data\":null}"
