package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func GenerateTokens(userID, email string, isAdmin bool, userName string) (*TokenDetails, error) {
	td := &TokenDetails{
		RefreshUuid: uuid.NewString(),
		AtExpires:   time.Now().Add(30 * time.Minute).Unix(),
		RtExpires:   time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	// Access Token 생성
	atClaims := jwt.MapClaims{
		"user_id":   userID,
		"email":     email,
		"is_admin":  isAdmin,
		"user_name": userName,
		"exp":       td.AtExpires,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	// Refresh Token 생성
	rtClaims := jwt.MapClaims{
		"refresh_uuid": td.RefreshUuid,
		"user_id":      userID,
		"exp":          td.RtExpires,
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// 토큰 만료 시간 검증
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid token expiration")
	}

	if float64(time.Now().Unix()) > exp {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
