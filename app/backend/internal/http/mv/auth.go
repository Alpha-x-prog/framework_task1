package mw

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"example/defects/app/backend/internal/auth"
)

// AuthRequired проверяет заголовок Authorization: Bearer <token>
// Если валидный — кладёт uid (int64) и role (string) в контекст.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.Parse(tokenStr, secret)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// sub в MapClaims обычно float64 (из-за JSON-чисел) — аккуратно приводим к int64
		var uid int64
		if v, ok := claims["sub"]; ok {
			switch t := v.(type) {
			case float64:
				uid = int64(t)
			case int64:
				uid = t
			}
		}
		var role string
		if v, ok := claims["role"]; ok {
			if s, ok := v.(string); ok {
				role = s
			}
		}

		if uid == 0 || role == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("uid", uid)
		c.Set("role", role)
		c.Next()
	}
}
