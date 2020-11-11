// @Author : Lik
// @Time   : 2020/10/18
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
)

var (
	aesKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`)
	aesIv = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`)
)

// 非对称加密，根据公钥和原始内容产生加密内容
func encryptRSA(key []byte, plainTest []byte) []byte {
	block, _ := pem.Decode(key)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	pubKey := pubInterface.(*rsa.PublicKey)

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plainTest)
	if err != nil {
		panic(err)
	}
	return cipherText
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//AesEncrypt 加密函数
func AesEncrypt(key, iv []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted := make([]byte, aes.BlockSize+len(plaintext))
	iiv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iiv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iiv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], plaintext)

	return encrypted, nil

}

//明文补码算法
//func PKCS5Padding(src []byte, blockSize int) []byte {
//	padding := blockSize - len(src)%blockSize
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(src, padtext...)
//}

//明文减码算法
//func PKCS5UnPadding(src []byte) []byte {
//	length := len(src)
//	unpadding := int(src[length-1])
//	return src[:(length - unpadding)]
//}
// AesDecrypt 解密函数
func AesDecrypt(key, iv []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func encryptSendData(publicKey []byte, privateKey []byte, sendData []byte) ([]byte, error) {

	//aesKey, _ := hex.DecodeString("0f90023fc9ae101e")
	//aesIv, _ := hex.DecodeString("0f90023fc9ae101d")
	//aesKey := []byte("pMDcg8YAONlvzxGZ")
	//aesIv := []byte("CTycf3Z6iVjtnxLL")

	aesKeyWithBase64 := Krand(32, 3)
	aesIVWithBase64 := RangeRand(1000000000000000, 10000000000000000)

	//key, _ := RSAEncrypt(aesKeyWithBase64, PublicKey)
	//salt, _ := RSAEncrypt([]byte(aesIVWithBase64), PublicKey)
	// 1DT3ioauTu3v6A1c3Dt7090H3LblT6Rf
	key, _ := RSAEncrypt(aesKeyWithBase64, PublicKey)
	// 4188184627492260
	salt, _ := RSAEncrypt([]byte(aesIVWithBase64), PublicKey)

	//key := encryptRSA(publicKey, aesKey)
	//salt := encryptRSA(publicKey, aesIv)

	ciphertext, err := AesEncrypt(aesKeyWithBase64, []byte(aesIVWithBase64), sendData)

	if err != nil {
		panic(err)
	}

	requestMap := make(map[string][]byte)
	requestMap["key"] = []byte(key)
	requestMap["salt"] = []byte(salt)
	requestMap["data"] = ciphertext
	requestMap["source"] = []byte("TEST")

	sign := sha256Map(requestMap)
	signb, _ := RSAEncrypt(sign, string(privateKey))
	requestMap["sign"] = []byte(signb)

	return json.Marshal(requestMap)
}

func encryptByResult(data []byte) {

	sendData, _ := encryptSendData(aesKey, aesIv, data)
	fmt.Println(sendData)
}

func sha256Map(paramMap map[string][]byte) []byte {
	concatStr := ""
	for k, v := range paramMap {
		if "sign" == k {
			continue
		}
		concatStr += k + string(v) + "&"
	}
	h := sha256.New()
	h.Write([]byte(concatStr))
	sum := h.Sum(nil)

	//由于是十六进制表示，因此需要转换
	s := hex.EncodeToString(sum)
	return []byte(s)
}

func mainss() {
	encryptByResult([]byte("我爱你"))
}
