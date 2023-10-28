package message

import (
	"time"

	"github.com/KuYaki/waffler_server/internal/models"
)

type UserInfo struct {
	Parser Parser `json:"parser,omitempty"`
	Locale string `json:"locale,omitempty"`
}

type Locale string

type Parser struct {
	Type  string `json:"type,omitempty"`
	Token string `json:"token,omitempty"`
}

type InfoRequest struct {
	Name string            `json:"name"`
	Type models.SourceType `json:"type"`
}

type SourceURL struct {
	SourceUrl string `json:"source_url"`
}

type ParserRequest struct {
	SourceURL string           `json:"source_url"`
	ScoreType models.ScoreType `json:"score_type"`
	Parser    *Parser          `json:"parser"`
	ClientID  string           `json:"client_id"`
}

type SearchResponse struct {
	Sources []models.SourceDTO `json:"sources"`
	Cursor  *Cursor            `json:"cursor"`
}

type Cursor struct {
	Offset    int `json:"offset"`
	Partition int `json:"partition"`
}

type Source struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SourceType   string `json:"source_type"`
	SourceUrl    string `json:"source_url"`
	WafflerScore string `json:"waffler_score"`
}

type ScoreRequest struct {
	SourceId int              `json:"source_id"`
	Type     models.ScoreType `json:"score_type"`
	Limit    int              `json:"limit"`
	Cursor   Cursor           `json:"cursor"`
	Order    []string         `json:"order"`
}

type ScoreResponse struct {
	Records []Record `json:"records"`
	Cursor  *Cursor  `json:"cursor"`
}

type Record struct {
	RecordText string    `json:"record_text,omitempty"`
	Score      int       `json:"score"`
	Timestamp  time.Time `json:"timestamp"`
}

type User struct {
	ID          int    `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	ParserToken string `json:"token_gpt,omitempty"`
	ParserType  string `json:"parser_type,omitempty"`
	Locale      string `json:"locale,omitempty"`
}

type Search struct {
	QueryForName string              `json:"query"`
	Limit        int                 `json:"limit"`
	Cursor       *Cursor             `json:"cursor"`
	Order        []string            `json:"order"`
	SourceType   []models.SourceType `json:"source_type"`
	ScoreType    []models.ScoreType  `json:"score_type"`
}

type SourceParse struct {
	SourceURL string `json:"source_url"`
	ScoreType string `json:"score_type"`
	Parser    `json:"parser"`
}

type PriceRequest struct {
	SourceUrl string           `json:"source_url"`
	ScoreType models.ScoreType `json:"score_type"`
	Parser    Parser           `json:"parser"`
	Limit     int              `json:"limit"`
	Currency  string           `json:"currency"`
}

type PriceResponse struct {
	Price    string `json:"price"`
	Currency string `json:"currency"`
}
