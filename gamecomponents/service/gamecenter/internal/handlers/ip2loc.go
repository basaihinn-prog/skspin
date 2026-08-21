package handlers

import (
	"net"

	"github.com/phuslu/iploc"
)

var (
	ip2locMap = &IP2locMap{
		// m: map[string]string{},
	}
)

type IP2locMap struct {
	// m   map[string]string
	// mtx sync.Mutex
}

func (*IP2locMap) Get(ip string) string {
	// m.mtx.Lock()
	// defer m.mtx.Unlock()

	// if loc, ok := m.m[ip]; ok {
	// 	return loc
	// }

	loc := string(iploc.Country(net.ParseIP(ip)))
	// m.m[ip] = loc
	return loc
}
