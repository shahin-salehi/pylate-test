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
	DeleteFile(ctx context.Context, id int64) error
	GetCategories(ctx context.Context) ([]string, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	RegisterUser(ctx context.Context, user types.User) (int64, error)
	GetUserGroup(ctx context.Context, userID int64 ) (int64, error) 
}


type crud struct{
	Conn *pgxpool.Pool
}

func (c *crud) ReadFiles(ctx context.Context) ([]types.File, error){
	// change to args email + group as $1 $2 then append i .Query ... arguments
	// this is ultra stupid
	sqlStmt := `SELECT * FROM read_files('tester@example.com', 2)`

	rows, err := c.Conn.Query(ctx, sqlStmt)
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

func (c *crud) DeleteFile(ctx context.Context, id int64) error {
	sqlStmt := `DELETE FROM pdfs WHERE id=$1`
	
	_, err := c.Conn.Exec(ctx, sqlStmt, id)
	
	if err != nil{
		slog.Error("failed to delete pdf", slog.Any("error", err))
		return err
	}
	
	return nil 

}

func (c *crud) GetCategories(ctx context.Context) ([]string, error){
	sqlStmt := `SELECT DISTINCT category FROM pdf_chunks WHERE category IS NOT NULL` 

	rows, err := c.Conn.Query(ctx, sqlStmt)
	if err!= nil {
		slog.Error("failed to execute query ", slog.Any("function", "GetCategories"), slog.Any("error", err))
		return nil, err 
	}
	
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			slog.Error("Error in scan", slog.Any("error", err))
			return nil, err
		}
		categories = append(categories, category)
	}

	if rows.Err() != nil {
		slog.Error("rows contain errors", slog.Any("error",err))
		return nil, err
	}

	return categories, nil


}
