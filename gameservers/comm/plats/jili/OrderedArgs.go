package jili

import (
	"fmt"
	"net/url"
	"strings"

	"serve/comm/db"

	"go.mongodb.org/mongo-driver/bson"
)

type OrderedArgs bson.D

func (oa OrderedArgs) Encode() string {
	var buf strings.Builder

	for _, e := range oa {
		keyEscaped := url.QueryEscape(e.Key)
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(keyEscaped)
		buf.WriteByte('=')
		v := ""
		if _, ok := e.Value.(float64); ok {
			v = fmt.Sprintf("%.2f", e.Value)
		} else {
			v = fmt.Sprint(e.Value)
		}
		buf.WriteString(v)
	}
	return buf.String()
}

func D(ps ...any) OrderedArgs {
	return OrderedArgs(db.D(ps...))
}

func (oa OrderedArgs) Append(key string, value interface{}) OrderedArgs {
	return append(oa, bson.E{key, value})
}
