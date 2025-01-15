package handler

import (
	"DBP/internal/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecommendHandler struct {
	repo repository.Repository // 공통 인터페이스 사용

}

func NewRecommendHandler(repo repository.Repository) *RecommendHandler {
	return &RecommendHandler{
		repo: repo,
	}
}

// Tag정보 가져옴
func (r *RecommendHandler) GetTagHandler(c *gin.Context) {
	Tags := r.repo.GetTagList()
	user := c.MustGet("user")
	c.HTML(http.StatusOK, "category_field.html", gin.H{
		"Tags": Tags,
		"user": user,
	})
}

// Tag에 맞는 책 가져옴
func (r *RecommendHandler) GetTagWithBookHandler(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("TagID"))
	if err != nil {
		log.Print("잘못된 태그입니다")
	}
	TagWithBook := r.repo.GetTagWithBookList(tagID)
	user := c.MustGet("user")
	c.HTML(http.StatusOK, "category_field_detail.html", gin.H{
		"TagWithBooks": TagWithBook,
		"user":         user,
	})
}

func (r *RecommendHandler) GetRecommendDevListHandler(c *gin.Context) {
	DevList := r.repo.GetRecommendDevList()
	user := c.MustGet("user")
	c.HTML(http.StatusOK, "category_developer.html", gin.H{
		"DevList": DevList,
		"user":    user,
	})
}

// 개발자가 추천도서 가져옴
func (r *RecommendHandler) GetDevWithBooksHandler(c *gin.Context) {
	devID := c.Param("devID")
	if devID == "" {
		log.Print("개발자 ID가 비어있습니다")
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "잘못된 개발자 ID입니다",
		})
		return
	}

	// 개발자 정보 조회
	developers := r.repo.GetDevInfo(devID)
	if len(developers) == 0 {
		log.Print("개발자 정보가 존재하지 않습니다")
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "개발자를 찾을 수 없습니다",
		})
		return
	}

	// 추천 도서 목록 조회
	books := r.repo.GetDevRecommendedBooks(devID)
	user := c.MustGet("user")
	c.HTML(http.StatusOK, "category_developer_detail.html", gin.H{
		"Developer":        developers[0],
		"RecommendedBooks": books,
		"user":             user,
	})
}

// 단계별 추천
func (r *RecommendHandler) GetStepWithBookHandler(c *gin.Context) {
	step := c.Param("stepName")
	StepWithBook := r.repo.GetStepWithBook(step)
	user := c.MustGet("user")
	c.HTML(http.StatusOK, "category_step.html", gin.H{
		"StepWithBook": StepWithBook,
		"Grade":        step,
		"user":         user,
	})
}
