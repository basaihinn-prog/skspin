package main

import (
	"encoding/json"
	"errors"
	"game/comm/mux"
	"game/comm/ut"
	"game/duck/ut2/fileutil"
	"net/http"
	"os"
)

func init() {

	// 默认配置
	RegMsgProc("/AdminInfo/EditConfig", "添加修改删除配置文件", "AdminInfo", ChangeConfig, &ConfigHandleParam{})
	mux.RegHttpWithSample("/AdminInfo/GetConfig", "获取Json配置文件", "AdminInfo", GetConfig, &ConfigHandleParam{})

}

type ConfigHandleParam struct {
	FileName string
	Context  []map[string]string
}

type GameBetJson struct {
	Context map[string]*GameBet
}

type ConfigList struct {
	Context []map[string]string
}

func GetConfig(r *http.Request, ps ConfigHandleParam, ret *ConfigList) (err error) {

	fileContext, err := os.ReadFile("./config/" + ps.FileName)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(fileContext, &ret.Context)
	if err != nil {
		return nil
	}

	// 语言不在配置不在这里  在admin_get_lang.go里面
	return err
}
func ChangeConfig(ctx *Context, ps ConfigHandleParam, ret *GameBetJson) (err error) {

	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New("没有权限")
	}

	err = ConfigJsonIOWrite(ps.FileName, ps.Context)

	return err
}

// 上传多语言
func uploadConfig(w http.ResponseWriter, r *http.Request) {
	//----------------------------------处理鉴权--------------------------
	cker := &checker{}
	pid, username := cker.GetPidAndUname(r.Header.Get("Authorization"))

	context := &Context{
		Request:  r,
		PID:      pid,
		Username: username,
	}

	_, ok := IsAdminUser(context)
	if !ok {
		ut.HttpReturnJson(w, map[string]string{
			"code": "0",
			"msg":  "权限不足",
			"data": "",
		})
		return
	}

	// -----------------------------读取上传文件-------------------------
	file, _, err := r.FormFile("file")
	fileName := r.FormValue("fileName")

	excel, _ := fileutil.ReadExcel(file)

	if err != nil {
		ut.HttpReturnJson(w, map[string]string{
			"code": "0",
			"msg":  "上传错误" + err.Error(),
			"data": "",
		})
		return
	}

	// ---------------------------写入文件-------------------------
	err = ConfigJsonIOWrite(fileName, excel)
	if err != nil {

		ut.HttpReturnJson(w, map[string]string{
			"code": "0",
			"msg":  "修改失败" + err.Error(),
			"data": "",
		})
		return
	}

	ut.HttpReturnJson(w, map[string]string{
		"code": "1",
		"msg":  "完成",
		"data": "",
	})
}

func ConfigJsonIOWrite(fileName string, writeContext []map[string]string) (err error) {

	create, err := os.Create("./config/" + fileName)
	if err != nil {
		return err
	}

	defer create.Close()

	encoder := json.NewEncoder(create)

	err = encoder.Encode(writeContext)
	return err
}
