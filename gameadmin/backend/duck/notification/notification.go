package notification

import (
	"encoding/json"
	"errors"
)

type Notification struct {
	NotifyTitle     string `bson:"notifyTitle" json:"notify_title"`
	NotifyType      string `bson:"notifyType" json:"notify_type"`
	OnlineLimitTime string `bson:"onlineLimitTime" json:"online_limit_time"`
	LowerLimitTime  string `bson:"lowerLimitTime" json:"lower_limit_time"`
	Status          string `bson:"notifyStatus" json:"status"`
	Operator        int64  `bson:"operator"`
	CreatedTime     string `bson:"created_time"`
	Sort            int32  `bson:"sort"`
}

func MarshalJSON(languageConfig string) (languageConfigMap map[string]string, err error) {
	marshal, err := json.Marshal(languageConfig)
	if err != nil {
		return nil, err
	}
	languageConfigMap, err = UnmarshalJSON(marshal)
	if err != nil {
		return nil, err
	}

	return languageConfigMap, nil
}
func UnmarshalJSON(languageConfig []byte) (languageConfigMap map[string]string, err error) {

	err = json.Unmarshal(languageConfig, &languageConfigMap)

	if err != nil {
		return nil, errors.New("语言配置获取失败")
	}

	return languageConfigMap, nil
}
