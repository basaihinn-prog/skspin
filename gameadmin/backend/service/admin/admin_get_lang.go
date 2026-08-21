package main

import (
	"game/comm/mux"
	"net/http"
)

func init() {
	mux.RegHttpWithSample("/AdminInfo/LangGet", "获取多语言", "AdminInfo", getLangList, &GetLangListParams{})
}

type GetLangListParams struct {
	Lang      string
	PageIndex int64
	PageSize  int64
}

type GetLangListResults struct {
	List  []map[string]string
	Count int64
}

func getLangList(_ *http.Request, ps GetLangListParams, ret *GetLangListResults) (err error) {
	ret.List = Language

	return err
}
