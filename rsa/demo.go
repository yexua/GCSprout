// @Author : Lik
// @Time   : 2020/10/19
package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"
)

type ZyBaseRequest struct {
	AppId     string      `json:"appId"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
	Sign      string      `json:"sign"`
	Source   string
}



const (
	PushFeatureDeal = "/ndToZy/pushFeatureDeal"
	PushUserInfo    = "/ndToZy/pushUserInfo"
	PushFeatureInfo = "/ndToZy/pushFeatureInfo"

	//PublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwTkJpYvb+n96f/heVUmdLX1dt+l8V9/NOkbFxG4nSdZmcAaAmrlE1/+Kw2i3i+LaqQTw8a80TTVCTl3w83mwxzwExtQkB1KkxdPNTxOINVKuurE4fD6tzi3ooRC/7QLBb1Q+NSsDOdheSkLlC67URmKtW2IMwX/uCwQ0mpPBI0IcAVF42xMDx4PK9PkwZQe16d8x9Aa+rVpn8AoFGSqT0OoSm5Z20QrUs6tRyWV+B0JGEzc1Mg1oCu6880nCCMlAfzHXC9QXRWnFjJthF99NQNZNww7pKYlCDhYe5G/nl1aWo/o3e5zVu1uo2vF/bT+/hXaVNcSXkM56Z057fCJMNQIDAQAB"
	PublicKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDI6d306Q8fIfCOaTXyiUeJHkrIvYISRcc73s3vF1ZT7XN8RNPwJxo8pWaJMmvyTn9N4HQ632qJBVHf8sxHi/fEsraprwCtzvzQETrNRwVxLO5jVmRGi60j8Ue1efIlzPXV9je9mkjzOmdssymZkh2QhUrCmZYI/FCEa3/cNMW0QIDAQAB"
	PrivateKey = "MIICXQIBAAKBgQC+L0rfjLl3neHleNMOsYTW8r0QXZ5RVb2p/vvY3fJNNugvJ7lo4+fdBz+LN4mDxTz4MTOhi5e2yeAqx+v3nKpNmPzC5LmDjhHZURhwbqFtIpZD51mOfno2c3MDwlrsVi6mTypbNu4uaQzw/TOpwufSLWF7k6p2pLoVmmqJzQiD0QIDAQABAoGAakB1risquv9D4zX7hCv9MTFwGyKSfpJOYhkIjwKAik7wrNeeqFEbisqv35FpjGq3Q1oJpGkem4pxaLVEyZOHONefZ9MGVChT/MNH5b0FJYWl392RZy8KCdq376Vt4gKVlABvaV1DkapL+nLh7LMo/bENudARsxD55IGObMU19lkCQQDwHmzWPMHfc3kdY6AqiLrOss+MVIAhQqZOHhDe0aW2gZtwiWeYK1wB/fRxJ5esk1sScOWgzvCN/oGJLhU3kipHAkEAysNoSdG2oWADxlIt4W9kUiiiqNgimHGMHPwp4JMxupHMTm7D9XtGUIiDijZxunHv3kvktNfWj3Yji0661zHVJwJBAM8TDf077F4NsVc9AXVs8N0sq3xzqwQD/HPFzfq6hdR8tVY5yRMb4X7+SX4EDPORKKsgnYcur5lk8MUi7r072iUCQQC8xQvUne+fcdpRyrR4StJlQvucogwjTKMbYRBDygXkIlTJOIorgudFlrKP/HwJDoY4uQNl8gQJb/1LdrKwIe7FAkBl0TNtfodGrDXBHwBgtN/t3pyi+sz7OpJdUklKE7zMSBuLd1E3O4JMzvWP9wEE7JDb+brjgK4/cxxUHUTkk592"

)

var configFile = flag.String("f", "config/job_risk_third.json", "the config file")

type FeatureDealRequest struct {
	FeatureId  string `json:"featureId"`
	UserId     string `json:"userId"`
	ObjectId   string `json:"objectId"`
	DealId     string `json:"dealId"`
	DealName   string `json:"dealName"`
	DealTime   string `json:"dealTime"`
	DealStatus string `json:"dealStatus"`
	DealResult string `json:"dealResult"`
}
type BodyStruct struct {
	Data string `json:"data"`
	Sign string `json:"sign"`
	Key string `json:"key"`
	Salt string `json:"salt"`
}

func main() {
	flag.Parse()

	body := new(FeatureDealRequest)
	body.FeatureId = "111111"
	body.UserId = "22222"
	body.ObjectId = "33333"
	body.DealId = "44444"
	body.DealName = "李四"
	body.DealTime = time.Now().Format("2006-01-02 15:04:05")
	body.DealStatus = "1"
	body.DealResult = "已处置，属于正常工作情况"
	jsBody, _ := json.Marshal(body)
	aesKeyWithBase64 := Krand(32, 3)
	aesIVWithBase64 := RangeRand(1000000000000000, 10000000000000000)
	//aesIVWithBase64 := []byte("m4mVWmNAZqfXfV5WxyI2WA")
	//aesKeyWithBase64 := []byte("b7xY4tV91fLhs0fRm1")
	key, _ := RSAEncrypt(aesKeyWithBase64, PublicKey)
	salt, _ := RSAEncrypt([]byte(aesIVWithBase64), PublicKey)
	cipherData, _ := AESEncrypt(jsBody, aesKeyWithBase64, []byte(aesIVWithBase64))
	reqMap := make(map[string]string)
	reqMap["key"] = key
	reqMap["salt"] = salt
	reqMap["data"] = cipherData
	reqMap["source"] = "TEST"
	sign, _ := RSAPriEncrypt(Sha256(reqMap), PrivateKey)
	reqMap["sign"] = sign
	//url := fmt.Sprintf("%s:%d%s", cfg.GetString("zy_host.address"), cfg.GetInt("zy_host.port"), PushFeatureDeal)
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/x-www-form-urlencoded"
	log.Printf("aesKeyWithBase64 %s", string(aesKeyWithBase64))
	log.Printf("aesIVWithBase64 %s", string(aesIVWithBase64))
	log.Printf("key %s", string(key))
	log.Printf("salt %s", string(salt))
	log.Printf("cipherData %s", string(cipherData))
	log.Printf("sign %s", sign)
	b,_ := json.Marshal(BodyStruct{
		Data: cipherData,
		Sign: sign,
		Key:  key,
		Salt: salt,

	})
	log.Println(string(b))
	//log.Printf("url %s", url)
	//b, err := utils.HttpDo(url,"POST",reqMap,headerMap)
	//if err!= nil{
	//  fmt.Println(err)
	//}
	//fmt.Println(string(b))
}
