package models

import "github.com/golang-jwt/jwt/v4"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// type RegisterResponse struct {
// 	Message string `json:"message"`
// 	UserID  string `json:"userId"`
// 	APIKey  string `json:"apiKey"`
// }

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message      string `json:"message,omitempty"`
	Token        string `json:"token ,omitempty"`
	RefreshToken string `json:"refreshToken ,omitempty"`
	APIKey       string `json:"apiKey,omitempty"`
	UserID       string `json:"userId ,omitempty"`
	Role         string `json:"role ,omitempty"`
}

// type RefreshRequest struct {
// 	RefreshToken string `json:"refreshToken"`
// }

// RefreshResponse is returned after token refresh.
// type RefreshResponse struct {
// 	Token        string `json:"token"`
// 	RefreshToken string `json:"refreshToken,omitempty"`
// }

type RefreshTokenRequest struct {
	Token string `json:"token"`
}

type RefreshTokenResponse struct {
	Message      string `json:"message"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// ErrorResponse is used to return errors.
// Duplicate ErrorResponse declaration removed.

// User defines a system user.
// type User struct {
// 	UserID       string `json:"userId"`
// 	Username     string `json:"username"`
// 	PasswordHash string `json:"password_hash"`
// 	Role         string `json:"role"` // e.g., "creator", "partner"
// 	APIKey       string `json:"apiKey"`
// }

type User struct {
	UserID       string `db:"user_id" json:"userId ,omitempty"`
	Username     string `db:"username" json:"username ,omitempty"`
	PasswordHash string `db:"password_hash" json:"password_hash ,omitempty"`
	Role         string `db:"role" json:"role ,omitempty"`
	APIKey       string `db:"api_key" json:"apiKey ,omitempty"`
	IsRevoke     bool   `db:"is_revoke" json:"is_revoke ,omitempty"`
}

// CustomClaims defines custom JWT claims.
type CustomClaims struct {
	UserID string   `json:"userId"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// type ErrorResponse struct {
// 	Error string `json:"error"`
// }
