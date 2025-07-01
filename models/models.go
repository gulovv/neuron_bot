package models

import (
	"context"
	"sync"
	"github.com/google/generative-ai-go/genai"
)

type Message struct {
	ChatID      int64  `json:"chat_id"`
	Username    string `json:"username"`
	Text        string `json:"text"`
	MessageID   int    `json:"message_id"`
	ReplyNeeded bool   `json:"reply_needed"`
}

const (
    TextModel_name   = "gemini-1.5-flash-latest"
)

var (
    Ctx         = context.Background()
    Client      *genai.Client
    TextModel   *genai.GenerativeModel
 
    ModelMap    = make(map[string]*genai.GenerativeModel)
    ChatSession sync.Map
)
