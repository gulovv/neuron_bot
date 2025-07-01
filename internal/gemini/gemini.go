package gemini

import (
	"errors"
	"fmt"
	"github.com/gulovv/neuron_bot/models"
	"github.com/gulovv/neuron_bot/utils"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)
func GenerateGeminiResponse(chatID int64, modelName string, parts ...genai.Part) string {
    model, ok := models.ModelMap[modelName]
    if !ok || model == nil {
        return "❌ Модель Gemini не инициализирована"
    }

    sessionID := SessionKey(chatID, modelName)
    cs := GetOrCreateSession(modelName, sessionID)
    if cs == nil {
        return "❌ Не удалось создать сессию Gemini"
    }

    iter := cs.SendMessageStream(models.Ctx, parts...)

    var response string
    ch := make(chan string)
    go func() {
        for {
            res, err := iter.Next()
            if errors.Is(err, iterator.Done) {
                break
            }
            if err != nil {
                ch <- err.Error()
                break
            }
            if res != nil && len(res.Candidates) > 0 {
                for _, p := range res.Candidates[0].Content.Parts {
                    ch <- fmt.Sprint(p)
                }
            }
        }
        close(ch)
    }()

    for str := range ch {
        response += str
    }

    if response == "" {
        response = "No response received."
    }

    return utils.EscapeMarkdownV2(response)
}

func GetOrCreateSession(modelName, key string) *genai.ChatSession {
    model, ok := models.ModelMap[modelName]
    if !ok || model == nil {
        return nil
    }

    if val, ok := models.ChatSession.Load(key); ok {
        return val.(*genai.ChatSession)
    }

    session := model.StartChat()
    if modelName == models.TextModel_name {
        models.ChatSession.Store(key, session)
    }
    return session
}

func SessionKey(chatID int64, modelName string) string {
    return fmt.Sprintf("%d-%s", chatID, modelName)
}

func DefaultSafety() []*genai.SafetySetting {
    return []*genai.SafetySetting{
        {Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockNone},
        {Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockNone},
        {Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockNone},
        {Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockNone},
    }
}