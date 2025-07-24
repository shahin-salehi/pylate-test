package db

import (
	"context"
	"log/slog"
	"shahin/webserver/internal/types"
	"golang.org/x/crypto/bcrypt"
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


// EnsureDefaultAdmin creates the 'aiops' group and 'admin@local' user if they don't exist,
// and associates them. This is for bootstrap only.
func (c *crud) EnsureDefaultAdmin(ctx context.Context) error {
	const (
		AdminEmail    = "admin@example.com"
		AdminUsername = "admin"
		AdminPassword = "admin"
	)

	// 1. Ensure 'aiops' group exists
	var groupID int64
	err := c.Conn.QueryRow(ctx, `
		INSERT INTO groups (name)
		VALUES ('aiops')
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id;
	`).Scan(&groupID)
	if err != nil {
		slog.Error("failed to ensure group aiops", slog.Any("error", err))
		return err
	}

	// 2. Check if admin user exists
	var userID int64
	err = c.Conn.QueryRow(ctx, `
		SELECT id FROM users WHERE email = $1;
	`, AdminEmail).Scan(&userID)

	if err != nil {
		// Admin doesn't exist — insert it
		hashed, _ := bcrypt.GenerateFromPassword([]byte(AdminPassword), bcrypt.DefaultCost)

		err = c.Conn.QueryRow(ctx, `
			INSERT INTO users (username, email, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id;
		`, AdminUsername, AdminEmail, string(hashed)).Scan(&userID)
		if err != nil {
			slog.Error("failed to insert default admin user", slog.Any("error", err))
			return err
		}
		slog.Info("✅ Created default admin user", slog.String("email", AdminEmail), slog.String("password", AdminPassword))
	} else {
		slog.Info("✅ Default admin user already exists", slog.String("email", AdminEmail))
	}

	// 3. Link user to group
	var exists bool
	err = c.Conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM user_to_group WHERE user_id = $1 AND group_id = $2
		);
	`, userID, groupID).Scan(&exists)
	if err != nil {
		slog.Error("failed to check user_to_group", slog.Any("error", err))
		return err
	}

	if !exists {
		_, err = c.Conn.Exec(ctx, `
			INSERT INTO user_to_group (user_id, group_id)
			VALUES ($1, $2);
		`, userID, groupID)
		if err != nil {
			slog.Error("failed to insert user_to_group link", slog.Any("error", err))
			return err
		}
		slog.Info("✅ Linked admin user to group aiops")
	}

	return nil
}
