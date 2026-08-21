package comm

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

var privateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAtDdVQ3Dt8gTpM8slvceExzGomzssHmImy09oN1U7gbXcZS4t
HUrewoFabz9qjIqu8omVk51D2Nbidg1Gh4ukAAj10NR+edAlLP22AVzmw6qnONhC
T8fRE6OUOokdi9Jy3vfHKpeQYLX/fKxyX+rs0chVuqw3aKBjWe3WJZm6kGS7KG0w
pcs+ToEZ93p0ji8qvOB98O+l9txrOgHGKGZjYfFWrdaxWa03vM6TV1eZzIjH/jbu
lMdYOxxDcrEqgwAkQ3ZcV5xGGELZuXfWkyMe2WoLNaCpAbA1ObMM1Gi5v4uyJLSX
Mkw/N7CEKbFhPte1YHX3lQKo0SRuVCl/Nf9PzQIDAQABAoIBAHbB2+MEYRjSWaay
4R0NhJcLR1N8C9e90Fi77C5CcWNJp4HZiws5kk/Uk/apcJpKrXzQY4wR32reN/+Z
QfgCckE/plVGIk49drIOQsjlIoCgTW/tOs1+HG33pq9oOdsxBFegKlQL6q2AYWsT
7I7+ra6UeMH5yM2em7ngO/UtN778NcxanY6+SbrPKJ6EAZlluRnmyg1TVuc18a0p
sBMqYzQzVwdZpI8CxL1vxhhItzJw58GbiWtj8O9JtvoWQScS8bj7Dr/3jln3hrm9
8vhjxpzcEdrEDnQmNucARtHGeckbZ3z64FyT7Fk8WHzzbOp5+UdRbX57TzXFSBtT
P290vCECgYEA7R11QF+MfZweloI0LJsdtLQVOxxJ8q/H0leNxQP3usiZjI4Nv+Pp
ut2K5rZy48KAVmi4Nnjm8iKHgjr4PnjCVCj6R1CHqyuwo862AYb7Qv87C1e59tjI
FoiMU2x6dCKsMAs+8yu1ekQ3GLsDz2u2/IlcOTSBcEKxFfgp+SlYSjUCgYEAwpHD
HgW3jPufUUIdXIQoqkQAJ39PRFrnU/Tf0nltTIhBjCPS2tkNArl6Y85lR79Ow0/r
c3RgCEbcXCO4jBMJoGXPNUAjiQHpX5xmk6s0HgaE0Q6AwwFPw/2nOegUhSYA607f
ORrdrtb7Sp7mgxI8IJXVEPrqtleLkorLmSsE4jkCgYAb6NW+SADfYBrxmE3P2ko6
1N+S35eMq0gX6BpV0Eu+fpIkSywvJAKE7kLFOUB4spIsmZLlRoHYilvs5kgGAmzN
Py2Ga2Issa3O+ivOLjcxAZ3PjfnjpkyW6meqAiC/vr0JwqkcMk7gH1tk285tAb6+
JuTmDtoVfqQdc+Js44Ly6QKBgEhZDPVv6MWKlr4PWH2bQse1C12kcCQZrSTBzCwm
LKclj0H93By2UqktsL3F9FEOaMolQIaowkCxoKS+P5QOTCkRUlAZrlz2kgGUVWwZ
YAK+J8rYmrZoGXHmMrVMf7zW2caliElinQWzOLORjGM2d5ciP5zVwErXGLX/2B73
KRS5AoGAa7LSQr8krWM1Q3UG2BXc8cSjB/dFQNML5TQty8Q3xqQICHML9XGIk83k
YAd67fjAhlE/QeS4XTOmcatJ4QvFbdpCie/3KbiYPy9StfyjmqBM5C3jwiOZI9sc
2ayWD8V/DNeXr88+5yv/xy+yyieUEiYSF29PwvI8Boj58MI1UkI=
-----END RSA PRIVATE KEY-----`

//func encrypt(data string, publicKeyStr string) ([]byte, error) {
//	// 解析公钥字符串
//	block, _ := pem.Decode([]byte(publicKeyStr))
//	if block == nil {
//		return nil, fmt.Errorf("failed to parse public key")
//	}
//
//	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//
//	// 加密数据
//	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey.(*rsa.PublicKey), []byte(data), nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return encryptedData, nil
//}

//func RsaDecrypt(encryptedData []byte) (string, error) {
//	// 解析私钥字符串
//	block, _ := pem.Decode([]byte(privateKeyStr))
//	if block == nil {
//		return "", fmt.Errorf("failed to parse private key")
//	}
//
//	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		return "", err
//	}
//
//	// 解密数据
//	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedData, nil)
//	if err != nil {
//		return "", err
//	}
//
//	return string(decryptedData), nil
//}

func RsaDecrypt(encryptedBase64 string) (string, error) {
	// 解码Base64
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", err
	}

	// 解析私钥
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("failed to parse PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 解密
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedBytes)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}

// 基本MD5加密（字符串）
func Md5Encrypt(str string) string {
	// 创建MD5哈希对象
	hash := md5.New()

	// 写入字符串
	hash.Write([]byte(str))

	// 计算哈希值并转换为16进制字符串
	return hex.EncodeToString(hash.Sum(nil))
}
