package services

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path"
	"strings"

	"github.com/google/uuid"

	"files/internal/repository"
)

// Пример белого списка (можно расширять)
var allowedExts = []string{".png", ".jpg", ".jpeg", ".gif"}

// isAllowedExt — проверяет, входит ли расширение в «белый список».
func isAllowedExt(ext string) bool {
	ext = strings.ToLower(ext)
	for _, e := range allowedExts {
		if e == ext {
			return true
		}
	}
	return false
}

// S3Service — слой бизнес-логики для работы с файлами.
type S3Service struct {
	repo *repository.S3Repository
}

// NewS3Service — конструктор, принимает репозиторий.
func NewS3Service(repo *repository.S3Repository) *S3Service {
	return &S3Service{repo: repo}
}

// UploadMultiple — читает файлы из multipart.Reader, заливает их в S3.
func (s *S3Service) UploadMultiple(
	ctx context.Context,
	idParam string,
	multipartReader *multipart.Reader, // <-- используем стандартный *multipart.Reader
) ([]string, error) {

	// Префикс (папка) в S3, напр. photos/123
	prefix := fmt.Sprintf("photos/%s", idParam)

	var fileURLs []string

	// Читаем части (part) из multipart.Reader
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			// Файлы закончились
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения part: %w", err)
		}

		fileName := part.FileName()
		if fileName == "" {
			// Это не файл, а поле формы — пропускаем
			continue
		}

		// Проверяем расширение
		ext := path.Ext(fileName)
		if !isAllowedExt(ext) {
			return nil, fmt.Errorf("расширение %q не поддерживается", ext)
		}

		// Генерируем UUID для имени файла
		fileUUID := uuid.New().String()
		// Формируем ключ в S3: photos/123/uuid.png
		s3Key := fmt.Sprintf("%s/%s%s", prefix, fileUUID, ext)

		// Определяем Content-Type
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// Загружаем файл в S3 (через репозиторий)
		fileURL, err := s.repo.UploadFile(ctx, s3Key, contentType, part)
		if err != nil {
			return nil, fmt.Errorf("ошибка загрузки в S3: %w", err)
		}

		fileURLs = append(fileURLs, fileURL)
	}

	if len(fileURLs) == 0 {
		return nil, fmt.Errorf("в multipart нет файлов")
	}

	return fileURLs, nil
}

// DeleteAllByID — удаляет все файлы в photos/:id/
func (s *S3Service) DeleteAllByID(ctx context.Context, idParam string) error {
	prefix := fmt.Sprintf("photos/%s/", idParam)

	objects, err := s.repo.ListFilesByPrefix(ctx, prefix)
	if err != nil {
		return fmt.Errorf("не удалось получить список файлов: %w", err)
	}
	if len(objects) == 0 {
		return fmt.Errorf("нет файлов с префиксом '%s'", prefix)
	}

	var keys []string
	for _, obj := range objects {
		keys = append(keys, *obj.Key)
	}

	if err := s.repo.DeleteFilesBatch(ctx, keys); err != nil {
		return fmt.Errorf("ошибка удаления файлов: %w", err)
	}
	return nil
}

// DeleteOneByUUID — удаляет один (или несколько) файлов с префиксом photos/:id/:uuid
func (s *S3Service) DeleteOneByUUID(ctx context.Context, idParam, uuidParam string) ([]string, error) {
	prefix := fmt.Sprintf("photos/%s/%s", idParam, uuidParam)

	objects, err := s.repo.ListFilesByPrefix(ctx, prefix)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список файлов: %w", err)
	}
	if len(objects) == 0 {
		return nil, fmt.Errorf("файл с префиксом %q не найден", prefix)
	}

	var keys []string
	for _, obj := range objects {
		keys = append(keys, *obj.Key)
	}

	if err := s.repo.DeleteFilesBatch(ctx, keys); err != nil {
		return nil, fmt.Errorf("ошибка удаления: %w", err)
	}
	return keys, nil
}
