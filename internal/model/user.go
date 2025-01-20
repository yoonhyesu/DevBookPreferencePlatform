package model

// 사용자
type User struct {
	UserId           string `json:"USER_ID"`
	Password         string `json:"PASSWORD"`
	UserName         string `json:"USER_NAME"`
	Email            string `json:"EMAIL"`
	PhoneNumber      string `json:"PHONE_NUMBER"`
	GithubLink       string `json:"GITHUB_LINK"`
	BlogLink         string `json:"BLOG_LINK"`
	Company          string `json:"COMPANY"`
	UserText         string `json:"USER_TEXT"`
	IsAdmin          bool   `json:"IS_ADMIN"`
	ProfileImagePath string `json:"PROFILE_IMAGE_PATH"`
	ProfileImageName string `json:"PROFILE_IMAGE_NAME"`
}

// 회원가입용
type RegisterRequest struct {
	UserID      string `json:"USER_ID"`
	Password    string `json:"PASSWORD"`
	UserName    string `json:"USER_NAME"`
	Email       string `json:"EMAIL"`
	PhoneNumber string `json:"PHONE_NUMBER"`
}

// 사용자계정 관리 비밀번호변경
type UserSettingRequest struct {
	CurrentPassword string `json:"CURRENT_PASSWORD"`
	NewPassword     string `json:"NEW_PASSWORD"`
}

// 로그인용
type LoginRequest struct {
	UserId   string `json:"USER_ID"`
	Password string `json:"PASSWORD"`
	Email    string `json:"EMAIL"`
	IsAdmin  bool   `json:"IS_ADMIN"`
	UserName string `json:"USER_NAME"`
}

// 템플릿에 전달할 사용자 정보
type UserTemplateData struct {
	UserId   string `json:"USER_ID"`
	Email    string `json:"EMAIL"`
	UserName string `json:"USER_NAME"`
	IsAdmin  bool   `json:"IS_ADMIN"`
}

// 프로필 수정 요청
type ProfileEditRequest struct {
	UserID           string `json:"-"` // 서버에서 설정
	UserName         string `json:"USER_NAME"`
	PhoneNumber      string `json:"PHONE_NUMBER"`
	UserText         string `json:"USER_TEXT"`
	Company          string `json:"COMPANY"`
	GithubLink       string `json:"GITHUB_LINK"`
	BlogLink         string `json:"BLOG_LINK"`
	ProfileImagePath string `json:"PROFILE_IMAGE_PATH"`
}

// 비밀번호 변경 요청
type PasswordChangeRequest struct {
	UserID      string `json:"-"` // 서버에서 설정
	OldPassword string `json:"OLD_PASSWORD"`
	NewPassword string `json:"NEW_PASSWORD"`
}

// 회원 탈퇴 요청
type LeaveRequest struct {
	UserID   string `json:"USER_ID"`
	Password string `json:"PASSWORD"`
}
