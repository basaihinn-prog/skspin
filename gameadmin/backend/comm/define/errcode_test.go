package define

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrCode(t *testing.T) {
	ec := NewErrCode("some error from db", 1201)
	s := ec.Error()
	fmt.Println(s)

	ec2, ok := NewErrCodeFromStr(s)
	assert.True(t, ok)
	fmt.Println(ec2.Error())
	assert.Equal(t, s, ec2.Error())
	// arr := errcodeReg.FindStringSubmatch("code:1201#some error from db")
	// fmt.Printf("%+q", arr)

	ec3, ok := NewErrCodeFromStr("some error")
	assert.False(t, ok)
	fmt.Println(ec3)
}
