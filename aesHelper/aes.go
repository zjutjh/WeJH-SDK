package aesHelper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var encryptKey []byte

var initKey = false

var (
	// ErrKeyNotInit 未初始化密钥
	ErrKeyNotInit = errors.New("请先初始化 AES 密钥")

	// ErrKeyLengthNotMatch 密钥长度不匹配
	ErrKeyLengthNotMatch = errors.New("AES 密钥长度必须为 16、24 或 32 字节")
)

// Init 读入 AES 密钥配置
func Init(key string) error {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return ErrKeyLengthNotMatch
	}
	encryptKey = []byte(key)
	initKey = true
	return nil
}

// Encrypt AES 加密
func Encrypt(orig string) (string, error) {
	if !initKey {
		return "", ErrKeyNotInit
	}

	origData := []byte(orig)

	// 分组秘钥
	block, err := aes.NewCipher(encryptKey)
	if err != nil {
		return "", err
	}

	// 进行 PKCS7 填充
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)

	// 使用 CBC 加密模式
	blockMode := cipher.NewCBCEncrypter(block, encryptKey[:blockSize])
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)

	// 使用 RawURLEncoding 编码为 Base64，适合放入 URL
	return base64.RawURLEncoding.EncodeToString(cryted), nil
}

// Decrypt AES 解密
func Decrypt(cryted string) (string, error) {
	if !initKey {
		return "", ErrKeyNotInit
	}

	// 解码 Base64 字符串
	crytedByte, err := base64.RawURLEncoding.DecodeString(cryted)
	if err != nil {
		return "", err
	}

	// 分组秘钥
	block, err := aes.NewCipher(encryptKey)
	if err != nil {
		return "", err
	}

	// CBC 模式解密
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, encryptKey[:blockSize])
	orig := make([]byte, len(crytedByte))
	blockMode.CryptBlocks(orig, crytedByte)

	// 去除 PKCS7 填充
	orig = PKCS7UnPadding(orig)

	return string(orig), nil
}

// PKCS7Padding 填充数据，使长度为 blockSize 的倍数
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7UnPadding 去除填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
