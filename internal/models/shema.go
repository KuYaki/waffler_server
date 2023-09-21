package models

import (
	"database/sql"
	"time"
)

type UserDTO struct {
	ID          int            `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username,omitempty"`
	Hash        string         `json:"password,omitempty"`
	ParserToken sql.NullString `json:"token_parser,omitempty"`
	ParserType  sql.NullString `json:"type_parser"`
	Locale      sql.NullString `json:"locale,omitempty"`
}

type SourceDTO struct {
	ID           int        `json:"id,omitempty" gorm:"primaryKey"`
	Name         string     `json:"name"`
	SourceType   SourceType `json:"source_type"`
	SourceUrl    string     `json:"source_url"`
	WafflerScore int        `json:"waffler_score"`
	RacismScore  int        `json:"racism_score"`
}

type RecordDTO struct {
	ID         int       `json:"id,omitempty" gorm:"primaryKey"`
	RecordText string    `json:"record_text"`
	Score      int       `json:"score"`
	ScoreType  ScoreType `json:"score_type"`
	CreatedAt  time.Time `json:"timestamp"`
	SourceID   int       `json:"source_id,omitempty"`
}

type ScoreType int

const (
	Waffler ScoreType = iota
	Racism
)

type SourceType int

const (
	Telegram SourceType = iota
	Youtube
)
