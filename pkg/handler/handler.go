package handler

import (
	"context"
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
)

type Handler interface {
	HandleMessage(chatID int64, content string) (*api.MessageConfig, error)
}

type OpenAIHandler struct {
	client *openai.Client
}

func NewOpenAIHandler(client *openai.Client) *OpenAIHandler {
	return &OpenAIHandler{
		client: client,
	}
}

func (handler *OpenAIHandler) HandleMessage(chatID int64, content string) (*api.MessageConfig, error) {
	resp, err := handler.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	answer := resp.Choices[0].Message.Content

	msg := api.NewMessage(chatID, answer)
	return &msg, nil
}
