package model

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}
