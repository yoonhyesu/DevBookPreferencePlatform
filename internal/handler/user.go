package handler

import (
	"DBP/config"
	"DBP/internal/model"
	"DBP/internal/repository"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	repo repository.Repository
}

func NewUserHandler(repo repository.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// 회원가입 시 아이디 중복체크
func (u *UserHandler) CheckDuplicateID(c *gin.Context) {
	var req struct {
		UserID string `json:"USER_ID"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON 바인딩 에러: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}

	registerReq := model.RegisterRequest{
		UserID: req.UserID,
	}

	// 레포지토리 호출
	err := u.repo.CheckDuplicateID(registerReq)
	if err != nil {
		if err.Error() == "이미 존재하는 아이디 입니다" {
			// 중복된 아이디인 경우
			c.JSON(http.StatusOK, gin.H{"exists": true})
			return
		}
		// 그 외 에러의 경우
		log.Printf("예상치 못한 에러: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 내부 오류가 발생했습니다"})
		return
	}

	// 사용 가능한 아이디인 경우
	c.JSON(http.StatusOK, gin.H{"exists": false})
}

// 회원가입 처리
func (u *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}
	// 레포지토리 호출
	if err := u.repo.Register(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "회원가입이 완료되었습니다"})
}

// 로그인
func (u *UserHandler) SignIn(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}

	// 로그인 검증 및 토큰 생성
	tokens, err := u.repo.SignIn(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Access Token을 쿠키에 저장
	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		30*60,
		"/",
		config.GetConfig().Domain,
		false,
		true,
	)

	log.Printf("로그인 성공 - 사용자: %s, 토큰: %s", req.UserId, tokens.AccessToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "로그인 성공",
		"user": gin.H{
			"user_id":  req.UserId,
			"email":    req.Email,
			"username": req.UserName,
			"is_admin": req.IsAdmin,
		},
	})
}

// 로그아웃
func (u *UserHandler) SignOut(c *gin.Context) {
	// 쿠키에서 access token 가져오기
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "로그인 상태가 아님"})
		return
	}

	// redis에 토큰 블랙리스트 추가-같은것으로 발행x
	if err := u.repo.SignOut(accessToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 쿠키삭제
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		config.GetConfig().Domain,
		false,
		true,
	)
	c.JSON(http.StatusOK, gin.H{"message": "로그아웃되었습니다"})
}

// 프로필 조회
func (u *UserHandler) GetUserProfile(c *gin.Context) {
	user := c.MustGet("user")
	userData, ok := user.(gin.H)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 정보 처리 오류가 발생했습니다"})
		return
	}

	userID, ok := userData["UserId"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 ID를 찾을 수 없습니다"})
		return
	}

	userInfo, err := u.repo.GetUserProfile(userID)
	if err != nil {
		log.Printf("프로필 조회 실패: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "프로필 조회에 실패했습니다"})
		return
	}

	log.Printf("프로필 조회 성공 - 데이터: %+v", userInfo)
	c.JSON(http.StatusOK, gin.H{
		"PHONE_NUMBER":       userInfo.PhoneNumber,
		"USER_TEXT":          userInfo.UserText,
		"COMPANY":            userInfo.Company,
		"GITHUB_LINK":        userInfo.GithubLink,
		"BLOG_LINK":          userInfo.BlogLink,
		"PROFILE_IMAGE_PATH": os.Getenv("SELECT_PROFILE_PATH") + "/" + userInfo.ProfileImagePath,
	})
}

// 프로필 이미지 환경변수 가져오기
func getUserProfilePath() string {
	return os.Getenv("PROFILE_PATH")
}

// 프로필 수정
func (h *UserHandler) ProfileEdit(c *gin.Context) {
	// 현재 로그인한 사용자 정보 가져오기
	user := c.MustGet("user")
	userData, ok := user.(gin.H)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 정보 처리 오류가 발생했습니다"})
		return
	}

	userID, ok := userData["UserId"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 ID를 찾을 수 없습니다"})
		return
	}

	// multiple files 파싱
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "폼 데이터 처리 중 오류가 발생했습니다"})
		return
	}

	var req model.ProfileEditRequest
	// 현재 로그인한 사용자의 ID 설정
	req.UserID = userID

	// 폼 데이터 바인딩
	if values, exists := form.Value["USER_NAME"]; exists && len(values) > 0 {
		req.UserName = values[0]
	}
	if values, exists := form.Value["PHONE_NUMBER"]; exists && len(values) > 0 {
		req.PhoneNumber = values[0]
	}
	if values, exists := form.Value["USER_TEXT"]; exists && len(values) > 0 {
		req.UserText = values[0]
	}
	if values, exists := form.Value["COMPANY"]; exists && len(values) > 0 {
		req.Company = values[0]
	}
	if values, exists := form.Value["GITHUB_LINK"]; exists && len(values) > 0 {
		req.GithubLink = values[0]
	}
	if values, exists := form.Value["BLOG_LINK"]; exists && len(values) > 0 {
		req.BlogLink = values[0]
	}

	// 이미지 처리
	if files := form.File["PROFILE_IMAGE"]; len(files) > 0 {
		file := files[0]
		ext := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + ext

		// DB에는 파일명만 저장
		profileImagePath := newFileName

		// 실제 파일 저장 (유저 프로필 디렉토리에)
		uploadPath := filepath.Join(getUserProfilePath(), newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "이미지 저장 실패"})
			return
		}
		req.ProfileImagePath = profileImagePath
	} else {
		// 이미지가 업로드되지 않은 경우 기존 프로필 이미지 경로 가져오기
		existingProfile, err := h.repo.GetUserProfile(userID)
		if err != nil {
			log.Printf("기존 프로필 정보 조회 실패: %v", err)
		} else {
			req.ProfileImagePath = existingProfile.ProfileImagePath
		}
	}

	if err := h.repo.UpdateProfile(req); err != nil {
		log.Printf("프로필 수정 실패: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "프로필이 성공적으로 수정되었습니다",
		"PROFILE_IMAGE_PATH": req.ProfileImagePath,
	})
}

// 비밀번호 변경
func (h *UserHandler) ChangePW(c *gin.Context) {
	user := c.MustGet("user")
	userData := user.(gin.H)
	userID := userData["UserId"].(string)
	var req model.PasswordChangeRequest
	req.UserID = userID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.ChangePassword(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "비밀번호가 성공적으로 변경되었습니다"})
}

// 회원 탈퇴
func (h *UserHandler) LeavePlattform(c *gin.Context) {
	user := c.MustGet("user")
	userData := user.(gin.H)
	userID := userData["UserId"].(string)
	var req model.LeaveRequest
	req.UserID = userID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.LeaveAccount(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 쿠키 삭제
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		config.GetConfig().Domain,
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "회원 탈퇴가 완료되었습니다"})
}

// 현재 비밀번호 확인
func (h *UserHandler) CheckCurrentPassword(c *gin.Context) {
	// 현재 로그인한 사용자 정보 가져오기
	user, exists := c.Get("user")
	if !exists {
		log.Printf("사용자 정보를 찾을 수 없음")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "로그인이 필요합니다"})
		return
	}

	// gin.H 타입으로 변환
	userData, ok := user.(gin.H)
	if !ok {
		log.Printf("사용자 정보 타입 변환 실패")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "사용자 정보 처리 오류"})
		return
	}

	log.Printf("전달된 사용자 데이터: %+v", userData)

	// user_id 추출 (UserId 키로 시도)
	userID, ok := userData["UserId"].(string)
	if !ok {
		// 첫 번째 시도 실패시 user_id로 다시 시도
		userID, ok = userData["user_id"].(string)
		if !ok {
			log.Printf("사용자 ID 추출 실패 - 전달된 데이터: %+v", userData)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "사용자 정보 처리 오류"})
			return
		}
	}

	log.Printf("현재 로그인한 사용자 ID: %s", userID)

	// 요청에서 현재 비밀번호 가져오기
	var req struct {
		CurrentPassword string `json:"CURRENT_PASSWORD"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON 바인딩 실패: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	log.Printf("입력받은 비밀번호 길이: %d", len(req.CurrentPassword))

	// 비밀번호 확인
	err := h.repo.CheckCurrentPassword(userID, req.CurrentPassword)
	if err != nil {
		log.Printf("비밀번호 확인 실패 - 사용자: %s, 에러: %v", userID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	log.Printf("비밀번호 확인 성공 - 사용자: %s", userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
