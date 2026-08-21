package jiliut

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"serve/comm/mq"
	"serve/comm/ut"
	"strconv"

	"serve/comm/jwtutil"
	"serve/comm/mux"

	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ProtoReader(m proto.Message) io.Reader {
	data := lo.Must(proto.Marshal(m))
	return bytes.NewReader(data)
}

func ProtoEncode(m proto.Message) []byte {
	data := lo.Must(proto.Marshal(m))
	return data
}

// 填充数据到16字节的倍数
func pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// AES-CBC 加密
func ProtoEncryption(token string, m proto.Message) ([]byte, error) {
	// 生成密钥
	key := []byte(token[:32])                    // 确保密钥长度为32字节
	dataBefore16 := ut.GenerateRandomString2(16) // 生成IV
	iv := []byte(dataBefore16)

	// 序列化 Protobuf 消息
	plainText, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("序列化错误: %v", err)
	}

	// 填充明文
	plainText = pad(plainText)

	// 创建 AES 加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建加密器错误: %v", err)
	}

	// 创建 CBC 模式的加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	// 返回 IV 和密文
	return append(iv, cipherText...), nil
}

func PostProto(ul string, req proto.Message, res protoreflect.ProtoMessage) (err error) {
	resp, err := http.Post(ul, "application/x-protobuf", ProtoReader(req))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if res != nil {
		err = proto.Unmarshal(body, res)
	}
	return
}
func PostProtoWithHeaders(ul string, ps proto.Message, res protoreflect.ProtoMessage, headers map[string]string) (err error) {
	req := lo.Must(http.NewRequest("POST", ul, ProtoReader(ps)))
	if len(headers) != 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/x-protobuf")

	// resp, err := http.Post(ul, "application/x-protobuf", ProtoReader(ps))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if res != nil {
		err = proto.Unmarshal(body, res)
	}
	return
}

func PostJson(ul string, ps any, res any) (err error) {
	var psReader io.Reader
	if ps != nil {
		data := lo.Must(json.Marshal(ps))
		psReader = bytes.NewReader(data)
	}
	resp, err := http.Post(ul, "application/json", psReader)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if res != nil {
		err = json.Unmarshal(body, res)
	}
	return
}

func GetBody(ul string) (body []byte, err error) {
	resp, err := http.Get(ul)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	body, err = io.ReadAll(resp.Body)

	return
}

func GetJsonWithHeaders(ul string, res any, headers map[string]string) (err error) {
	req := lo.Must(http.NewRequest("GET", ul, nil))
	if len(headers) != 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if res != nil {
		err = json.Unmarshal(body, res)
	}
	return
}

func PostJsonWithHeaders(ul string, ps any, res any, headers map[string]string) (err error) {
	var psReader io.Reader
	if ps != nil {
		data := lo.Must(json.Marshal(ps))
		psReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest("POST", ul, psReader)
	if len(headers) != 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	// resp, err := http.Post(ul, "application/json", psReader)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if res != nil {
		err = json.Unmarshal(body, res)
	}
	return
}

func LaunchGame(username, game, lang string) (url string, err error) {
	var ret struct {
		Url string
	}

	err = mux.HttpInvoke("https://plats.rpgamestest.com/plat/JILI/LaunchGame", map[string]string{
		"UID":  username,
		"Game": game,
		"Lang": lang,
	}, &ret)
	if err != nil {
		return
	}

	url = ret.Url
	return
}

type ReturnObj struct {
	Code    int
	Message string
	Data    any
}

func MarshalJsonReturn(data any) []byte {
	return lo.Must(json.Marshal(ReturnObj{
		Code:    0,
		Message: "成功",
		Data:    data,
	}))
}

func ParseAuthToken(token string) (pid int64, err error) {
	o := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, &o, func(t *jwt.Token) (interface{}, error) {
		return jwtutil.SignatureKey(), nil
	})
	if err != nil {
		return
	}

	aidstr, ok := o["AID"].(string)
	if !ok {
		err = errors.New("not string type")
		return
	}
	n, err := strconv.Atoi(aidstr)

	pid = int64(n)
	return
}

func GetFetchMongoAddr() string {
	return "mongodb://myUserAdmin:aajohsujie9vecieSohqu4weivai7oxayei9rie5Yoh4vuojohwaothee0waethi@127.0.0.1:27017/?authSource=admin"
}

type GetOperatorInfoRes struct {
	//商户币种        AdminOperator_copy1 / CurrencyKey
	CurrencyKey                   string         `json:"CurrencyKey" bson:"CurrencyKey"`
	CurrencyManufactureVisibleOff map[string]int `json:"CurrencyManufactureVisibleOff" bson:"CurrencyVisibleOff"`
}

func GetOperatorInfo(appId string) GetOperatorInfoRes {
	resInfo := GetOperatorInfoRes{}
	err := mq.Invoke("/AdminInfo/Interior/operatorInfo", map[string]any{
		"AppID": appId,
	}, &resInfo)

	if err != nil {
		return resInfo
	}
	return resInfo
}
