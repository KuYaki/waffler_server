package gpt

import (
	"context"
	"errors"
	"fmt"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"regexp"
	"strconv"
)

type ChatGPT struct {
	gpt *openai.Client
	log *zap.Logger
}

func NewChatGPT(token string, log *zap.Logger) *ChatGPT {
	client := openai.NewClient(token)
	return &ChatGPT{
		gpt: client,
		log: log,
	}
}

const answerGPTRacism = `Не пиши ничего кроме дробного числа - доля расизма в субъективной оценке гипотетического оценщик, если не уверен пиши 0.5. Оцени долю расизма следующего высказывания:`

func (g *ChatGPT) ConstructQuestionGPT(mess string, typeAnswer models.ScoreType) (*int, error) {
	var answerGPT *openai.ChatCompletionResponse
	var err, errWarn error
	switch typeAnswer {

	case models.Racism:
		answerGPT, errWarn = g.questionForGPT(answerGPTRacism + " " + mess)
		if errWarn != nil {
			g.log.Warn("error: questionForGPT", zap.Error(errWarn))
			return nil, nil
		}

	default:
		answerGPT, err = nil, errors.New(fmt.Sprintf("error: unknown score type %v", typeAnswer))
	}
	if err != nil {
		return nil, err
	}

	score, errWarn := parseAnswerGPT(answerGPT.Choices[0].Message.Content)
	if errWarn != nil {
		g.log.Warn("error: questionForGPT", zap.Error(errWarn))
		return nil, nil
	}

	return &score, nil
}

func parseAnswerGPT(answer string) (int, error) {
	var scoreFloat = []string{"1.0", "0.9", "0.8", "0.7", "0.6", "0.5", "0.4", "0.3", "0.2", "0.1", "0.0", "0"}
	var resRaw string
	var find bool
	for _, v := range scoreFloat {
		re, err := regexp.Compile(string(v))
		if err != nil {
			return 0, err
		}
		res := re.MatchString(answer)
		if res {
			resRaw = v
			find = true
			break
		}
	}
	if !find {
		return 0, errors.New("error: parseAnswerGPT" + answer)
	}
	if resRaw == "0" {
		return 0, nil

	}

	var result int
	float, err := strconv.ParseFloat(resRaw, 64)
	if err != nil {
		return 0, err
	}
	switch float {
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
	default:
		return 0, errors.New(fmt.Sprintf("unknown float: %.2f", float))

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
