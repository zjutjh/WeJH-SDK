package oauth

import (
	"crypto/rsa"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type publicKeyData struct {
	Modulus  string `json:"modulus"`
	Exponent string `json:"exponent"`
}

// getPublicKey 获取加密密钥
func getPublicKey(client *resty.Client) (publicKey *rsa.PublicKey, err error) {
	var data publicKeyData
	_, err = client.R().
		SetResult(&data).
		Get(PublicKeyUrl)
	if err != nil {
		return nil, err
	}
	modulus := new(big.Int)
	modulus.SetString(data.Modulus, 16)

	e, err := strconv.ParseInt(data.Exponent, 16, 32)
	if err != nil {
		return nil, err
	}
	exponent := int(e)
	return &rsa.PublicKey{
		N: modulus,
		E: exponent,
	}, nil
}

// rsaEncrypt 加密
func rsaEncrypt(publicKey *rsa.PublicKey, text []byte) []byte {
	chunkSize := 2 * (publicKey.N.BitLen()/16 - 1)
	textLen := len(text)
	result := make([][]byte, 0)
	// 分块加密
	for i := textLen; i > 0; i -= chunkSize {
		textChunk := new(big.Int)
		textChunk.SetBytes(text[max(i-chunkSize, 0):i])
		textChunk.Exp(textChunk, big.NewInt(int64(publicKey.E)), publicKey.N)
		result = append(result, textChunk.Bytes())
	}

	return result[0]
}

// getEncryptedPassword 密码加密
func getEncryptedPassword(client *resty.Client, password string) (string, error) {
	key, err := getPublicKey(client)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rsaEncrypt(key, []byte(password))), nil
}
