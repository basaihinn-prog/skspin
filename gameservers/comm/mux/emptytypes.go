package mux

import (
	"reflect"
)

type EmptyResult struct{}
type EmptyParams struct{}

var EmptyResultType = reflect.TypeOf(EmptyResult{})
var EmptyResultValue = reflect.ValueOf(EmptyResult{})
var EmptyParamsValue = reflect.ValueOf(EmptyParams{})
var EmptyParamsType = reflect.TypeOf(EmptyParams{})

func IsEmptyParamsType(t reflect.Type) bool {
	t = Indirect(t)

	return EmptyParamsType == t
}
