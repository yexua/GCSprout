package main

import (
	"GCSprout/db"
	"GCSprout/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminLoginInput struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func main() {
	r := gin.Default()
	r.Use(middleware.TimeMiddleware())

	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*")

	r.POST("/test", func(c *gin.Context) {
		admin := &AdminLoginInput{}
		c.ShouldBind(admin)
		fmt.Printf("%#v",admin)

	})


	r.POST("/auth", middleware.AuthHandler)

	r.GET("/home", middleware.JWTAuthMiddleware(), middleware.HomeHandler)


	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})



	r.Run()

	db.InitMySQL()
	db := db.DB.DB()
	defer db.Close()

}
