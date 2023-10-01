package language_model

import (
	"fmt"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/service"
	"github.com/go-faster/errors"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"regexp"
	"strconv"
)

const (
	answerGPTRacism  = `Не пиши ничего кроме дробного числа - доля расизма в субъективной оценке гипотетического оценщик, если не уверен пиши 0.0. Оцени долю расизма следующего высказывания:`
	answerGPTWaffler = `Оцени по шкале от 0 до 100 на сколько логически противоречат друг другу следующие два блока текста. 0 - это не противоречат, 100 точное логическое противоречие. Если не уверен или не можешь оценить, просто пиши -1 и не надо никак пояснений):`
)

func (w *service.WafflerService) ConstructQuestionGPT(mess string, search *message.ParserRequest) (*int, error) {
	var answerGPT *openai.ChatCompletionResponse
	var err, errWarn error

	chatGPT := gpt.NewChatGPT(search.Parser.Token, w.log)
	switch search.ScoreType {

	case models.Racism:
		answerGPT, errWarn = chatGPT.QuestionForGPT(answerGPTRacism + " " + mess)
		if errWarn != nil {
			w.log.Warn("error: QuestionForGPT", zap.Error(errWarn))
			return nil, nil
		}

	default:
		answerGPT, err = nil, errors.New(fmt.Sprintf("error: unknown score type %v", search.ScoreType))
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
