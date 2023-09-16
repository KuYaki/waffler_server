package message

import "github.com/KuYaki/waffler_server/internal/models"

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
	Cursor       int                 `json:"cursor"`
	Order        string              `json:"order"`
	SourceType   []models.SourceType `json:"source_type"`
	ScoreType    []models.ScoreType  `json:"score_type"`
}

type SourceParse struct {
	SourceURL string `json:"source_url"`
	ScoreType string `json:"score_type"`
	Parser    `json:"parser"`
}
