package jwt

import (
	"github.com/Samoy/bill_backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var Username string

func Jwt() gin.HandlerFunc {
	return func(context *gin.Context) {
		var message string
		token := context.GetHeader("token")
		claims, err := utils.ParseToken(token)
		if err != nil {
			message = "token验证失败"
		} else if time.Now().Unix() > claims.ExpiresAt {
			message = "token验证超时"
		}
		if len(message) != 0 {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": message,
			})
			context.Abort()
			return
		}
		if claims != nil {
			Username = claims.Username
		}
		context.Next()
	}
}
