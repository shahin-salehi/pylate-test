package db

import (
	"context"
	"log/slog"
	"shahin/webserver/internal/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Crud interface{
	ReadFiles(ctx context.Context, groupID int64) ([]types.File, error)
	DeleteFile(ctx context.Context, ID int64, groupID int64) error
	GetCategories(ctx context.Context, groupID int64) ([]string, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	RegisterUser(ctx context.Context, user types.User) (int64, error)
	GetUserGroup(ctx context.Context, userID int64 ) (int64, error) 
	EnsureDefaultAdmin(ctx context.Context) error 
}


type crud struct{
	Conn *pgxpool.Pool
}

func (c *crud) ReadFiles(ctx context.Context, groupID int64) ([]types.File, error){
	// change to args email + group as $1 $2 then append i .Query ... arguments
	// this is ultra stupid
	sqlStmt := `SELECT * FROM read_files($1)`

	rows, err := c.Conn.Query(ctx, sqlStmt, groupID)
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

func (c *crud) DeleteFile(ctx context.Context, ID int64, groupID int64) error {
	sqlStmt := `DELETE FROM pdfs WHERE id=$1 AND owner=$2`
	
	_, err := c.Conn.Exec(ctx, sqlStmt, ID, groupID)
	
	if err != nil{
		slog.Error("failed to delete pdf", slog.Any("error", err))
		return err
	}
	
	return nil 

}

func (c *crud) GetCategories(ctx context.Context, groupID int64) ([]string, error){
	sqlStmt := `
		SELECT DISTINCT category
		FROM pdf_chunks
		JOIN pdfs ON pdf_chunks.pdf_id = pdfs.id
		WHERE pdfs.owner = $1 AND category IS NOT NULL;
	`


	rows, err := c.Conn.Query(ctx, sqlStmt, groupID)
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


func NewCrud(pool *pgxpool.Pool) *crud{
	return &crud{
		Conn: pool,
	}

}
