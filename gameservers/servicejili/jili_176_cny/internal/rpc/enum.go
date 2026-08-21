package rpc

var AckType = map[string]int32{
	"login":         0,
	"info":          11,
	"spin":          12,
	"free":          36,
	"logout":        1,
	"exchangeChips": 2,
	"last":          3,
	"heartBeat":     98,
	"serverMsg":     99,
}
