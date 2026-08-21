package rpc

import (
	"serve/servicepp/ppcomm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoInit(t *testing.T) {

	ps := ppcomm.Variables{}

	ps.SetStr("test", "test111")
	assert.Equal(t, "test111", ps.Str("test"))

	ps.SetInt("test", 123)
	assert.Equal(t, 123, ps.Int("test"))

	ps.SetFloat("test", 3.14)
	assert.Equal(t, 3.14, ps.Float("test"))

	ps.SetCurrency("test", 1234567.8)
	assert.Equal(t, 1234567.8, ps.Currency("test"))
	assert.Equal(t, "1,234,567.80", ps.Str("test"))

}
