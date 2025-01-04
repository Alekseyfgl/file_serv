package http_error

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// HTTPError represents a structure for HTTP errors that includes a status code, message, and optional details.
type HTTPError struct {
	StatusCode int         `json:"status_code"`       // HTTP status code of the error.
	Message    string      `json:"message"`           // Error message describing the issue.
	Details    []ErrorItem `json:"details,omitempty"` // Optional array of error details for additional context.
}

// ErrorItem represents a detailed error structure for specific fields.
type ErrorItem struct {
	Field string `json:"field"` // The field in the request where the error occurred.
	Error string `json:"error"` // The specific error message for the field.
}

// NewHTTPError creates a new instance of HTTPError.
//
// Parameters:
//   - statusCode: The HTTP status code (e.g., 400, 500).
//   - message: A brief error message describing the issue.
//   - details: An optional array of ErrorItem for detailed error information.
//
// Returns:
//   - A pointer to an HTTPError instance.
func NewHTTPError(statusCode int, message string, details []ErrorItem) *HTTPError {
	// Если details == nil, устанавливаем пустой массив
	if details == nil {
		details = []ErrorItem{}
	} else {
		// Преобразуем названия полей в lowercase
		for i := range details {
			details[i].Field = strings.ToLower(details[i].Field)
		}
	}

	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

// Send sends the HTTPError as a JSON response to the client.
//
// Parameters:
//   - ctx: The Fiber context used to send the response.
//
// Returns:
//   - An error if sending the response fails; otherwise, nil.
func (e *HTTPError) Send(ctx *gin.Context) {
	ctx.JSON(e.StatusCode, gin.H{
		"error":   e.Message,
		"details": e.Details,
	})
	// Обычно возвращать error не нужно, так как gin.Context сам обрабатывает ответ.
	// Если хотите прервать цепочку обработчиков, можно вызвать ctx.Abort().
}
