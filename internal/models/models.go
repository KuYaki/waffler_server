package models

type User struct {
	ID          int    `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	ParserToken string `json:"token_gpt,omitempty"`
	ParserType  string `json:"parser_type,omitempty"`
	Locale      string `json:"locale,omitempty"`
}

type Search struct {
	QueryForName string `json:"query"`
	Limit        int    `json:"limit"`
	Cursor       int    `json:"cursor"`
	Order        string `json:"order"`
}

type SourceParse struct {
	SourceURL string `json:"source_url"`
	ScoreType string `json:"score_type"`
	Parser    `json:"parser"`
}
type Parser struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}
