package main

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func init() {
	RegMsgProc("/AdminInfo/AddLang", "添加多语言", "AdminInfo", AddLange, &LangeParams{})
	RegMsgProc("/AdminInfo/RemoveLang", "删除多语言", "AdminInfo", removeLange, &LangeParams{})
	RegMsgProc("/AdminInfo/EditLang", "修改多语言", "AdminInfo", editLange, &LangeParams{})
	RegMsgProc("/AdminInfo/SelectLang", "查询多语言", "AdminInfo", selectLange, &LangeParams{})
	//RegMsgProc("/AdminInfo/UploadLang", "上传多语言", "AdminInfo", uploadLange, &LangeParams{})
}

type LangeParams struct {
	PageSize       int64               `json:"PageSize"`
	Page           int64               `json:"Page"`
	Filter         string              `json:"filter"`
	CollList       string              `json:"CollList"`
	LanguageConfig []map[string]string `json:"LanguageConfig"`
}

type LangeResponse struct {
	List []map[string]string
}

func AddLange(ctx *Context, ps LangeParams, ret *GetplayRTPListResults) (err error) {

	return errors.New("7777")
}

func removeLange(ctx *Context, ps LangeParams, ret *GetplayRTPListResults) (err error) {
	return errors.New("7777")

}

func editLange(ctx *Context, ps LangeParams, ret *GetplayRTPListResults) (err error) {

	fmt.Println(ps.LanguageConfig)

	return err
}

func selectLange(ctx *Context, ps LangeParams, ret *LangeResponse) (err error) {

	coll := DB.Collection("Lange")
	query := bson.M{}
	if ps.Filter != "" {
		fieldList := strings.Split(ps.Filter, ",")
		fieldQueryList := []bson.M{}
		for _, s := range fieldList {
			fieldQueryList = append(fieldQueryList, bson.M{s: bson.M{"$regex": ps.Filter}})
		}
		query["$or"] = fieldList
	}
	err = coll.FindAll(query, &ret.List)

	return err

}

// 上传多语言
//func uploadLange(w http.ResponseWriter, r *http.Request) {
//
//}
