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
	}
	BookDetail := b.repo.GetBookDetail(BookID)
	user := c.MustGet("user")
	c.HTML(http.StatusOK, "category_book_detail.html", gin.H{
		"BookDetail": BookDetail[0],
		"user":       user,
	})
}
