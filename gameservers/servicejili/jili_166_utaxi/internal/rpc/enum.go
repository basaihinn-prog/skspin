package rpc

var AckType = map[string]int32{
	"login":         0,
	"exchangeChips": 2,
	"info":          11,
	"spin":          12,
	"buyBonus":      36,
	"heartBeat":     98,
	"serverMsg":     99,
}
