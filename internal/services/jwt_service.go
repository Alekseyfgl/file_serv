package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// JWTServiceInterface интерфейс, определяющий методы для работы с JWT
type JWTServiceInterface interface {
	// GenerateAccessToken генерирует JWT access токен
	// userId - уникальный идентификатор пользователя
	// expiresIn - продолжительность времени действия токена
	GenerateAccessToken(userId int, expiresIn time.Duration) (string, error)

	// ValidateToken проверяет валидность предоставленного токена
	// Возвращает объект токена или ошибку
	ValidateToken(token string) (*jwt.Token, error)
}

// jwtService структура, содержащая настройки для работы с JWT
type jwtService struct {
	secretKey string      // секретный ключ для подписи токенов
	logger    *zap.Logger // логгер для записи действий
}

// Claims определяет пользовательские данные для хранения в JWT токене
// UserID - идентификатор пользователя
// RegisteredClaims - встроенные поля JWT (например, время истечения)
type Claims struct {
	UserId               int `json:"userId"` // Уникальный идентификатор пользователя
	jwt.RegisteredClaims     // Встроенные стандартные claims (exp, iat и т.д.)
}

// NewJWTService создает новый экземпляр JWTServiceInterface
// secretKey - строка, используемая для подписи токенов
// logger - объект логгера для записи действий
func NewJWTService(secretKey string, logger *zap.Logger) JWTServiceInterface {
	return &jwtService{
		secretKey: secretKey,
		logger:    logger,
	}
}

// GenerateAccessToken генерирует JWT токен с заданным userID и временем действия
// Возвращает подписанный токен в виде строки или ошибку
func (s *jwtService) GenerateAccessToken(userId int, expiresIn time.Duration) (string, error) {
	// Логируем начало генерации токена
	s.logger.Info("Generating access token", zap.Int("userId", userId), zap.Duration("expiresIn", expiresIn))

	// Создаем claims с пользовательскими и стандартными данными
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)), // Устанавливаем время истечения токена
			IssuedAt:  jwt.NewNumericDate(time.Now()),                // Устанавливаем время создания токена
		},
	}

	// Создаем новый токен с методом подписи HS256 и нашими claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с помощью секретного ключа
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		s.logger.Error("Failed to sign token", zap.Error(err))
		return "", err
	}

	// Логируем успешную генерацию токена
	s.logger.Info("Token generated successfully", zap.Int("userId", userId))
	return signedToken, nil
}

// ValidateToken проверяет валидность предоставленного токена
// tokenStr - строковое представление токена
// Возвращает объект токена и nil, если токен валиден, или ошибку, если он недействителен
func (s *jwtService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	// Логируем начало валидации токена
	s.logger.Info("Validating token")

	// Разбираем токен и проверяем его подпись с использованием секретного ключа
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		s.logger.Error("Failed to validate token", zap.Error(err))
		return nil, err
	}

	// Логируем успешную валидацию токена
	s.logger.Info("Token validated successfully")
	return token, nil
}
