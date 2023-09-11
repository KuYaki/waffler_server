package message

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/go-faster/errors"
	"github.com/goccy/go-json"
	"time"
)

type UserInfo struct {
	Parser Parser `json:"parser,omitempty"`
	Locale `json:"locale,omitempty"`
}

func (l *UserInfo) UnmarshalJSON(data []byte) error {
	required := struct {
		Parser struct {
			TypePars string `json:"type,omitempty"`
			Token    string `json:"token,omitempty"`
		} `json:"parser,omitempty"`
		Locale string `json:"locale,omitempty"`
	}{}

	err := json.Unmarshal(data, &required)
	if err != nil {
		return err
	}
	switch required.Parser.TypePars {
	case "GPT":
		l.Parser.Type = GPT
	default:
		return errors.New("unknown parser type")
	}

	l.Parser.Token = required.Parser.Token

	switch required.Locale {
	case "RU":
		l.Locale = Russian
	case "EN":
		l.Locale = English
	default:
		return errors.New("unknown locale")
	}

	return nil
}

type Locale string

const (
	Russian Locale = "RU"
	English Locale = "EN"
)

type Parser struct {
	Type  ParserType `json:"type,omitempty"`
	Token string     `json:"token,omitempty"`
}

type ParserType string

const (
	GPT ParserType = "GPT"
)

type InfoRequest struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type SourceURL struct {
	SourceUrl string `json:"source_url"`
}

type ParserRequest struct {
	SourceURL string            `json:"source_url"`
	ScoreType gpt.AnswerTypeGPT `json:"score_type"`
	Parser    *Parser           `json:"parser"`
	ClientID  string            `json:"client_id"`
}

type SearchResponse struct {
	Sources []models.SourceDTO `json:"sources"`
	Cursor  int                `json:"cursor"`
}

type Source struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SourceType   string `json:"source_type"`
	SourceUrl    string `json:"source_url"`
	WafflerScore string `json:"waffler_score"`
}

type ScoreRequest struct {
	SourceId int    `json:"source_id"`
	Type     string `json:"type"`
	Limit    int    `json:"limit"`
	Cursor   int    `json:"cursor"`
	Order    string `json:"order"`
}

type ScoreResponse struct {
	Records []Record `json:"records"`
	Cursor  int      `json:"cursor"`
}

type Record struct {
	RecordText string    `json:"record_text,omitempty"`
	Score      int       `json:"score"`
	Timestamp  time.Time `json:"timestamp"`
}
