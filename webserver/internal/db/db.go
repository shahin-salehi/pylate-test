package db

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/001_create_tables.sql
var createTables string

func NewDatabase(connectionSting string) (*pgxpool.Pool, error){
	pool, err := pgxpool.New(context.Background(), connectionSting)
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("error", err))
		return nil, err
	}
	
	_, err = pool.Exec(context.Background(), createTables)
	if err != nil {
		slog.Error("failed to create tables", slog.Any("error", err))
		return nil, err
	}

	//todo add functions

	return pool, nil
}
