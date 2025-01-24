package handler

import (
	"DBP/internal/model"
	"DBP/internal/repository"
	"DBP/internal/service"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	repo repository.Repository
}

func NewAdminHandler(repo repository.Repository) *AdminHandler {
	return &AdminHandler{
		repo: repo,
	}
}

// 개발자 목록 조회 API
func (a *AdminHandler) GetDevList(c *gin.Context) {
	devs := a.repo.GetDevList()
	c.JSON(http.StatusOK, devs)
}

// 개발자 프로필 이미지 환경변수 가져오기
func getDevImagePath() string {
	//return "/docker/dbpstorage/image/dev"
	return "/storage/image/dev"
}

// 개발자 등록 API
func (a *AdminHandler) AddDevs(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "폼 데이터 처리 중 오류가 발생했습니다"})
		return
	}

	var dev model.AddDevs

	// 폼 데이터 바인딩
	if values, exists := form.Value["DEV_NAME"]; exists && len(values) > 0 {
		dev.DevName = values[0]
	}
	if values, exists := form.Value["DEV_DETAIL_NAME"]; exists && len(values) > 0 {
		dev.DevDetailName = values[0]
	}
	if values, exists := form.Value["DEV_HISTORY"]; exists && len(values) > 0 {
		dev.DevHistory = values[0]
	}
	if values, exists := form.Value["VIEW_YN"]; exists && len(values) > 0 {
		viewYN, err := strconv.ParseBool(values[0])
		if err == nil {
			dev.ViewYN = viewYN
		}
	}

	// 이미지 처리
	if files := form.File["PROFILE_IMAGE"]; len(files) > 0 {
		file := files[0]
		ext := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + ext

		// DB에는 파일명만 저장 (앞에 / 포함)
		dev.ProfileImagePath = "/" + newFileName

		// 실제 파일 저장 (개발자 프로필 디렉토리에)
		uploadPath := filepath.Join(getDevImagePath(), newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "이미지 저장 실패"})
			return
		}

		if err := a.repo.AddDev(dev); err != nil {
			log.Printf("개발자 프로필 수정 실패: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"messsage":           "프로필이 성공적으로 수정되었습니다",
			"PROFILE_IMAGE_PATH": dev.ProfileImagePath,
		})
	}
}

// 개발자 수정
func (a *AdminHandler) UpdateDevs(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "폼 데이터 처리 중 오류가 발생했습니다"})
		return
	}

	var dev model.AddDevs

	// 폼 데이터 바인딩
	if values, exists := form.Value["DEV_NAME"]; exists && len(values) > 0 {
		dev.DevName = values[0]
	}
	if values, exists := form.Value["DEV_DETAIL_NAME"]; exists && len(values) > 0 {
		dev.DevDetailName = values[0]
	}
	if values, exists := form.Value["DEV_HISTORY"]; exists && len(values) > 0 {
		dev.DevHistory = values[0]
	}
	if values, exists := form.Value["VIEW_YN"]; exists && len(values) > 0 {
		viewYN, err := strconv.ParseBool(values[0])
		if err == nil {
			dev.ViewYN = viewYN
		}
	}
	if values, exists := form.Value["DEV_ID"]; exists && len(values) > 0 {
		dev.DevID = values[0]

	}

	// 이미지 처리
	if files := form.File["PROFILE_IMAGE"]; len(files) > 0 {
		file := files[0]
		ext := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + ext

		// DB에는 파일명만 저장 (앞에 / 포함)
		dev.ProfileImagePath = "/" + newFileName

		// 실제 파일 저장 (개발자 프로필 디렉토리에)
		uploadPath := filepath.Join(getDevImagePath(), newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "이미지 저장 실패"})
			return
		}
	} else {
		// 이미지가 업로드되지 않은 경우 기존 이미지 경로 가져오기
		existingDev, err := a.repo.GetDevID(dev.DevID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "기존 개발자 정보 조회 중 오류 발생"})
			return
		}
		dev.ProfileImagePath = existingDev.ProfileImagePath
	}

	if err := a.repo.UpdateDev(dev); err != nil {
		log.Printf("개발자 프로필 수정 실패: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "개발자 정보가 성공적으로 수정되었습니다",
		"PROFILE_IMAGE_PATH": dev.ProfileImagePath,
	})
}

// 개발자 삭제
func (a *AdminHandler) DeleteDevs(c *gin.Context) {
	var req struct {
		ID string `json:"ID"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.DeleteDev(req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "개발자가 성공적으로 삭제되었습니다"})
}

// 공지사항 목록 조회
func (a *AdminHandler) GetNoticeList(c *gin.Context) {
	notices := a.repo.GetNoticeList()
	log.Printf("조회된 공지사항: %+v", notices)
	c.JSON(http.StatusOK, notices)
}

// 공지사항 등록
func (a *AdminHandler) AddNotice(c *gin.Context) {
	var notice model.Notice
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.AddNotice(notice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "공지사항이 성공적으로 등록되었습니다"})
}

// 공지사항 수정
func (a *AdminHandler) UpdateNotice(c *gin.Context) {
	var notice model.Notice
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.UpdateNotice(notice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "공지사항이 성공적으로 수정되었습니다"})
}

// 공지사항 삭제
func (a *AdminHandler) DeleteNotice(c *gin.Context) {
	var notice model.Notice
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.DeleteNotice(notice.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "공지사항이 성공적으로 삭제되었습니다"})
}

// 책 목록 조회 API
func (a *AdminHandler) GetBookList(c *gin.Context) {
	books := a.repo.GetBookList()
	log.Printf("조회된 책 목록: %+v", books)
	c.JSON(http.StatusOK, books)
}

type BookAddRequest struct {
	Book        model.AddBooks       `json:"book"`
	DevContents []model.DevRecommend `json:"DEV_CONTENTS"`
}

// 책 등록 API
func (a *AdminHandler) AddBook(c *gin.Context) {
	var req BookAddRequest

	// 요청 바디 로깅
	body, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	log.Printf("받은 요청 바디: %s", string(body))

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON 바인딩 오류: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.AddBook(req.Book, req.DevContents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "책이 성공적으로 등록되었습니다"})
}

// 책 수정 API
func (a *AdminHandler) UpdateBook(c *gin.Context) {
	var req struct {
		Book          model.AddBooks       `json:"book"`
		DevRecommends []model.DevRecommend `json:"devRecommends"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.UpdateBook(req.Book, req.DevRecommends); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "책이 성공적으로 수정되었습니다"})
}

// 책 삭제 API
func (a *AdminHandler) DeleteBook(c *gin.Context) {
	var req struct {
		ID int `json:"BOOK_ID"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.DeleteBook(req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "책이 성공적으로 삭제되었습니다"})
}

// 태그 목록 조회 API
func (a *AdminHandler) GetTagList(c *gin.Context) {
	tags := a.repo.GetAdminTagList()
	log.Printf("조회된 태그 목록: %+v", tags)
	c.JSON(http.StatusOK, tags)
}

// 책 검색 API
func (a *AdminHandler) SearchBooks(c *gin.Context) {
	title := strings.TrimSpace(c.Query("title"))
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "검색어를 입력해주세요"})
		return
	}

	bookService := service.NewBookSearchService()
	result, err := bookService.SearchBooks(title)
	if err != nil {
		log.Printf("도서 검색 중 오류 발생: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
