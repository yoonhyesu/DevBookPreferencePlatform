package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// 비밀번호 해싱
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 비밀번호 검증 (통합 버전)
func ComparePassword(inputPW, hashedPW string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(inputPW))
	if err != nil {
		return fmt.Errorf("비밀번호가 일치하지 않습니다")
	}
	return nil
}
