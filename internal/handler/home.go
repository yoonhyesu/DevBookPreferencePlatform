package handler

import (
	"DBP/internal/repository"
	service "DBP/internal/service/rag_book_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler 로봇 설계도: repo는 로봇이 가지고 있어야할 도구(망치)임
type HomeHandler struct {
	repo       repository.Repository
	ragService *service.RagServer
}

// 하나의 통합된 생성자 함수
func NewHomeHandler(repo repository.Repository, ragService *service.RagServer) *HomeHandler {
	return &HomeHandler{
		repo:       repo,
		ragService: ragService,
	}
}

// 공지사항 조회
// 리시버 사용 "이 기능은 이 로봇이 수행할 수 있는 동작 가리킴"
func (h *HomeHandler) GetHomeData(c *gin.Context) {
	notices := h.repo.GetNotice()
	devBook := h.repo.GetHomeDevBook()
	books := h.ragService.GetCasheBook()

	// templateData 대신 새로운 데이터 맵 생성
	templateData := gin.H{
		"user":        c.MustGet("user"), // 미들웨어에서 설정한 user 데이터
		"notices":     notices,
		"homeDevBook": devBook,
		"latestBooks": books.Documents,
	}

	// log.Printf("홈 화면 템플릿 데이터: %+v", templateData)
	c.HTML(http.StatusOK, "index.html", templateData)
}
