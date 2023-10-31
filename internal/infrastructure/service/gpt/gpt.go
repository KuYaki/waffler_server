package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type ChatGPT struct {
	gpt *openai.Client
}

type AiLanguageModel interface {
	QuestionForGPT(answer string) (*openai.ChatCompletionResponse, error)
}

func NewAiLanguageModel(token string) AiLanguageModel {
	client := openai.NewClient(token)
	return &ChatGPT{
		gpt: client,
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
