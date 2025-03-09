package handler

import (
	"DBP/internal/model"
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

func (b *BookHandler) GetBookContentsList(c *gin.Context) {
	// 책 아이디 검사
	BookID, err := strconv.Atoi(c.Param("BookID"))
	if err != nil {
		log.Print("책 정보가 잘못되었습니다")
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "잘못된 책 ID입니다",
		})
		return
	}
	ContentList := b.repo.GetBookContentsList(BookID)
	c.JSON(http.StatusOK, gin.H{
		"contents": ContentList,
	})
}

// 책 좋아요 상태 확인
func (b *BookHandler) GetLikeStatus(c *gin.Context) {
	bookID := c.Query("bookId")
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "로그인이 필요합니다"})
		return
	}

	// 사용자 정보 타입 확인 및 처리
	var userId string

	// UserTemplateData 타입으로 시도
	userData, ok := user.(model.UserTemplateData)
	if ok {
		userId = userData.UserId
	} else {
		// gin.H 타입으로 시도 (chat.go에서 사용하는 방식)
		userMap, ok := user.(gin.H)
		if ok {
			userId = userMap["UserId"].(string)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 정보를 가져올 수 없습니다"})
			return
		}
	}

	// 사용자의 좋아요 목록 확인
	liked := b.repo.CheckBookLiked(userId, bookID)

	c.JSON(http.StatusOK, gin.H{
		"liked": liked,
	})
}

// 책 좋아요 토글
func (b *BookHandler) ToggleLike(c *gin.Context) {
	bookID := c.PostForm("bookId")
	liked, _ := strconv.ParseBool(c.PostForm("liked"))

	log.Printf("ToggleLike 요청: bookID=%s, liked=%v", bookID, liked)

	user, exists := c.Get("user")
	if !exists {
		log.Println("사용자 정보 없음")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "로그인이 필요합니다"})
		return
	}

	log.Printf("사용자 정보 타입: %T, 값: %+v", user, user)

	// 사용자 정보 타입 확인 및 처리
	var userId string

	// UserTemplateData 타입으로 시도
	userData, ok := user.(model.UserTemplateData)
	if ok {
		userId = userData.UserId
		log.Printf("UserTemplateData 타입으로 변환 성공: userId=%s", userId)
	} else {
		// gin.H 타입으로 시도 (chat.go에서 사용하는 방식)
		userMap, ok := user.(gin.H)
		if ok {
			userId = userMap["UserId"].(string)
			log.Printf("gin.H 타입으로 변환 성공: userId=%s", userId)
		} else {
			log.Printf("사용자 정보 변환 실패: %T", user)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 정보를 가져올 수 없습니다"})
			return
		}
	}

	// 좋아요 상태 업데이트
	err := b.repo.UpdateBookLike(userId, bookID, liked)
	if err != nil {
		log.Printf("좋아요 상태 업데이트 실패: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "좋아요 상태 업데이트 실패"})
		return
	}

	log.Printf("좋아요 상태 업데이트 성공: userId=%s, bookID=%s, liked=%v", userId, bookID, liked)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// 좋아요한 책 목록 페이지
func (b *BookHandler) GetLikedBooks(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 사용자 정보 타입 확인 및 처리
	var userId string

	// UserTemplateData 타입으로 시도
	userData, ok := user.(model.UserTemplateData)
	if ok {
		userId = userData.UserId
	} else {
		// gin.H 타입으로 시도 (chat.go에서 사용하는 방식)
		userMap, ok := user.(gin.H)
		if ok {
			userId = userMap["UserId"].(string)
		} else {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "사용자 정보를 가져올 수 없습니다",
			})
			return
		}
	}

	// 사용자가 좋아요한 책 목록 가져오기
	likedBooks := b.repo.GetLikedBooks(userId)

	c.HTML(http.StatusOK, "like_book.html", gin.H{
		"user":       user,
		"LikedBooks": likedBooks,
	})
}
