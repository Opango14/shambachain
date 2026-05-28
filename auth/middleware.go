package auth

import (
	"net/http"
	"shambachain/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}
		part := strings.SplitN(authHeader, " ", 2)
		if len(part) != 2 || strings.ToLower(part[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			return
		}
		tokenString := strings.TrimSpace(part[1])

		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token missing",
			})
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token: " + err.Error(),
			})
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Set("username", claims["username"])
		ctx.Next()
	}
}
