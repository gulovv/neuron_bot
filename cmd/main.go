package main

import (
    "log"
    "os"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/joho/godotenv"

    "github.com/gulovv/neuron_bot/internal/handlers"
)

func main() {
    log.Println("Загрузка .env...")
    _ = godotenv.Load()

    botToken := os.Getenv("BOT_TOKEN")
    if botToken == "" {
        log.Fatal("BOT_TOKEN не указан в .env")
    }

    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Fatal("Ошибка Telegram:", err)
    }
    log.Printf("Бот запущен как %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates := bot.GetUpdatesChan(u)

    log.Println("Ожидание сообщений...")

    for update := range updates {
        if update.Message == nil {
            continue
        }

        switch {
        case update.Message.IsCommand():
            handlers.HandleCommand(update, bot)

        case update.Message.Text != "":
            handlers.HandleText(update, bot)
        }
    }
}