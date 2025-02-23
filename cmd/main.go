package main

import (
	"files/configs/env"
	"files/internal/api/middlewares"
	"files/internal/ioc"
	"files/internal/routes"
	"files/pkg/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func main() {
	container := ioc.NewContainer()

	// Отключаем режим отладки, чтобы не выводились лишние сообщения
	gin.SetMode(gin.ReleaseMode)

	// Инициализируем новый роутер (без встроенных логов)
	r := gin.New()

	// Настройка CORS для разрешения всех запросов и поддержки отправки файлов
	corsConfig := cors.Config{
		// Разрешаем запросы с любых источников
		AllowAllOrigins: true,
		// Разрешаем все основные HTTP-методы
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		// Разрешаем основные заголовки, включая необходимые для отправки файлов
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		// Заголовки, которые могут быть видны на стороне клиента
		ExposeHeaders: []string{"Content-Length"},
		// Если требуется, можно передавать куки
		AllowCredentials: true,
		// Время, в течение которого результаты preflight-запроса кэшируются
		MaxAge: 12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Ограничение: при multipart-файлах > 8 MB данные пойдут во временный файл
	r.MaxMultipartMemory = 8 << 20

	// Добавляем остальные middleware
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestLoggerMiddleware(container.Logger))

	apiGroup := r.Group("/files")
	//set auth for group
	//apiGroup.Use(auth.JwtAuthMiddleware(container.JwtService))

	routes.S3Routes(apiGroup, container.S3Handler)

	port := env.GetEnv("SERV_PORT", "3000")
	log.Info("Starting server", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
