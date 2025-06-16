package db

import (
	"context"
	"log/slog"
	"shahin/webserver/internal/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Crud interface{
	ReadFiles(ctx context.Context) ([]types.File, error)
}


type crud struct{
	Conn *pgxpool.Pool
}

func (c *crud) ReadFiles(ctx context.Context) ([]types.File, error){
	// change to args email + group as $1 $2 then append i .Query ... arguments
	sql_stmt := `SELECT * FROM read_files('shahin@example.com', 1)`

	rows, err := c.Conn.Query(ctx, sql_stmt)
	if err != nil {
		slog.Error("failed to execute query ", slog.Any("function", "ReadFiles"), slog.Any("error", err))
		return nil, err
	}
	
	structuredRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.File])
	if err != nil {
		slog.Error("failed to collect rows", slog.Any("function", "ReadFiles"), slog.Any("error", err))
		return nil, err
	}

	return structuredRows, nil
	
}
