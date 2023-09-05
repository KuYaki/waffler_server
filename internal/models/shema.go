package models

import (
	"database/sql"
	"time"
)

func (d *UserDTO) GetShemaName() string {
	return "users"
}

const Telegram = iota

type UserDTO struct {
	ID          int            `json:"id"`
	Username    string         `json:"username,omitempty"`
	Hash        string         `json:"password,omitempty"`
	ParserToken sql.NullString `json:"token_parser,omitempty"`
	ParserType  sql.NullString `json:"type_parser"`
	Locale      sql.NullString `json:"locale,omitempty"`
}

type Source struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	SourceType   int    `json:"source_type"`
	SourceUrl    string `json:"source_url"`
	WafflerScore int    `json:"waffler_score"`
	Records      []Record
}

type Record struct {
	ID         int       `json:"id,omitempty"`
	RecordText string    `json:"record_text"`
	Score      int       `json:"score"`
	CreatedAt  time.Time `json:"timestamp"`
	RecordID   int       `json:"source_id,omitempty"`
}
