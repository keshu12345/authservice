package dao

import (
	"fmt"

	"github.com/keshucs12345/authservice/models"
)

func (a authDAO) GetUserByAPIKey(apiKey string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE api_key = $1`

	err := a.DB.Get(&user, query, apiKey)
	if err != nil {
		a.Logger.Errorf("failed to get user by API key: %v", err)
		return nil, fmt.Errorf("failed to get user by API key: %w", err)
	}

	a.Logger.Infof("User authenticated successfully with API Key :%s with username : %s", apiKey, user.Username)
	return &user, nil
}


func (a authDAO) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := a.DB.Get(&user, "SELECT * FROM users WHERE user_id = $1", userID)
	if err != nil {
		a.Logger.Errorf("get user by ID failed: %v", err)
		return nil, fmt.Errorf("get user by ID failed: %w", err)
	}

	a.Logger.Infof("User authenticated successfully with ID :%s with username : %s", userID, user.Username)
	return &user, nil
}

func (a authDAO) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := a.DB.Get(&user, "SELECT * FROM public.users WHERE username = $1", username)
	if err != nil {
		a.Logger.Errorf("failed to get user by username: %v", err)
		return nil, err
	}

	a.Logger.Infof("User authenticated successfully with username :%s with ID : %s", username, user.UserID)
	return &user, nil
}

func (a authDAO) InsertUser(user *models.User) error {
	query := `INSERT INTO users (user_id, username, password_hash, role, api_key) 
              VALUES ($1, $2, $3, $4, $5)`

	result, err := a.DB.Exec(query,
		user.UserID, user.Username, user.PasswordHash, user.Role, user.APIKey)

	if err != nil {
		a.Logger.Errorf("insert user failed: %v", err)
		return fmt.Errorf("insert user failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		a.Logger.Errorf("failed to fetch rows affected: %v", err)
		return fmt.Errorf("failed to fetch rows affected: %w", err)
	}

	if rowsAffected == 0 {
		a.Logger.Warnf("no rows inserted for user: %s", user.Username)
		return fmt.Errorf("no rows inserted")
	}

	a.Logger.Infof("User registered successfully: %s", user.Username)
	return nil
}

// UpdateUser implements AuthDAO.
func (a authDAO) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET username = $1, 
		    password_hash = $2, 
		    role = $3, 
		    api_key = $4, 
		    is_revoke = $5
		WHERE user_id = $6
	`

	_, err := a.DB.Exec(query,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.APIKey,
		user.IsRevoke,
		user.UserID,
	)

	if err != nil {
		a.Logger.Errorf("failed to update user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	a.Logger.Infof("User updated successfully: %s", user.Username)
	return nil
}
