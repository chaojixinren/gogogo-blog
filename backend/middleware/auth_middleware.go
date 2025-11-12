package middleware

import (
	"net/http"
	"strings"

	"gogogo/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawHeader := ctx.GetHeader("Authorization")
		if rawHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			return
		}

		token := rawHeader
		if parts := strings.SplitN(rawHeader, " ", 2); len(parts) == 2 {
			token = parts[1]
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}
