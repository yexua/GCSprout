// @Author : Lik
// @Time   : 2021/1/26
package zap

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

/**
 * @Description: gin 结合zap的日志中间件
 * @param logger zap Logger 对象
 * @param timeFormat 时间格式字符串
 * @param utc 是否为utc时区
 * @return gin.HandlerFunc gin中间件
 */
func GinZapMiddleware(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(timeFormat)),
				zap.Duration("latency", latency),
			)
		}
	}
}

//
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
/**
 * @Description: that recovers from any panics and logs requests using uber-go/zap.
 * @param logger
 * @param stack 是否输出堆栈消息
 * @return gin.HandlerFunc returns a gin.HandlerFunc (middleware)
 */
func RecoveryWithZap(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errCheck
					c.Abort()
					return
				}
				//var fieldList []zap.Field
				fieldList := make([]zap.Field, 4, 5)
				//fieldList = append(fieldList, zap.Time("time", time.Now()))
				//fieldList = append(fieldList, zap.Any("error", err))
				//fieldList = append(fieldList, zap.Time("time", time.Now()))
				//fieldList = append(fieldList, zap.Time("time", time.Now()))
				fieldList[0] = zap.Time("time", time.Now())
				fieldList[1] = zap.Any("error", err)
				fieldList[2] = zap.Time("time", time.Now())
				fieldList[3] = zap.Time("time", time.Now())
				if stack {
					//fieldList[4] = zap.String("stack", string(debug.Stack()))
					fieldList = append(fieldList, zap.String("stack", string(debug.Stack())))
				}
				logger.Error("[Recovery from panic]", fieldList...)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
