package ut

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopN(t *testing.T) {
	arr := [][]int{{1, 2, 3, 4, 5}}

	ret := PopN(&arr[0], 7)
	fmt.Println(ret)
	fmt.Println(arr)
}

func TestJoinStr(t *testing.T) {
	ans := JoinStr('-', "hello", "world", "1234")
	fmt.Println(ans)
}

func TestSetField(t *testing.T) {
	plr := struct {
		PID int64
	}{}

	pid := int64(123)
	SetField(&plr, "PID", pid)
	assert.Equal(t, int64(123), plr.PID)
}

func TestMaxLen(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7}
	ans := MaxLen(arr, 5)
	assert.Equal(t, []int{3, 4, 5, 6, 7}, ans)

	arr = []int{1, 2, 3, 4, 5, 6, 7}
	ans = MaxLen(arr, 1)
	assert.Equal(t, []int{7}, ans)

	arr = []int{1, 2, 3, 4, 5, 6, 7}
	ans = MaxLen(arr, 6)
	assert.Equal(t, []int{2, 3, 4, 5, 6, 7}, ans)

	arr = []int{1, 2, 3, 4, 5, 6, 7}
	ans = MaxLen(arr, 7)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, ans)

	arr = []int{1, 2, 3, 4, 5, 6, 7}
	ans = MaxLen(arr, 8)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, ans)
}

func TestClip(t *testing.T) {
	assert.Equal(t, 123, ClipInt(123, 100, 1000))
	assert.Equal(t, 100, ClipInt(99, 100, 1000))
	assert.Equal(t, 1000, ClipInt(1001, 100, 1000))
}

func TestSampleByWeightsPred(t *testing.T) {
	arr := []string{"a", "b", "c"}

	m := map[string]int{}

	for i := 0; i < 60000; i++ {

		s, ok := SampleByWeightsPred(arr, func(s string) int {
			if s == "a" {
				return 1
			}
			if s == "b" {
				return 2
			}
			return 3
		})
		if ok {
			m[s]++
		}
	}

	fmt.Println(m)
}

func TestPopRand(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	for i := 0; i < 3; i++ {
		fmt.Println(PopRand(&arr))
	}
	// fmt.Println(PopRand(&arr))
	fmt.Println(arr)
	assert.Empty(t, arr)
}

func TestErrString(t *testing.T) {
	var err error
	assert.Equal(t, "", ErrString(err))
	err = errors.New("err123")
	assert.Equal(t, "err123", ErrString(err))
}

func TestRound6(t *testing.T) {
	fmt.Println(Round6(1.7999999999999998))
}

func TestFloatStrPtrMul(t *testing.T) {
	f := 12323.14

	s := strconv.FormatFloat(f, 'f', 4, 64)
	fmt.Println(s)

	ss := "3.1415926"
	FloatStrPtrMul(&ss, 2)
	fmt.Println(ss)
}
