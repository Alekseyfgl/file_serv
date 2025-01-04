package ioc

import (
	"files/configs/env"
	"files/internal/api/handlers"
	"files/internal/repository"
	"files/internal/services"
	"files/pkg/log"
	"go.uber.org/zap"
)

type Container struct {
	Logger     *zap.Logger
	S3Repo     *repository.S3Repository
	S3Service  *services.S3Service
	JwtService services.JWTServiceInterface
	S3Handler  *handlers.S3Handlers
}

// NewContainer - создаем контейнер с зависимостями.
func NewContainer() *Container {
	env.LoadEnv()

	// Initialize logger
	log.InitLogger()
	// Get global logger
	logger := log.GetLogger()

	bucketName := env.GetEnv("BUCKET_NAME", "")
	S3AccessKey := env.GetEnv("S3_ACCESS_KEY", "")
	S3SecretAccessKey := env.GetEnv("S3_SECRET_ACCESS_KEY", "")
	if bucketName == "" || S3AccessKey == "" || S3SecretAccessKey == "" {
		log.Fatal("S3 credentials or bucket name are not provided")
	}

	// Create repositories
	s3Repo := repository.NewS3Repository(bucketName, S3AccessKey, S3SecretAccessKey)
	// Create services
	s3Service := services.NewS3Service(s3Repo)
	jwtService := services.NewJWTService(env.GetEnv("JWT_KEY", ""), logger)

	// Create handlers
	s3Handler := handlers.NewS3Handler(s3Service)
	// Return the container with all dependencies
	return &Container{
		Logger:     logger,
		S3Repo:     s3Repo,
		JwtService: jwtService,
		S3Handler:  s3Handler,
	}
}
