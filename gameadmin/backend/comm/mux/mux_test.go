package mux

import (
	"testing"
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
