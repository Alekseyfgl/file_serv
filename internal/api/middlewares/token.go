package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware проверяет заголовок Authorization.
// Здесь показан условный пример, где токен строго равен "Bearer SomeSecretToken".
// В реальном проекте логику проверки можно усложнить (JWT, БД и т.д.).
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// Простейшая проверка — должен быть непустой заголовок
		// и совпадать со строкой "Bearer SomeSecretToken".
		if token == "" || token != "Bearer SomeSecretToken" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or missing token",
			})
			return
		}

		// Всё окей — идём дальше
		c.Next()
	}
}
