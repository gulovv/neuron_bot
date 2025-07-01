package utils

import (
	"strings"
)



func EscapeHTML(text string) string {
    text = strings.ReplaceAll(text, "&", "&amp;")
    text = strings.ReplaceAll(text, "<", "&lt;")
    text = strings.ReplaceAll(text, ">", "&gt;")
    text = strings.ReplaceAll(text, `"`, "&quot;")
    text = strings.ReplaceAll(text, `'`, "&#39;")
    return text
}



// EscapeMarkdownV2 экранирует специальные символы для Telegram MarkdownV2 и восстанавливает поддерживаемые конструкции форматирования.
// escapeAllMarkdownV2 экранирует все специальные символы MarkdownV2.
func EscapeMarkdownV2(text string) string {
	specialChars := []string{
		"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|",
		"{", "}", ".", "!", "\\",
	}

	// Экранируем обратный слеш в первую очередь
	text = strings.ReplaceAll(text, "\\", "\\\\")

	// Экранируем остальные символы
	for _, char := range specialChars {
		if char != "\\" {
			text = strings.ReplaceAll(text, char, "\\"+char)
		}
	}

	return text
}



