package models

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	UserID string   `json:"userId"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}


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

type RefreshTokenRequest struct {
	Token string `json:"token"`
}

type RefreshTokenResponse struct {
	Message      string `json:"message"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type User struct {
	UserID       string `db:"user_id" json:"userId ,omitempty"`
	Username     string `db:"username" json:"username ,omitempty"`
	PasswordHash string `db:"password_hash" json:"password_hash ,omitempty"`
	Role         string `db:"role" json:"role ,omitempty"`
	APIKey       string `db:"api_key" json:"apiKey ,omitempty"`
	IsRevoke     bool   `db:"is_revoke" json:"is_revoke ,omitempty"`
}
