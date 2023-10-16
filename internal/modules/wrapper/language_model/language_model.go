package language_model

import (
	"fmt"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/go-faster/errors"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
)

const (
	answerGPTRacism  = `Оцени по шкале от 0 до 100 на сколько рассисткий следующий  текст. 0 - это не рассисткие, 100 точное расисткие. Если не уверен или не можешь оценить, просто пиши -1 и не надо никак пояснений:`
	answerGPTWaffler = `Оцени по шкале от 0 до 100 на сколько логически противоречат друг другу следующие два блока текста. 0 - это не противоречат, 100 точное логическое противоречие. Если не уверен или не можешь оценить, просто пиши -1 и не надо никак пояснений:`
)

type LanguageModel interface {
	ConstructQuestionGPT(mess string, scoreType models.ScoreType) (*int, error)
}

func NewChatGPTWrapper(token string, log *zap.Logger) LanguageModel {
	return &ChatGPT{
		gpt: gpt.NewChatGPT(token),
		log: log,
	}
}

type ChatGPT struct {
	gpt gpt.AiLanguageModel
	log *zap.Logger
}

func (w *ChatGPT) ConstructQuestionGPT(mess string, scoreType models.ScoreType) (*int, error) {
	var answerGPT *openai.ChatCompletionResponse
	var err, errWarn error

	switch scoreType {

	case models.Racism:
		answerGPT, errWarn = w.gpt.QuestionForGPT(answerGPTRacism + " " + mess)
		if errWarn != nil {
			w.log.Warn("error: QuestionForGPT", zap.Error(errWarn))
			return nil, nil
		}
	case models.Waffler:
		answerGPT, errWarn = w.gpt.QuestionForGPT(answerGPTWaffler + " " + mess)
		if errWarn != nil {
			w.log.Warn("error: QuestionForGPT", zap.Error(errWarn))
			return nil, nil
		}
	default:
		answerGPT, err = nil, errors.New(fmt.Sprintf("error: unknown score type %v", scoreType))
	}
	if err != nil {
		return nil, err
	}

	score, errWarn := parseAnswerGPT(answerGPT.Choices[0].Message.Content)
	if errWarn != nil {
		w.log.Warn("error: QuestionForGPT", zap.Error(errWarn))
		return nil, nil
	}

	return &score, nil
}

func parseAnswerGPT(answer string) (int, error) {
	if strings.Contains(answer, "-1") {
		return -1, nil
	}

	// Само регулярный выражение для поиска числа от 0 до 100
	re := regexp.MustCompile(`\b(?:100|[1-9]?[0-9])\b`)

	// Используем регулярное выражение для поиска всех соответствий в строке
	matches := re.FindAllString(answer, 1)

	if len(matches) == 0 {
		return -1, nil
	}
	res, err := strconv.Atoi(matches[0])
	if err != nil {
		return 0, err
	}

	return res, nil
}
