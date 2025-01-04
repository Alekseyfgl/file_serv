package middlewares

import (
	"github.com/gin-gonic/gin"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Генерируем уникальный идентификатор запроса
		requestID := uuid.New().String()
		// Записываем его в заголовок
		c.Writer.Header().Set("X-Request-ID", requestID)

		start := time.Now()

		// Передаём управление обработке маршрута
		c.Next()

		duration := time.Since(start)

		// Получаем статус ответа
		statusCode := c.Writer.Status()
		// Получаем IP клиента
		clientIP := c.ClientIP()

		logger.Info("Request processed",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", clientIP),
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration),
		)
	}
}
