package models

import "time"

func (d *UserDTO) GetShemaName() string {
	return "users"
}

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username,omitempty"`
	Hash     string `json:"password,omitempty"`
	TokenGPT string `json:"token_gpt,omitempty"`
}

type Source struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	SourceType   string `json:"source_type"`
	SourceUrl    string `json:"source_url"`
	WafflerScore int    `json:"waffler_score"`
}

type Record struct {
	ID         int       `json:"id,omitempty"`
	RecordText string    `json:"record_text"`
	Score      int       `json:"score"`
	Timestamp  time.Time `json:"timestamp"`
	RecordID   int       `json:"source_id,omitempty"`
}
