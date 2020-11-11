package router

import (
	"github.com/gin-gonic/gin"
)

func InitToDoRouter(Router *gin.RouterGroup) {

	ToDoRouter := Router.Group("todo")
	{
		// 添加代办
		ToDoRouter.POST("/todo", func(c *gin.Context) {
		})

		// 查看所有代办
		ToDoRouter.GET("todo", func(c *gin.Context) {

		})

		// 查看某一个代办事项
		ToDoRouter.GET("todo/:id", func(c *gin.Context) {

		})

		ToDoRouter.GET("todo", func(c *gin.Context) {

		})

		ToDoRouter.GET("todo", func(c *gin.Context) {

		})
	}
}