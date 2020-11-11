package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func TimeMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		timeStart := time.Now()

		c.Next()

		timeElapsed := time.Since(timeStart)
		url := c.Request.URL.String()
		log.Printf("the request URL %s cost %v\n", url, timeElapsed)
	}
}
