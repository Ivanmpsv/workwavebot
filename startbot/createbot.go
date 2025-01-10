package startbot

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func Createbot() (*tgbotapi.BotAPI, error) {
	// Загрузка переменных из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Получаем токен из переменных окружения
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Panic("Ошибка: токен бота не установлен!")
	}

	// Хэширование токена для логирования (не используется в NewBotAPI)
	hashedToken := hashToken(token)
	log.Printf("Хэш токена для проверки: %s", hashedToken)

	// Создаём объект бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	// Включаем отладочный режим
	bot.Debug = true

	return bot, nil
}
