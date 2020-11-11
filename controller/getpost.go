package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func request() {
	gin.DisableConsoleColor()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 路径参数获取
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
	r.POST("/user/:name/*action", func(c *gin.Context) {
		c.FullPath()
	})

	// Query参数
	r.GET("/query", func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		rows := c.DefaultQuery("rows", "10")

		c.String(http.StatusOK, "第 %s 页 每页 %s 条", page, rows)
	})

	r.Run()
}
