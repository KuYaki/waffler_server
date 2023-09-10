package models

import (
	"database/sql"
	"time"
)

const Telegram = iota

type UserDTO struct {
	ID          int            `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username,omitempty"`
	Hash        string         `json:"password,omitempty"`
	ParserToken sql.NullString `json:"token_parser,omitempty"`
	ParserType  sql.NullString `json:"type_parser"`
	Locale      sql.NullString `json:"locale,omitempty"`
}

type SourceDTO struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	SourceType   int    `json:"source_type"`
	SourceUrl    string `json:"source_url"`
	WafflerScore int    `json:"waffler_score"`
}

type RecordDTO struct {
	ID         int       `json:"id,omitempty"`
	RecordText string    `json:"record_text"`
	Score      int       `json:"score"`
	CreatedAt  time.Time `json:"timestamp"`
	SourceID   int       `json:"source_id,omitempty"`
}
