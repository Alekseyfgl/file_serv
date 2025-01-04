package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// LoadEnv загружает переменные из файла .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
}

// GetEnv получает значение переменной окружения
func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
