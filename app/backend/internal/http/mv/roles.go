package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	// 1) Строим "множество" разрешённых ролей:
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{} // O(1) проверка принадлежности
	}

	// 2) Возвращаем функцию-мидлварь (замыкание),
	// которая будет выполнена ПЕРЕД хэндлером
	return func(c *gin.Context) {
		// роль заранее должна быть положена в контекст (например, auth-мидлварью)
		role, _ := c.Get("role")

		// пытаемся привести к string
		if roleStr, ok := role.(string); ok {
			// если роль есть в allowed → пропускаем дальше по цепочке
			if _, ok := allowed[roleStr]; ok {
				c.Next() // передаём управление следующему middleware/хэндлеру
				return
			}
		}

		// иначе — запрещаем доступ и прерываем цепочку
		c.AbortWithStatus(http.StatusForbidden) // 403
	}
}
