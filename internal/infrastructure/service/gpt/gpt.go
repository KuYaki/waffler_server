package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type ChatGPT struct {
	gpt *openai.Client
}

func NewChatGPT(token string) *ChatGPT {
	client := openai.NewClient(token)
	return &ChatGPT{
		gpt: client,
	}
}

type answerTypeGPT string

const (
	SexismGPT answerTypeGPT = "sexism"
	RacismGPT answerTypeGPT = "racism"
)

const answerGPTRacism = `Не пиши ничего кроме дробного числа - доля расизма в субъективной оценке гипотетического оценщик, если не уверен пиши 0.5. Оцени долю расизма следующего высказывания:`

func (g *ChatGPT) ConstructAnswerGPT(answer string, typeAnswer answerTypeGPT) (*openai.ChatCompletionResponse, error) {
	switch typeAnswer {

	case RacismGPT:
		return g.AnswerForGPT(answerGPTRacism + " " + answer)
	}
	return nil, nil
}
func (g *ChatGPT) AnswerForGPT(answer string) (*openai.ChatCompletionResponse, error) {
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
