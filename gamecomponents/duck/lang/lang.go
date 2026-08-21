package lang

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"game/duck/cfgmgr"
	"log/slog"
	"strings"
	"sync"
)

type Lang struct {
	ZH         string            `bson:"_id"`        // 中文
	Permission int32             `bson:"permission"` // 权限
	LangMap    map[string]string ` bson:"inline"`
}

func (l *Lang) MarshalJSON() ([]byte, error) {
	langMap := map[string]any{
		"ZH":         l.ZH,
		"Permission": l.Permission,
	}
	for k, v := range l.LangMap {
		langMap[strings.ToUpper(k)] = v
	}
	data, err := json.Marshal(langMap)
	return data, err
}

func (l *Lang) UnmarshalJSON(data []byte) error {
	langMap := map[string]any{}
	err := json.Unmarshal(data, &langMap)
	if err != nil {
		return err
	}
	if langMap["ZH"] == nil {
		return errors.New("lang unmarshal err, zh")
	}
	l.ZH = langMap["ZH"].(string)
	permission := float64(0)
	if langMap["Permission"] != nil {
		permission = langMap["Permission"].(float64)
	}
	l.Permission = int32(permission)
	delete(langMap, "ZH")
	delete(langMap, "Permission")
	l.LangMap = map[string]string{}
	for k, v := range langMap {
		_, ok := v.(string)
		if !ok {
			continue
		}
		l.LangMap[strings.ToLower(k)] = v.(string)
	}
	return nil
}

var Langs map[string]*Lang
var langMutex sync.Mutex

func Init(name string, cm *cfgmgr.ConfigManager) {
	cm.WatchAndLoad(name, loadLang)
}

func loadLang(buf []byte) error {
	langMutex.Lock()
	defer langMutex.Unlock()

	reader := csv.NewReader(bytes.NewReader(buf))
	all, err := reader.ReadAll()
	if err != nil {
		return err
	}
	if len(all) <= 1 {
		return errors.New("empty")
	}

	header := all[0]
	if len(header) == 0 {
		return errors.New("empty header")
	}
	body := all[1:]
	Langs = map[string]*Lang{}
	for _, row := range body {
		langMap := map[string]string{}
		for j := 1; j < len(row); j++ {
			langMap[header[j]] = row[j]
		}
		_, ok := Langs[row[0]]
		if ok {
			slog.Info("重复:", "key", row[0])
		}
		Langs[row[0]] = &Lang{
			ZH:         row[0],
			Permission: 0,
			LangMap:    langMap,
		}
	}
	return err
}

// 仅获取翻译文本
func Get(language string, id string) string {
	langMutex.Lock()
	defer langMutex.Unlock()
	lang := Langs[id]
	if lang == nil {
		return id
	}
	if language == "zh" {
		return lang.ZH
	}
	if language == "id" {
		language = "idr"
	}
	l, ok := lang.LangMap[language]
	if !ok {
		return id
	}
	return l
}
func GetLang(lang string, id string) string {
	return Get(lang, id)
}

func Error(lang string, id string) error {
	msg := GetLang(lang, id)
	return errors.New(msg)
}

func GetAllArr() []*Lang {
	ret := make([]*Lang, 0, len(Langs))
	for _, v := range Langs {
		ret = append(ret, v)
	}
	return ret
}

// 带参数的翻译
func Translate(language string, id string, data any, plural any) string {
	return Get(language, id)
}
