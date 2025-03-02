package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime"
	"mime/multipart"
	"path"

	"files/internal/repository"
)

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
	multipartReader *multipart.Reader,
) ([]string, error) {
	// Префикс (папка) в S3, напр. photos/123
	prefix := fmt.Sprintf("photos/%s", idParam)
	var fileURLs []string

	// Читаем части (part) из multipart.Reader
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения part: %w", err)
		}

		fileName := part.FileName()
		if fileName == "" {
			continue
		}

		ext := path.Ext(fileName)
		fileUUID := uuid.New().String()
		// Формируем ключ в S3: photos/123/uuid.png
		s3Key := fmt.Sprintf("%s/%s%s", prefix, fileUUID, ext)

		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

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

// ListAllFiles — возвращает список URL всех файлов из S3 бакета.
func (s *S3Service) ListAllFiles(ctx context.Context) ([]string, error) {
	objects, err := s.repo.ListAllFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список файлов: %w", err)
	}

	var fileURLs []string
	for _, obj := range objects {
		if obj.Key == nil {
			continue
		}
		fileURL := fmt.Sprintf("https://%s.s3.timeweb.cloud/%s", s.repo.BucketName, *obj.Key)
		fileURLs = append(fileURLs, fileURL)
	}
	return fileURLs, nil
}

// ListFilesInFolder — возвращает список URL файлов по заданному префиксу (папке).
func (s *S3Service) ListFilesInFolder(ctx context.Context, folderName string) ([]string, error) {
	// Если folderName не заканчивается слэшем, дополняем его.
	if folderName[len(folderName)-1] != '/' {
		folderName += "/"
	}
	objects, err := s.repo.ListFilesByPrefix(ctx, folderName)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить файлы по префиксу %s: %w", folderName, err)
	}
	var fileURLs []string
	for _, obj := range objects {
		if obj.Key != nil {
			fileURL := fmt.Sprintf("https://%s.s3.timeweb.cloud/%s", s.repo.BucketName, *obj.Key)
			fileURLs = append(fileURLs, fileURL)
		}
	}
	return fileURLs, nil
}

// GetFolderInfo — возвращает информацию о папке: существует ли папка и список файлов в ней.
func (s *S3Service) GetFolderInfo(ctx context.Context, folderName string) (bool, []string, error) {
	// Приводим folderName к корректному виду: заканчивается слэшом.
	if folderName[len(folderName)-1] != '/' {
		folderName += "/"
	}
	fileURLs, err := s.ListFilesInFolder(ctx, folderName)
	if err != nil {
		return false, nil, fmt.Errorf("ошибка получения файлов для папки %s: %w", folderName, err)
	}
	// Если файлов нет, считаем, что папка не существует.
	exists := len(fileURLs) > 0
	return exists, fileURLs, nil
}
