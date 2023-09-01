package message

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	"time"
)

type UserInfo struct {
	Parser Parser `json:"parser,omitempty"`
	Locale string `json:"locale,omitempty"`
}

type Parser struct {
	Type  string `json:"type,omitempty"`
	Token string `json:"token,omitempty"`
}

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
	Sources []Source `json:"sources"`
	Cursor  int      `json:"cursor"`
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
	Score      string    `json:"score"`
	Timestamp  time.Time `json:"timestamp"`
}
