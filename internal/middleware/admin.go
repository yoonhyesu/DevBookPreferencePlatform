package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware 관리자 권한 확인
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// context에서 user 정보 가져오기
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "인증 정보가 없습니다",
			})
			c.Abort()
			return
		}

		// user 정보를 map으로 타입 변환
		userData, ok := user.(gin.H)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "사용자 정보 형식이 잘못되었습니다",
			})
			c.Abort()
			return
		}

		// IsAdmin 필드 확인
		isAdmin, exists := userData["IsAdmin"]
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "관리자 권한이 필요합니다",
			})
			c.Abort()
			return
		}

		// 관리자 권한이 확인되면 다음 핸들러로 진행
		c.Next()
	}
}
