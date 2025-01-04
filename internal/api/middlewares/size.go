package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LimitRequestSizeMiddleware ограничивает общий размер тела запроса.
// Например, передаём в параметр 50<<20 (50 МБ).
func LimitRequestSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}
