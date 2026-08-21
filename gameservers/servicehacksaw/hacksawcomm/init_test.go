package hacksawcomm

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/phuslu/iploc"
	"github.com/stretchr/testify/assert"
)

func TestDoInit(t *testing.T) {

	ps := Variables{}

	ps.SetStr("teststr", "test111")
	assert.Equal(t, "test111", ps.Str("teststr"))

	ps.SetInt("testint", 123)
	assert.Equal(t, 123, ps.Int("testint"))

	ps.SetFloat("testfloat", 3.14)
	assert.Equal(t, 3.14, ps.Float("testfloat"))

	ps.SetCurrency("testcurrency", 1234567.8)
	assert.Equal(t, 1234567.8, ps.Currency("testcurrency"))
	assert.Equal(t, "1,234,567.80", ps.Str("testcurrency"))

	fmt.Println(ps, len(ps.Encode()))
}

func TestFetchMgckey(t *testing.T) {
	// FetchMgckey("vs20olympx", "123456")
}

func TestGetPubIP(t *testing.T) {
	assert.Equal(t, "127.0.0.1", getPubIP(http.DefaultClient))

	loc := string(iploc.Country(net.ParseIP("127.0.0.1")))
	fmt.Println(loc)
}
