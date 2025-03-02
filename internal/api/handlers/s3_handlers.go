package handlers

import (
	"context"
	"net/http"

	"files/internal/services"
	"files/pkg/http_error" // <-- Импортируем ваш модуль с ошибками
	"github.com/gin-gonic/gin"
)

type S3Handlers struct {
	S3Service *services.S3Service
}

func NewS3Handler(svc *services.S3Service) *S3Handlers {
	return &S3Handlers{S3Service: svc}
}

// UploadMultipleHandler — POST /upload/:id
func (h *S3Handlers) UploadMultipleHandler(c *gin.Context) {
	multipartReader, err := c.Request.MultipartReader()
	if err != nil {
		http_error.NewHTTPError(
			http.StatusBadRequest,
			"Ошибка чтения multipart-данных",
			[]http_error.ErrorItem{
				{Field: "multipart_data", Error: err.Error()},
			},
		).Send(c)
		return
	}

	idParam := c.Param("id")
	if idParam == "" {
		http_error.NewHTTPError(
			http.StatusBadRequest,
			"Не указан :id в пути",
			[]http_error.ErrorItem{
				{Field: "id", Error: "missing"},
			},
		).Send(c)
		return
	}

	urls, err := h.S3Service.UploadMultiple(context.Background(), idParam, multipartReader)
	if err != nil {
		http_error.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"urls": urls})
}

// DeleteAllByIDHandler — DELETE /upload/:id
func (h *S3Handlers) DeleteAllByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		http_error.NewHTTPError(
			http.StatusBadRequest,
			"Не указан :id в пути",
			[]http_error.ErrorItem{
				{Field: "id", Error: "missing"},
			},
		).Send(c)
		return
	}

	err := h.S3Service.DeleteAllByID(context.Background(), idParam)
	if err != nil {
		http_error.NewHTTPError(
			http.StatusNotFound,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Все файлы удалены"})
}

// DeleteOneByUUIDHandler — DELETE /upload/:id/:uuid
func (h *S3Handlers) DeleteOneByUUIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		http_error.NewHTTPError(
			http.StatusBadRequest,
			"Не указан :id в пути",
			[]http_error.ErrorItem{
				{Field: "id", Error: "missing"},
			},
		).Send(c)
		return
	}

	uuidParam := c.Param("uuid")
	if uuidParam == "" {
		http_error.NewHTTPError(
			http.StatusBadRequest,
			"Не указан :uuid в пути",
			[]http_error.ErrorItem{
				{Field: "uuid", Error: "missing"},
			},
		).Send(c)
		return
	}

	keys, err := h.S3Service.DeleteOneByUUID(context.Background(), idParam, uuidParam)
	if err != nil {
		http_error.NewHTTPError(
			http.StatusNotFound,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Удалён файл(ы) по UUID",
		"keys":    keys,
	})
}

// ListAllFilesHandler — GET /files
// Возвращает список URL всех файлов из S3 бакета.
func (h *S3Handlers) ListAllFilesHandler(c *gin.Context) {
	files, err := h.S3Service.ListAllFiles(context.Background())
	if err != nil {
		http_error.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
			nil,
		).Send(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

// FolderExistsHandler — GET /folder_exists?folder=<folderName>
// Возвращает информацию о папке: существует ли она и список файлов по заданному пути.
// Если файлов нет, возвращается ошибка 404.
func (h *S3Handlers) FolderExistsHandler(c *gin.Context) {
	folderName := c.Query("folder")
	if folderName == "" {
		http_error.NewHTTPError(
			http.StatusBadRequest,
			"Не указан параметр folder",
			[]http_error.ErrorItem{
				{Field: "folder", Error: "missing"},
			},
		).Send(c)
		return
	}

	exists, files, err := h.S3Service.GetFolderInfo(context.Background(), folderName)
	if err != nil {
		http_error.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	if !exists {
		http_error.NewHTTPError(
			http.StatusNotFound,
			"Папка не найдена",
			[]http_error.ErrorItem{
				{Field: "folder", Error: "not found"},
			},
		).Send(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"folder": folderName,
		"exists": exists,
		"files":  files,
	})
}
