package routes

import (
	"files/internal/api/handlers"
	"files/internal/api/middlewares"
	"github.com/gin-gonic/gin"
)

func S3Routes(r *gin.RouterGroup, s3Handlers *handlers.S3Handlers) {
	r.POST("/upload/:id",
		middlewares.LimitRequestSizeMiddleware(50<<20),
		middlewares.CheckExtensionsMiddleware([]string{".png", ".jpg", ".jpeg", ".gif", ".webp"}),
		s3Handlers.UploadMultipleHandler,
	)

	r.DELETE("/upload/:id", s3Handlers.DeleteAllByIDHandler)

	r.DELETE("/upload/:id/:uuid", s3Handlers.DeleteOneByUUIDHandler)

	// Новый маршрут для получения списка всех файлов
	r.GET("/objects", s3Handlers.ListAllFilesHandler)

	// Новый маршрут для проверки существования папки в S3
	r.GET("/objects/exists", s3Handlers.FolderExistsHandler)
}
