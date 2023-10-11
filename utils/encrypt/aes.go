package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

type AesEncrypt struct {
	error   error
	secret  []byte
	iv      []byte
	outType OutType
}

type AesEncode struct {
	error     error
	outType   OutType
	encrypted []byte
}
type AesDecode struct {
	error     error
	outType   OutType
	decrypted []byte
}

type OutType int

const (
	OutBase64 OutType = iota //base64
	OutHex                   //16进制
)

func NewAes(secret string, iv string, outType OutType) *AesEncrypt {
	return &AesEncrypt{
		secret:  []byte(secret),
		iv:      []byte(iv),
		outType: outType,
	}
}

func (e *AesEncrypt) OutType(outType OutType) *AesEncrypt {
	e.outType = outType
	return e
}

func (e *AesEncrypt) Encode(origData []byte) *AesEncode {
	aesCipher, _ := aes.NewCipher(e.secret)
	blockSize := aesCipher.BlockSize()
	origData = e.pkcs7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(aesCipher, e.iv)
	// 创建数组
	encrypted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(encrypted, origData)
	return &AesEncode{
		outType:   e.outType,
		encrypted: encrypted,
		error:     e.error,
	}
}

func (e *AesEncode) Error() error {
	return e.error
}

func (e *AesEncode) String() string {
	if e.outType == OutBase64 {
		return base64.StdEncoding.EncodeToString(e.encrypted)
	}
	return hex.EncodeToString(e.encrypted)
}

func (e *AesEncrypt) Decode(encryptedStr string) *AesDecode {
	result := &AesDecode{
		outType:   e.outType,
		decrypted: []byte(encryptedStr),
	}
	encrypted, err := e.stringToEncrypted(encryptedStr)
	if err != nil {
		result.error = err
		return result
	}
	// 分组秘钥
	aesCipher, _ := aes.NewCipher(e.secret)
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(aesCipher, e.iv)
	// 排除未加密参数
	if len(encrypted)%blockMode.BlockSize() != 0 {
		result.error = errors.New("crypto/cipher: input not full blocks")
		return result
	}
	// 创建数组
	orig := make([]byte, len(encrypted))
	// 解密
	blockMode.CryptBlocks(orig, encrypted)
	// 去补全码
	decrypted := e.pkcs7UnPadding(orig)
	return &AesDecode{
		outType:   e.outType,
		decrypted: decrypted,
		error:     e.error,
	}
}

func (e *AesDecode) Error() error {
	return e.error
}
func (e *AesDecode) Data() string {
	return string(e.decrypted)
}

// stringToEncrypted base64或hex 解码
func (e *AesEncrypt) stringToEncrypted(encryptedStr string) ([]byte, error) {
	if e.outType == OutBase64 {
		return base64.StdEncoding.DecodeString(encryptedStr)
	}
	return hex.DecodeString(encryptedStr)
}

//补码
func (e *AesEncrypt) pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func (e *AesEncrypt) pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
