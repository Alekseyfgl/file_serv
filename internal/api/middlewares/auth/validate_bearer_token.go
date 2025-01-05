package auth

import (
	"files/internal/services"
	"files/pkg/http_error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JwtAuthMiddleware проверяет токен и добавляет userId в контекст Gin.
func JwtAuthMiddleware(jwtService services.JWTServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httpErr := http_error.NewHTTPError(http.StatusUnauthorized, "Authorization header is missing", nil)
			c.JSON(httpErr.StatusCode, httpErr)
			c.Abort()
			return
		}

		// Проверяем, что заголовок начинается с "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			httpErr := http_error.NewHTTPError(http.StatusUnauthorized, "Invalid authorization format", nil)
			c.JSON(httpErr.StatusCode, httpErr)
			c.Abort()
			return
		}

		// Извлекаем сам токен
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидируем токен
		token, err := jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			httpErr := http_error.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token", nil)
			c.JSON(httpErr.StatusCode, httpErr)
			c.Abort()
			return
		}

		// Извлекаем из токена ваши кастомные claims
		claims, ok := token.Claims.(*services.Claims)
		if !ok || claims.UserId <= 0 {
			// Здесь вы сами решаете, что считать "валидным" идентификатором
			httpErr := http_error.NewHTTPError(http.StatusUnauthorized, "Invalid token claims", nil)
			c.JSON(httpErr.StatusCode, httpErr)
			c.Abort()
			return
		}

		// Сохраняем userId (число) в контекст Gin
		c.Set("userId", claims.UserId)

		// Двигаемся дальше
		c.Next()
	}
}

func GetUserId(c *gin.Context) (int, bool) {
	userID, exists := c.Get("userId")
	if !exists {
		return 0, false
	}

	// Приводим к int64
	idNum, ok := userID.(int)
	if !ok {
		return 0, false
	}

	return idNum, true
}
