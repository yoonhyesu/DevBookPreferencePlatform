package handler

import (
	"DBP/internal/model"
	"DBP/internal/repository"
	"DBP/internal/service"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

// 개발자 등록 API
func (a *AdminHandler) AddDevs(c *gin.Context) {
	var dev model.AddDevs
	if err := c.ShouldBindJSON(&dev); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.AddDev(dev); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "개발자가 성공적으로 등록되었습니다"})
}

// 개발자 수정
func (a *AdminHandler) UpdateDevs(c *gin.Context) {
	var dev model.AddDevs
	if err := c.ShouldBindJSON(&dev); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.repo.UpdateDev(dev); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "개발자 정보가 성공적으로 수정되었습니다"})
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
