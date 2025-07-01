package handlers

import (
	"github.com/gulovv/neuron_bot/internal/gemini"
	"github.com/gulovv/neuron_bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func HandleCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
    chatID := update.Message.Chat.ID
    switch update.Message.Command() {
    case "start":
        msg := tgbotapi.NewMessage(chatID, "Привет! Я бот Gemini. Отправь текст или фото, чтобы получить ответ.")
        bot.Send(msg)
    case "help":
        msg := tgbotapi.NewMessage(chatID, "Отправь мне сообщение или фото, и я отвечу с помощью Gemini AI.")
        bot.Send(msg)
    case "clear":
        key := gemini.SessionKey(chatID, models.TextModel_name)
        models.ChatSession.Delete(key)
        msg := tgbotapi.NewMessage(chatID, "Сессия очищена.")
        bot.Send(msg)
    default:
        msg := tgbotapi.NewMessage(chatID, "Неизвестная команда.")
        bot.Send(msg)
    }
}