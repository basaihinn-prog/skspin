package mux

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"serve/comm/define"
)

type A struct{}

func (a A) Do() error {
	return nil
}

func (a A) Hello() error {
	return nil
}

func TestFunction(t *testing.T) {
	// v := reflect.ValueOf(A{})

	// obj := A{}
	// typ := v.Type()
	// nmethod := typ.NumMethod()

	// for i := 0; i < nmethod; i++ {
	// 	method := typ.Method(i)
	// 	fmt.Println(method.Name)

	// 	methodv := v.Method(i)
	// }
}

func TestJson(t *testing.T) {
	buf := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	m := define.M{
		"data": buf,
		"name": "xiaohong",
		"age":  19,
	}

	out, _ := json.Marshal(m)
	// assert.Nil(t, err)
	os.Stdout.Write(out)

	var n struct {
		Data []byte
		Name string
		Age  int
	}
	json.Unmarshal(out, &n)

	fmt.Println(n)
}
