package gamedata

import (
	"cmp"
	"log/slog"
	"sync"
)

type TransferLimitMap map[string]float64

var (
	transferLimitMap    TransferLimitMap
	transferLimitMapMtx sync.Mutex
)

func loadTransferLimitMap(tmp TransferLimitMap) (err error) {
	slog.Info("loadTransferLimitMap", "data", tmp)
	transferLimitMapMtx.Lock()
	transferLimitMap = tmp
	transferLimitMapMtx.Unlock()
	return nil
}

func GetTransferLimitMap() TransferLimitMap {
	transferLimitMapMtx.Lock()
	m := transferLimitMap
	transferLimitMapMtx.Unlock()

	return m
}

func (m TransferLimitMap) GetLimit(appID string) float64 {
	if len(m) == 0 {
		return 1e8
	}

	return cmp.Or(m[appID], m["default"], 1e8)
}
