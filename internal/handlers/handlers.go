package handlers

import (
	"context"
	"log"
	"github.com/gulovv/neuron_bot/internal/kafka"
	"github.com/gulovv/neuron_bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleText(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID
	userText := update.Message.Text
	username := update.Message.From.UserName

	log.Printf("[HandleText] Получено сообщение от пользователя %d: %s", chatID, userText)

	// Отправка "Обрабатываю..." — сразу
	replyMsg := tgbotapi.NewMessage(chatID, "Обрабатываю запрос...")
	sent, err := bot.Send(replyMsg)
	if err != nil {
		log.Printf("[HandleText] Ошибка при отправке сообщения: %v", err)
		return
	}

	// Подготовка сообщения для Kafka
	msg := models.Message{
		ChatID:    chatID,
		Username:  username,
		Text:      userText,
		MessageID: sent.MessageID,
	}

	if err := kafka.SendToKafka(context.Background(), msg); err != nil {
		log.Printf("[HandleText] Ошибка при отправке в Kafka: %v", err)
		return
	}

	log.Printf("[HandleText] Сообщение отправлено в Kafka для обработки. (MessageID: %d)", sent.MessageID)
}