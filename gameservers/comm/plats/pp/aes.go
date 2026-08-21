package pp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, paddingText...)
}

// 去除 PKCS7 填充
func pkcs7UnPadding(data []byte) []byte {
	padding := data[len(data)-1]
	return data[:len(data)-int(padding)]
}

// AES 加密函数（使用 ECB 模式和 PKCS7 填充）
func Encrypt(plaintext, key []byte) (string, error) {
	// 创建一个 AES 块加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充明文，使其成为块大小的倍数
	blockSize := block.BlockSize()
	plaintext = pkcs7Padding(plaintext, blockSize)

	// 创建加密密文的缓冲区
	ciphertext := make([]byte, len(plaintext))

	// 使用 ECB 模式进行加密
	mode := NewECBEncrypter(block)
	mode.CryptBlocks(ciphertext, plaintext)

	// 将加密后的字节数据转换为 Base64 格式
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

// AES 解密函数（使用 ECB 模式和 PKCS7 填充）
func Decrypt(ciphertextBase64, key []byte) ([]byte, error) {
	// 创建一个 AES 块加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 解码 Base64 密文
	ciphertext, err := base64.StdEncoding.DecodeString(string(ciphertextBase64))
	if err != nil {
		return nil, err
	}

	// 创建解密明文的缓冲区
	plaintext := make([]byte, len(ciphertext))

	// 使用 ECB 模式进行解密
	mode := NewECBDecrypter(block)
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	plaintext = pkcs7UnPadding(plaintext)
	return plaintext, nil
}

// 创建自定义的 ECB 加密器
func NewECBEncrypter(block cipher.Block) cipher.BlockMode {
	return &ecbEncrypter{block}
}

// 创建自定义的 ECB 解密器
func NewECBDecrypter(block cipher.Block) cipher.BlockMode {
	return &ecbDecrypter{block}
}

// 实现 ECB 加密器
type ecbEncrypter struct {
	block cipher.Block
}

func (e *ecbEncrypter) BlockSize() int {
	return e.block.BlockSize()
}

func (e *ecbEncrypter) CryptBlocks(dst, src []byte) {
	blockSize := e.block.BlockSize()
	if len(src)%blockSize != 0 {
		panic("src length is not a multiple of block size")
	}
	for i := 0; i < len(src); i += blockSize {
		e.block.Encrypt(dst[i:i+blockSize], src[i:i+blockSize])
	}
}

// 实现 ECB 解密器
type ecbDecrypter struct {
	block cipher.Block
}

func (e *ecbDecrypter) BlockSize() int {
	return e.block.BlockSize()
}

func (e *ecbDecrypter) CryptBlocks(dst, src []byte) {
	blockSize := e.block.BlockSize()
	if len(src)%blockSize != 0 {
		panic("src length is not a multiple of block size")
	}
	for i := 0; i < len(src); i += blockSize {
		e.block.Decrypt(dst[i:i+blockSize], src[i:i+blockSize])
	}
}
