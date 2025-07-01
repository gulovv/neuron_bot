package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"google.golang.org/api/option"

	"github.com/gulovv/neuron_bot/internal/gemini"
	"github.com/gulovv/neuron_bot/models"
)

func main() {
	log.Println("Загрузка переменных окружения...")
	_ = godotenv.Load()

	botToken := os.Getenv("BOT_TOKEN")
	geminiKey := os.Getenv("GEMINI_API_KEY")

	if botToken == "" || geminiKey == "" {
		log.Fatal("BOT_TOKEN или GEMINI_API_KEY не указаны в .env")
	}

	// Инициализация Gemini
	log.Println("Инициализация клиента Gemini...")
	client, err := genai.NewClient(models.Ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		log.Fatal("Ошибка при создании клиента Gemini:", err)
	}

	models.Client = client
	models.TextModel = client.GenerativeModel(models.TextModel_name)
	models.TextModel.SafetySettings = gemini.DefaultSafety()
	models.ModelMap[models.TextModel_name] = models.TextModel
	

	// Инициализация Telegram-бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Ошибка при создании Telegram-бота:", err)
	}

	// Инициализация Kafka reader
	log.Println("Запуск Kafka consumer...")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "messages",
		GroupID: "tg-bot-consumer",
	})
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Ошибка чтения из Kafka:", err)
			continue
		}

		var msg models.Message
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Println("Ошибка распаковки сообщения:", err)
			continue
		}

		log.Printf("[Consumer] Получено сообщение от %s: %s", msg.Username, msg.Text)

		// Обработка Gemini
		response := gemini.GenerateGeminiResponse(msg.ChatID, models.TextModel_name, genai.Text(msg.Text))
		log.Printf("[Consumer] Ответ от Gemini: %s", response)

		// Редактирование исходного сообщения
		edit := tgbotapi.NewEditMessageText(msg.ChatID, msg.MessageID, response)
		edit.ParseMode = "MarkdownV2"

		_, err = bot.Send(edit)
		if err != nil {
			log.Printf("[Consumer] Ошибка при отправке ответа: %v", err)
		} else {
			log.Printf("[Consumer] Ответ отправлен пользователю %d", msg.ChatID)
		}
	}
}