package middleware

import (
	"api-gateway/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// 支持 "Bearer token" 格式
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 开发模式：接受 mock token
		if strings.HasPrefix(tokenString, "mock_jwt_token_") {
			// 从 mock token 中提取用户名
			username := strings.TrimPrefix(tokenString, "mock_jwt_token_")
			c.Set("userID", uint(1)) // 使用固定的用户ID
			c.Set("username", username)
			c.Set("token", tokenString)
			c.Next()
			return
		}

		// 生产模式：验证真实的 JWT token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Global.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok {
			c.Set("userID", claims.UserID)
			c.Set("token", tokenString)
		}

		c.Next()
	}
}
