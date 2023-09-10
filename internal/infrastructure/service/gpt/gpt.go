package gpt

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"strconv"
	"strings"
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

type AnswerTypeGPT int

const (
	RacismGPT AnswerTypeGPT = iota
	SexismGPT
)

const answerGPTRacism = `Не пиши ничего кроме дробного числа - доля расизма в субъективной оценке гипотетического оценщик, если не уверен пиши 0.5. Оцени долю расизма следующего высказывания:`

func (g *ChatGPT) ConstructQuestionGPT(mess string, typeAnswer AnswerTypeGPT) (int, error) {
	var answerGPT *openai.ChatCompletionResponse
	var err error
	switch typeAnswer {

	case RacismGPT:
		answerGPT, err = g.questionForGPT(answerGPTRacism + " " + mess)

	default:
		answerGPT, err = nil, errors.New(fmt.Sprintf("error: unknown type mess %v", typeAnswer))
	}
	if err != nil {
		return 0, err
	}

	score, err := parseAnswerGPT(answerGPT.Choices[0].Message.Content)

	return score, err
}

func parseAnswerGPT(answer string) (int, error) {
	var result int
	answer = strings.TrimSpace(answer)
	scoreFloat, err := strconv.ParseFloat(answer, 64)
	if err != nil {
		return 0, err
	}
	switch scoreFloat {
	case 1.0:
		result = 1
	case 0.9:
		result = 9
	case 0.8:
		result = 8
	case 0.7:
		result = 7
	case 0.6:
		result = 6
	case 0.5:
		result = 5
	case 0.4:
		result = 4
	case 0.3:
		result = 3
	case 0.2:
		result = 2
	case 0.1:
		result = 1
	case 0.0:
		result = 0

	}
	return result, nil

}
func (g *ChatGPT) questionForGPT(answer string) (*openai.ChatCompletionResponse, error) {
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