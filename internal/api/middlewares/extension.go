package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// AllowedExtensions — структура для хранения «белого списка» расширений.
// Например: AllowedExtensions{[]string{".jpg", ".png", ".gif"}}
type AllowedExtensions struct {
	Extensions []string
}

func CheckExtensionsMiddleware(allowedExts []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 0) Считываем весь Body целиком.
		//    При больших файлах это "убивает" стриминг, так как всё помещается в память (или в tmp).
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Ошибка чтения тела запроса",
				"details": err.Error(),
			})
			return
		}

		// 1) Восстанавливаем заголовок Content-Type, чтобы понять boundary (если multipart)
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			// Если не multipart/form-data, не проверяем.
			// Можно и abort'ить, если хотите строго только multipart.
			// В данном примере просто восстанавливаем Body и пропускаем дальше.
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			c.Next()
			return
		}

		// 2) Парсим boundary
		_, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Ошибка парсинга Content-Type",
				"details": err.Error(),
			})
			return
		}
		boundary, ok := params["boundary"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Не найден boundary в заголовке Content-Type",
			})
			return
		}

		// 3) Создаём multipartReader на основе считанных bodyBytes
		mr := multipart.NewReader(bytes.NewReader(bodyBytes), boundary)

		// 4) Перебираем все части (файлы), проверяем расширения
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   "Ошибка чтения part",
					"details": err.Error(),
				})
				return
			}

			fileName := part.FileName()
			if fileName == "" {
				// Поле формы (не файл) — пропускаем
				continue
			}

			ext := strings.ToLower(path.Ext(fileName))

			// Проверяем, есть ли ext в allowedExts
			if !inSlice(ext, allowedExts) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Файл с расширением %q не разрешён", ext),
				})
				return
			}
		}

		// 5) Если всё ок, восстанавливаем Body, чтобы хендлер мог заново прочитать
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// Идём дальше
		c.Next()
	}
}

// Вспомогательная функция
func inSlice(val string, list []string) bool {
	for _, x := range list {
		if x == val {
			return true
		}
	}
	return false
}
