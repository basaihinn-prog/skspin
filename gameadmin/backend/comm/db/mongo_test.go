package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMisc(t *testing.T) {
	// DialToMongo("mongodb://127.0.0.1:27017/YingCaiShen", "")
	mongoaddr := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"
	DialToMongo(mongoaddr, "YingCaiShen")

	err := MiscSet("name", "zhangsan")
	assert.Nil(t, err)
	var name string
	err = MiscGet("name", &name)
	assert.Nil(t, err)

	assert.Equal(t, "zhangsan", name)

	MiscSet("age", 19)
	var age int
	MiscGet("age", &age)
	assert.Equal(t, 19, age)

	type objt struct {
		A int
		B string
		C float64
	}
	var obj = objt{
		123, "hello", 3.14,
	}
	MiscSet("obj", obj)

	var objptr = &objt{}
	MiscGet("obj", objptr)

	assert.Equal(t, obj, *objptr)

	var boyibo_balance int64
	MiscGet2("YingCaiShen", "boyibo_balance", &boyibo_balance)
	fmt.Println(boyibo_balance)
}
