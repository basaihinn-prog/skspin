package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/mux"
	"game/comm/ut"
	"game/duck/lang"
	"game/duck/logger"
	"game/duck/ut2/fileutil"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/LangAdd", "添加多语言", "AdminInfo", addLang, &addLangParams{})
}

type addLangParams struct {
	Langs map[string]lang.Lang
}

type addLangResults struct {
	ErrLangs []lang.Lang
}

func addLang(ctx *Context, ps addLangParams, ret *addLangResults) (err error) {
	if _, ok := IsAdminUser(ctx); !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	ret.ErrLangs = make([]lang.Lang, 0, 1)
	arr := make([]interface{}, 0, 1000)
	coll := db.Collection2("game", "Lang")
	for _, l := range ps.Langs {
		if l.LangMap["th"] == "" || l.ZH == "" {
			ret.ErrLangs = append(ret.ErrLangs, l)
			continue
		}
		arr = append(arr, l)
	}
	if len(arr) != 0 {
		for _, v := range ps.Langs {
			for k, _ := range v.LangMap {
				coll.UpdateMany(context.TODO(), bson.M{k: bson.M{"$exists": false}}, bson.M{"$set": bson.M{k: ""}})
			}
		}

		// 向输入库中插入中文解释   由于数据不完善 暂时先不处理

		//reslut, err := coll.InsertMany(context.T ODO(), arr, opts)
		//if err != nil {
		//	for k, l := range ps.Langs {
		//		var ik interface{} = k
		//		if !lo.Contains(reslut.InsertedIDs, ik) {
		//			ret.ErrLangs = append(ret.ErrLangs, l)
		//		}
		//	}
		//}
		return err
	}
	return nil
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {
	formFile, _, _ := r.FormFile("file")
	excel, err := fileutil.ReadExcel(formFile)
	if err != nil {
		return
	}

	// 对应的国家编码
	WORLD_COUNTRY := map[string]string{
		"英文":    "EN",
		"丹麦文":   "DA",
		"德文":    "DE",
		"西班牙文":  "ES",
		"芬兰文":   "FI",
		"法文":    "FR",
		"印尼文":   "IDR",
		"意大利文":  "IT",
		"日文":    "JA",
		"韩文":    "KO",
		"荷兰文":   "NL",
		"挪威文":   "NO",
		"波兰文":   "PL",
		"葡萄牙文":  "PT",
		"罗马尼亚文": "RO",
		"俄文":    "RU",
		"瑞典文":   "SV",
		"泰文":    "TH",
		"土耳其文":  "TR",
		"越南文":   "VI",
		"中文":    "ZH",
		"缅甸文":   "MY",
	}

	var toDB []interface{}
	for _, item := range excel {
		cellItem := make(map[string]interface{})
		for key, value := range item {

			countryCod := WORLD_COUNTRY[key]

			if countryCod != "" {
				cellItem[countryCod] = value
			} else {
				if key == "类型" {

					cellItem["Permission"] = value
				}
			}
		}
		toDB = append(toDB, cellItem)
	}
	_, err = db.Collection2("game", "Lang2").DeleteMany(context.TODO(), bson.M{})
	_, err = db.Collection2("game", "Lang2").InsertMany(context.TODO(), toDB)

	if err != nil {
		return
	}

	ut.HttpReturnJson(w, excel)
}

func saveUploadFile(w http.ResponseWriter, r *http.Request) {
	var ret mux.Response

	timestrap := time.Now().Unix()

	file, fileHeader, _ := r.FormFile("file")
	gameId := r.FormValue("game_id")    // 获取上传游戏ID
	language := r.FormValue("language") // 获取上传的语言
	manufacturer := r.FormValue("manufacturer")
	id := r.FormValue("id")
	fileType := path.Ext(fileHeader.Filename) // 获取文件类型
	fileName := fmt.Sprintf("%s-%s-%d%s", gameId, language, timestrap, fileType)

	// 确定保存图片的目录
	gameFileDir := "./BHdownload/"
	if err := os.MkdirAll(gameFileDir, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if manufacturer != "" {
		fileName = fmt.Sprintf("%s-%s%s", strings.ToUpper(manufacturer), gameId, fileType)
	}

	// 构建图片保存的完整路径
	filename := filepath.Join(gameFileDir, fileName) // 假设上传的图片总是jpg格式
	out, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 将文件内容写入到保存的文件
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sprintf := fmt.Sprintf("/BHdownload/%s", fileName)

	ret.Data = map[string]interface{}{"url": sprintf}

	if manufacturer != "" {
		go func() {
			var gameData *comm.Game2
			coll := db.Collection2("game", "Games")
			one := coll.FindOne(context.TODO(), bson.M{"_id": id})
			err := one.Decode(&gameData)
			if err != nil {
				logger.Err("同步游戏图标失败：err:%v", err)
				return
			}

			for _, item := range gameData.GameNameConfig {
				item.Icon = sprintf
			}

			_, err = coll.UpdateOne(
				context.TODO(),
				bson.M{"_id": id},
				bson.M{"$set": bson.M{"GameNameConfig": gameData.GameNameConfig}})
		}()
	}

	ut.HttpReturnJson(w, ret)
}
