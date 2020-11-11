package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("NingDaKeJi")

// 生成jwt
func GenToken(username string) (string, error) {
	c := MyClaims{
		username,
		jwt.StandardClaims{
			//过期时间
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			//签发人
			Issuer: "NingDa",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}


//解析jwt
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString,&MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil{
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid{
		return claims, nil
	}
	return nil, errors.New("invalid token")

}

func AuthHandler(c *gin.Context) {
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil{
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}
	// 校验用户名和密码是否正确
	msg  := ""
	if user.Username == "lik" && user.Password == "yik"{
		//生成token
		if tokenString, err := GenToken(user.Username); err == nil{
			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"msg":  "success",
				"data": gin.H{"token": tokenString},
			})
			return
		}else {
			log.Println(err)
			msg = "服务器出错"
		}
	}else{
		msg = "鉴权失败"
	}
	c.JSON(http.StatusOK,gin.H{
		"code": 2002,
		"msg":  msg,
	})
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带token的三种方式 1放在请求头 2.放在请求体， 放在URI
		// 这里假设token放在header的authorization中，并使用bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == ""{
			c.JSON(http.StatusOK, gin.H{
				"code":2003,
				"msg":"请求头中auth为空",
			})
			c.Abort()
			return
		}

		// 按照空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 20004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}

		// 将当前请求的username信息baocun到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next()

	}
}

func HomeHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"username": username},
	})
}