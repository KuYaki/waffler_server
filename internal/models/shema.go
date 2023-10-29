package models

import (
	"time"
)

type UserDTO struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username,omitempty"`
	PwdHash  string `json:"password,omitempty"`
}

type SourceDTO struct {
	ID           int         `json:"id,omitempty" gorm:"primaryKey"`
	Name         string      `json:"name"`
	SourceType   SourceType  `json:"source_type"`
	SourceUrl    string      `json:"source_url"`
	WafflerScore NullFloat64 `json:"waffler_score"`
	RacismScore  NullFloat64 `json:"racism_score"`
}

type RecordDTO struct {
	ID         int       `json:"id,omitempty" gorm:"primaryKey"`
	RecordText string    `json:"record_text"`
	RecordURL  string    `json:"record_url"`
	CreatedTs  time.Time `json:"created_ts"`
	SessionTs  time.Time `json:"session_ts"`
	SourceID   int       `json:"source_id,omitempty"`
}

type WafflerDTO struct {
	ID              int        `json:"id,omitempty" gorm:"primaryKey"`
	Score           NullInt64  `json:"score"`
	ParserType      ParserType `json:"parser"`
	RecordIDBefore  int        `json:"record_id_before"`
	RecordIDAfter   int        `json:"record_id_after"`
	CreatedTsBefore time.Time  `json:"timestamp_before"`
	CreatedTsAfter  time.Time  `json:"timestamp_after"`
	SourceID        int        `json:"source_id,omitempty"`
}

type RacismDTO struct {
	ID         int        `json:"id,omitempty" gorm:"primaryKey"`
	Score      NullInt64  `json:"score"`
	ParserType ParserType `json:"score_type"`
	CreatedTs  time.Time  `json:"created_ts"`
	RecordID   int        `json:"record_id"`
	SourceID   int        `json:"source_id"`
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

type ParserType int

const (
	GPT3_5TURBO ParserType = iota
	GPT4
)
