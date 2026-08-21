package main

import (
	"bytes"
	"fmt"
	"testing"

	"serve/comm/db"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

func TestSpindata(t *testing.T) {
	s := `{"wp":null,"lw":null,"frl":[3,7,5,4,5,3,5,6,5],"pc":null}`
	d, err := db.Json2BsonD([]byte(s))
	assert.Nil(t, err)

	d = append(d, bson.E{"name", "zzzz"})

	obuf, _ := bson.Marshal(d)

	fmt.Println(bson.Raw(obuf).String())

	// fmt.Println(obj.Lookup("sid").StringValue())

	// fmt.Println(obj.String())

	// db.BsonAppend(obj)
	// obj.String()

}

func TestAppendDoc(t *testing.T) {
	// dst, _ := bson.Marshal(bson.M{"name": "zzz", "age": 19})

	buf := bytes.NewBuffer(nil)
	vw, err := bsonrw.NewBSONValueWriter(buf)
	if err != nil {
		panic(err)
	}

	// vw.WriteDocument()
	// vw.WriteDocument()
	dw, _ := vw.WriteDocument()
	// dw.WriteDocumentElement("aa")
	// dw.WriteDocumentElement()
	fuckw, _ := dw.WriteDocumentElement("fuck")
	fuckw.WriteDouble(3.1415)
	// fuckw.WriteString("fdsafkd")
	namew, _ := dw.WriteDocumentElement("name")

	tt, b, _ := bson.MarshalValue("aaabbbf")
	namew.WriteBinaryWithSubtype(b, byte(tt))
	dw.WriteDocumentEnd()
	// enc, err := bson.NewEncoder(vw)
	// if err != nil {
	// 	panic(err)
	// }

	// enc.Encode(bson.M{"fuck": 3.1415})

	fmt.Println(bson.Raw(buf.Bytes()).String())
}
