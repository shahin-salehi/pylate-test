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

func NewDatabase(connectionSting string) (*crud, error){
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

	_, err = pool.Exec(context.Background(), createFunctions)
	if err != nil {
		slog.Error("failed to create functions", slog.Any("error", err))
		return nil, err
	}
	
	return &crud{Conn: pool}, err
}
