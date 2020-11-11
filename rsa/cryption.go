package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	r "math/rand"
	"time"
)

// RSAEncrypt RSA加密
// plainText 要加密的数据
// publicKey 公钥匙内容
func RSAEncrypt(plainText []byte, publicKey string) (string,error) {
	key, _ := base64.StdEncoding.DecodeString(publicKey)
	pubKey, _ := x509.ParsePKIXPublicKey(key)
	//logx.Infof("%v",pubKey)
	//解密pem格式的公钥
	//block, _ := pem.Decode([]byte(publicKey))
	//if block == nil {
	//	return "", fmt.Errorf("public key error")
	//}
	//// 解析公钥
	//pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	//if err != nil {
	//	return "", err
	//}
	// 类型断言
	pub := pubKey.(*rsa.PublicKey)
	encryptedData,err:=rsa.EncryptPKCS1v15(rand.Reader, pub, plainText)
	return base64.StdEncoding.EncodeToString(encryptedData),err
}

// RSADecrypt RSA解密
// cipherText 需要解密的byte数据
// privateKey 私钥匙内容
func RSADecrypt(cipherText ,privateKey string) (string,error){
	encryptedDecodeBytes,err:=base64.StdEncoding.DecodeString(cipherText)
	if err!=nil {
		return "",err
	}
	key,_:=base64.StdEncoding.DecodeString(privateKey)
	prvKey,_:=x509.ParsePKCS1PrivateKey(key)
	originalData,err:=rsa.DecryptPKCS1v15(rand.Reader,prvKey,encryptedDecodeBytes)
	return string(originalData),err
}

func RSAPriEncrypt(cipherText ,privateKey string)(string,error) {
	key,_:=base64.StdEncoding.DecodeString(privateKey)
	prvKey,_:=x509.ParsePKCS1PrivateKey(key)
	rng := rand.Reader
	hashed := sha256.Sum256([]byte(cipherText))
	signature, err := rsa.SignPKCS1v15(rng, prvKey, crypto.SHA256, hashed[:])
	if err != nil {
		//logx.Errorf("Error from signing: %s\n", err)
		return "",err
	}

	return fmt.Sprintf("%x", signature), nil
}

//RangeRand 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) string {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return fmt.Sprintf("%d", result.Int64() - i64Min)
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return fmt.Sprintf("%d", min + result.Int64())
	}
}


// Krand 随机字符串
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	r.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = r.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + r.Intn(scope))
	}
	return result
}


//AESEncrypt AES加密
func AESEncrypt(origData, key,iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCFBEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.XORKeyStream(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}


func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}


/**
 * 计算sha256值
 *
 * @param paramMap
 * @return 签名后的所有数据，原始数据+签名
 */

func Sha256(requestMap map[string]string)string{
	str := ""
	for k,v := range requestMap {
		str += fmt.Sprintf("%s=%s&",k,v)
	}
	//logx.Infof("requestMap %s",str)
	sum := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sum)
}

//func main()  {
//
//}