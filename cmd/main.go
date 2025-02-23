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
)

func main() {
	container := ioc.NewContainer()

	// Отключаем режим отладки, чтобы не выводились лишние сообщения
	gin.SetMode(gin.ReleaseMode)

	// Инициализируем новый роутер (без встроенных логов)
	r := gin.New()

	// Добавляем CORS middleware, разрешающий запросы с любых источников
	r.Use(cors.Default())

	// Ограничение, что при multipart-файлах > 8 MB данные пойдут во временный файл
	r.MaxMultipartMemory = 8 << 20

	// Добавляем остальные middleware
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestLoggerMiddleware(container.Logger))

	apiGroup := r.Group("/files")
	//set auth for group
	//apiGroup.Use(auth.JwtAuthMiddleware(container.JwtService))

	routes.S3Routes(apiGroup, container.S3Handler)

	port := env.GetEnv("SERV_PORT", "8080")
	log.Info("Starting server", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
