package routes

import (
	"files/internal/api/handlers"
	"files/internal/api/middlewares"
	"github.com/gin-gonic/gin"
)

func S3Routes(r *gin.Engine, s3Handlers *handlers.S3Handlers) {

	// POST /upload/:id
	r.POST("/upload/:id",
		middlewares.LimitRequestSizeMiddleware(50<<20),
		middlewares.CheckExtensionsMiddleware([]string{".png", ".jpg", ".jpeg", ".gif"}),
		s3Handlers.UploadMultipleHandler,
	)

	// DELETE /upload/:id
	r.DELETE("/upload/:id", s3Handlers.DeleteAllByIDHandler)

	// DELETE /upload/:id/:uuid
	r.DELETE("/upload/:id/:uuid", s3Handlers.DeleteOneByUUIDHandler)
}
