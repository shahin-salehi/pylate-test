package db

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/001_create_tables.sql
var createTables string

//go:embed migrations/002_create_functions.sql
var createFunctions string

func NewDatabase(connectionSting string) (*pgxpool.Pool, error){
	tmp, err := pgxpool.New(context.Background(), connectionSting)
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("error", err))
		return nil, err
	}
	pool, err := pgxpool.New(context.Background(), connectionSting)
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("error", err))
		return nil, err
	}
	
	_, err = tmp.Exec(context.Background(), createTables)
	if err != nil {
		slog.Error("failed to create tables", slog.Any("error", err))
		return nil, err
	}

	_, err = tmp.Exec(context.Background(), createFunctions)
	if err != nil {
		slog.Error("failed to create functions", slog.Any("error", err))
		return nil, err
	}
	defer tmp.Close()
	

	return pool, nil
}
