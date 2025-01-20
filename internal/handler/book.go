package handler

import (
	"DBP/internal/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	repo repository.Repository // 공통 인터페이스 사용

}

func NewBookHandler(repo repository.Repository) *BookHandler {
	return &BookHandler{
		repo: repo,
	}
}

func (b *BookHandler) GetBookDetail(c *gin.Context) {
	BookID, err := strconv.Atoi(c.Param("BookID"))
	if err != nil {
		log.Print("책 정보가 잘못되었습니다")
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "잘못된 책 ID입니다",
		})
		return
	}

	BookDetail := b.repo.GetBookDetail(BookID)
	if len(BookDetail) == 0 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "책을 찾을 수 없습니다",
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		user = nil
	}

	// 데이터를 명확하게 구조화
	c.HTML(http.StatusOK, "category_book_detail.html", gin.H{
		"BookDetail": BookDetail[0], // 단일 책 정보
		"user":       user,          // 사용자 정보
	})
}
