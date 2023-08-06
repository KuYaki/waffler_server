package models

type User struct {
	ID       int    `json:"ID,omitempty"`
	Username string `json:"username,omitempty"`
	Pass     string `json:"pass,omitempty"`
}

type Search struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`
	Order  string `json:"order"`
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
