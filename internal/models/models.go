package models

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	TokenGPT string `json:"token_gpt,omitempty"`
}

type Search struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`
	Order  string `json:"order"`
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
