package main

import (
	"files/configs/env"
	"files/internal/api/handlers"
	"files/internal/api/middlewares"
	"files/internal/repository"
	"files/internal/routes"
	"files/internal/services"
	"files/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	initApp()

	bucketName := env.GetEnv("BUCKET_NAME", "")
	S3AccessKey := env.GetEnv("S3_ACCESS_KEY", "")
	S3SecretAccessKey := env.GetEnv("S3_SECRET_ACCESS_KEY", "")
	if bucketName == "" || S3AccessKey == "" || S3SecretAccessKey == "" {
		log.Fatal("S3 credentials or bucket name are not provided")
	}

	s3Repo := repository.NewS3Repository(bucketName, S3AccessKey, S3SecretAccessKey)
	s3Service := services.NewS3Service(s3Repo)
	s3Handler := handlers.NewS3Handler(s3Service)

	// Отключаем режим отладки, чтобы не выводились лишние сообщения
	gin.SetMode(gin.ReleaseMode)

	// Инициализируем новый роутер (без встроенных логов)
	r := gin.New()
	// Ограничение, что при multipart-файлах > 8 MB данные пойдут во временный файл
	r.MaxMultipartMemory = 8 << 20
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestLoggerMiddleware(log.GetLogger()))

	routes.S3Routes(r, s3Handler)

	port := env.GetEnv("SERV_PORT", "8080")

	log.Info("Starting server", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}

func initApp() {
	env.LoadEnv()
	log.InitLogger()
}
