package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SqlClient struct {
	db *sqlx.DB
}

func NewSqlClient(db *sqlx.DB) *SqlClient {
	return &SqlClient{
		db: db,
	}
}

func (s *SqlClient) GetItem(ctx context.Context, id string) (*Item, error) {
	var result Item

	query := sq.Select("*").From("inventory.items").Where(sq.Eq{"id": id})

	//Use $1 format for psql
	sql, args, _ := query.PlaceholderFormat(sq.Dollar).ToSql()

	//Allow sqlx to inject variables
	if err := s.db.GetContext(ctx, &result, sql, args...); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *SqlClient) UpsertItem(ctx context.Context, item Item) (*Item, error) {

	insertCols := []string{"name", "description", "price"}
	insertVals := []interface{}{item.Name, item.Description, item.Price}

	if item.ID != "" {
		insertCols = append(insertCols, "id")
		insertVals = append(insertVals, item.ID)
	}

	//upsert squirrel query returning all fields
	query := sq.Insert("inventory.items").Columns(insertCols...).
		Values(insertVals...).
		Suffix("ON CONFLICT (id) DO UPDATE SET name = ?, description = ?, price = ?, updated_at = now()", item.Name, item.Description, item.Price).
		Suffix("RETURNING *")

	//Use $1 format for psql
	sql, args, _ := query.PlaceholderFormat(sq.Dollar).ToSql()

	var result Item
	//Allow sqlx to inject variables
	if err := s.db.GetContext(ctx, &result, sql, args...); err != nil {
		return nil, err
	}

	return &result, nil
}
