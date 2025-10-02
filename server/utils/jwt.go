package utils

import (
	"time"

	"server/model"

	"github.com/golang-jwt/jwt/v4"
)

// CreateToken 创建Token
func CreateToken(userID int, username string, email string, isAdmin bool) (string, error) {
	// hmac的密钥类型是字节数组
	secret := []byte("image_hub_secret")

	claims := model.Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // 1小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Server",
			Subject:   "user_token",
		},
	}

	// 创建Token，使用HS256算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整令牌
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*model.Claims, error) {
	secret := []byte("image_hub_secret")

	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// RefreshToken 刷新Token
func RefreshToken(tokenString string) (string, error) {
	// 先解析原Token
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 创建新Token（延长过期时间）
	return CreateToken(claims.UserID, claims.Username, claims.Email, claims.IsAdmin)
}
