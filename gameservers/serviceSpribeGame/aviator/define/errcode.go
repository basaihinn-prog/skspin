package define

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/samber/lo"
)

type ErrCode struct {
	err  string
	code int
}

func NewErrCode(err string, code int) ErrCode {
	lo.Must0(err != "")
	lo.Must0(code != 0)
	return ErrCode{err, code}
}

var (
	errcodeReg = regexp.MustCompile(`^code:(\d+)#(.*)$`)
)

func NewErrCodeFromStr(s string) (ec ErrCode, ok bool) {
	arr := errcodeReg.FindStringSubmatch(s)
	if len(arr) != 3 {
		return
	}
	ec.code, _ = strconv.Atoi(arr[1])
	ec.err = arr[2]
	ok = true
	return
}

func (ec ErrCode) Error() string {
	return fmt.Sprintf("code:%d#%s", ec.code, ec.err)
}

func (ec ErrCode) Code() int {
	return ec.code
}

type IErrcode interface {
	Code() int
}

func CodeErrorEq(err error, code int) bool {
	ec, ok := err.(IErrcode)
	return ok && code == ec.Code()
}
