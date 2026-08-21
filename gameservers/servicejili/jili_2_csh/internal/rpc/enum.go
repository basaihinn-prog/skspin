package rpc

var AckType = map[string]int32{
	"login":         0,
	"logout":        1,
	"exchangeChips": 2,
	"last":          3,
	"info":          11,
	"spin":          12,
	"spinEnd":       13,
	"buyBonus":      36,
	"heartBeat":     98,
	"serverMsg":     99,
}
