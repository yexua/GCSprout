package main

import (
	"GCSprout/zap"
	"fmt"
	"github.com/gin-gonic/gin"
	uberZap "go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	gin.DisableConsoleColor()
	r := gin.New()
	logger, _ := uberZap.NewProduction()
	defer logger.Sync()
	r.Use(zap.GinZapMiddleware(logger, "2006-01-02 15:04:05", false))
	//r.Use(zap.RecoveryWithZap(logger, true))
	r.GET("/ping", func(c *gin.Context) {

		time.Sleep(time.Second * 2)
		logger.Info("业务执行成功")

		c.JSON(http.StatusOK, gin.H{
			"res": "Hello World",
		})
	})

	r.GET("/dlv", TestDlvHandler)

	r.GET("/panic", func(c *gin.Context) {
		panic("This is a panic!")
	})

	r.Run(":9090")
}

func TestDlvHandler(c *gin.Context) {
	fmt.Println("接受到请求")
	a := 1
	b := a + 1
	d := b + 2
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(d)

	c.JSON(http.StatusOK, gin.H{
		"res": "Successes!",
	})
}
