package storage

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WafflerStorager interface {
	Create(*models.Source) error
	SearchBySourceName(search *models.Search) ([]message.Source, error)
	SelectRecords(idSource int) ([]message.Record, error)
}

func NewWafflerStorage(conn *pgxpool.Pool) WafflerStorager {
	return &WafflerStorage{conn: conn}
}

type WafflerStorage struct {
	conn *pgxpool.Pool
}

func (s WafflerStorage) Create(source *models.Source) error {
	ctx := context.Background()
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}
	sqlSource, argsSource, err := sq.Insert("sources").PlaceholderFormat(sq.Dollar).
		Columns("name", "source_type", "source_url", "waffel_score").
		Values(source.Name, source.SourceType, source.SourceUrl, source.WafflerScore).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sqlSource, argsSource...)
	if err != nil {
		return err
	}

	for _, record := range source.Records {
		sqlRecord, argsRecord, err := sq.Insert("records").PlaceholderFormat(sq.Dollar).
			Columns("source_id", "record_text", "score", "created_at").
			Values(source.ID, record.RecordText, record.Score, record.CreatedAt).ToSql()

		_, err = tx.Exec(ctx, sqlRecord, argsRecord...)
		if err != nil {
			return err
		}

	}

	err = tx.Commit(ctx)

	return err
}

func (s WafflerStorage) SearchBySourceName(search *models.Search) ([]message.Source, error) {
	ctx := context.Background()
	sql, args, err := sq.Select("id", "name", "source_type", "source_url", "waffel_score").
		PlaceholderFormat(sq.Dollar).
		From("sources").Where(sq.Like{"name": "%" + search.QueryForName + "%"}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []message.Source
	for rows.Next() {
		var searchResponse message.Source
		err := rows.Scan(&searchResponse.ID, &searchResponse.Name,
			&searchResponse.SourceType, &searchResponse.SourceUrl, &searchResponse.WafflerScore)
		if err != nil {
			return nil, err
		}
		res = append(res, searchResponse)
	}

	return res, nil
}

func (s WafflerStorage) SelectRecords(idSource int) ([]message.Record, error) {
	ctx := context.Background()
	sql, args, err := sq.Select("record_text", "score", "created_at").
		PlaceholderFormat(sq.Dollar).
		From("records").Where(sq.Eq{"source_id": idSource}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []message.Record
	for rows.Next() {
		var searchResponse message.Record
		err := rows.Scan(&searchResponse.RecordText,
			&searchResponse.Score, &searchResponse.Timestamp)
		if err != nil {
			return nil, err
		}
		res = append(res, searchResponse)
	}

	return res, nil
}
