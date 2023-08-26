package models

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Hash     string `json:"password,omitempty"`
	TokenGPT string `json:"token_gpt,omitempty"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	TokenGPT string `json:"token_gpt,omitempty"`
}

func (d *UserDTO) GetShemaName() string {
	return "users"
}

type Search struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`
	Order  string `json:"order"`
}

type Score struct {
	SourceId string `json:"source_id"`
	Type     string `json:"type"`
	Limit    int    `json:"limit"`
	Cursor   string `json:"cursor"`
	Order    string `json:"order"`
}

type SearchResult struct {
	Sources []Source `json:"sources"`
	Cursor  string   `json:"cursor"`
}

type Source struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	SourceType   string `json:"source_type"`
	SourceUrl    string `json:"source_url"`
	WafflerScore string `json:"waffler_score"`
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
