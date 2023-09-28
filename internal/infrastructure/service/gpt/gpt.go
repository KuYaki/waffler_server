package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type ChatGPT struct {
	gpt *openai.Client
	log *zap.Logger
}

type AiLanguageModel interface {
	QuestionForGPT(answer string) (*openai.ChatCompletionResponse, error)
}

func NewChatGPT(token string, log *zap.Logger) AiLanguageModel {
	client := openai.NewClient(token)
	return &ChatGPT{
		gpt: client,
		log: log,
	}
}

func (g *ChatGPT) QuestionForGPT(answer string) (*openai.ChatCompletionResponse, error) {
	resp, err := g.gpt.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: answer,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
