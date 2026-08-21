package handlers

import (
	"fmt"
	"game/duck/ut2/jwtutil"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestIp2loc(t *testing.T) {
	loc := ip2locMap.Get("47.108.238.254")  // CN
	loc2 := ip2locMap.Get("47.108.238.254") // CN
	fmt.Println(loc, loc2)
}

func BenchmarkIp2loc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf [4]byte
		lo.Must(rand.Read(buf[:]))
		ip := net.IPv4(buf[0], buf[1], buf[2], buf[3])
		ip2locMap.Get(ip.String())
		// ip.String()
		// iploc.Country(ip)
	}
}

func TestToken(t *testing.T) {
	token, err := jwtutil.NewTokenWithData(100124, time.Now().Add(12*time.Hour), "89")

	assert.Nil(t, err)
	fmt.Println(token)

	fmt.Println(jwtutil.ParseTokenData(token))
}
