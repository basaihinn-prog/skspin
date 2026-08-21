package pg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// Aes/ECB模式的加密方法，PKCS7填充方式
// func AesEcbEncrypt(src, key string) ([]byte, error) {
// 	block, err := aes.NewCipher([]byte(key))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if src == "" {
// 		return nil, errors.New("plain content empty")
// 	}
// 	ecb := NewECBEncrypter(block)
// 	content := []byte(src)
// 	content = PKCS5Padding(content, block.BlockSize())
// 	crypted := make([]byte, len(content))
// 	ecb.CryptBlocks(crypted, content)

// 	return crypted, nil
// }

// Aes/ECB模式的解密方法，PKCS7填充方式
// func AesEcbDecrypt(crypted, key []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	blockMode := NewECBDecrypter(block)
// 	origData := make([]byte, len(crypted))
// 	blockMode.CryptBlocks(origData, crypted)
// 	origData = PKCS5UnPadding(origData)

// 	return origData, nil
// }

func AesCBCDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	if len(ciphertext) < blockSize {
		return nil, errors.New("ciphertext too short")
	}
	// iv := []byte(aesIvDefValue)
	if len(ciphertext)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	blockModel := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plaintext, ciphertext)
	plaintext = PKCS5UnPadding(plaintext)
	return plaintext, nil
}

func AesCBCEncrypt(plaintext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		// fmt.Println("err=", err)
		return nil, errors.New("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	plaintext = PKCS5Padding(plaintext, blockSize)
	// iv := []byte(aesIvDefValue)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func AesEncrypt(orig string, key string) (string, error) {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}

	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS5Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return hex.EncodeToString(cryted), nil
}

func AesDecrypt(cryted, key string) (string, error) {

	crytedByte, err := hex.DecodeString(cryted)
	if err != nil {
		return "", err
	}

	// 转成字节数组
	k := []byte(key)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS5UnPadding(orig)
	return string(orig), nil
}

const pg2_enc_k = "you can't guess me!!!@2021-10-11" // len = 32

func GenToken(uid string) (string, error) {
	return uid, nil
	// tk := fmt.Sprintf("%v", uid)

	// enc, err := AesEncrypt(uid, pg2_enc_k)
	// if err != nil {
	// 	return "", err
	// }

	// return enc, nil
}
