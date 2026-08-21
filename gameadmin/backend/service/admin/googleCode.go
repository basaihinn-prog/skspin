package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type GoogleAuth struct {
}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

func (auth *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (auth *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (auth *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (auth *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (auth *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (auth *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (auth *GoogleAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := auth.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := auth.toUint32(hashParts)
	return number % 1000000
}

func (auth *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, auth.un())
	return strings.ToUpper(auth.base32encode(auth.hmacSha1(buf.Bytes(), nil)))
}

func (auth *GoogleAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := auth.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := auth.oneTimePassword(secretKey, auth.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}

func (auth *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

func (auth *GoogleAuth) GetQrcodeUrl(user, secret string) string {
	qrcode := auth.GetQrcode(user, secret)
	//	return fmt.Sprintf("https://www.google.com/chart?chs=200x200&chld=M%%7C0&cht=qr&chl=%s", qrcode)
	width := "200"
	height := "200"
	data := url.Values{}
	data.Set("data", qrcode)
	return "https://api.qrserver.com/v1/create-qr-code/?" + data.Encode() + "&size=" + width + "x" + height + "&ecc=M"
}

func (auth *GoogleAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := auth.GetCode(secret)
	if err != nil {
		return false, err
	}
	return _code == code, nil
	//return true,nil
}
