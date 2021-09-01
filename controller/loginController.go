// Package controller
// @Author : Lik
// @Time   : 2020/11/11
package controller

import "github.com/gin-gonic/gin"

type LoginController struct {
}

func (*LoginController) login(rg *gin.RouterGroup) {
	rg.GET("/login", func(c *gin.Context) {
		name, _ := c.Params.Get("name")
		c.JSON(200, gin.H{
			"name": name,
			"age":  12,
		})
	})
}
