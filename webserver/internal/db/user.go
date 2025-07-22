package db

import (
	"context"
	"log/slog"
	"shahin/webserver/internal/types"
)


func (c *crud) GetUserByEmail(ctx context.Context, email string) (*types.User, error){
	sqlStmt := `SELECT * FROM get_user_by_email($1);`
	
	row := c.Conn.QueryRow(ctx, sqlStmt, email)
	user := new(types.User)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		slog.Error("failed to fetch user from database", slog.Any("function", "GetUserByEmail"), slog.Any("error", err))
		return nil, err
	}
	return user, nil
}

// returns id of user so session can be set
// hey only allow aiops
func (c *crud) RegisterUser(ctx context.Context, user types.User) (int64, error){
	sqlStmt := `SELECT register_user($1, $2, $3);`

	var ID int64
	err := c.Conn.QueryRow(ctx, sqlStmt, user.Username, user.Email, user.PasswordHash).Scan(&ID)
	if err != nil {
		slog.Error("failed to register user", slog.Any("function","RegisterUser"), slog.Any("error", err))
		return 0, err
	}
	
	slog.Info("user registered")
	return ID, nil
}


func (c *crud) GetUserGroup(ctx context.Context, userID int64) (int64, error){
	sqlStmt := `SELECT group_id FROM user_to_group WHERE user_id=$1`

	var groupID int64
	err := c.Conn.QueryRow(ctx, sqlStmt, userID).Scan(&groupID)
	if err != nil {
		slog.Error("failed to get user group", slog.Any("function","GetUserGroup"), slog.Any("error", err))
		return 0, err
	}
	
	slog.Info("user registered")
	return groupID, nil
}
