package middleware

import (
	"DBP/config"
	"DBP/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// 로그인한 사용자 세션 유지
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 쿠키에서 access token 추출
		accessToken, err := c.Cookie("access_token")
		log.Printf("쿠키에서 추출한 토큰: %v", accessToken)

		if err != nil {
			log.Printf("쿠키 없음: %v", err)
			c.Set("user", nil)
			c.Next()
			return
		}

		// 토큰 검증
		claims, err := utils.ValidateToken(accessToken)
		if err != nil {
			log.Printf("토큰 검증 실패: %v", err)
			c.SetCookie("access_token", "", -1, "/", config.GetConfig().Domain, false, true)
			c.Set("user", nil)
			c.Next()
			return
		}

		// 사용자 데이터 전달
		userData := gin.H{
			"UserId":   claims["user_id"],
			"Email":    claims["email"],
			"Username": claims["user_name"],
			"IsAdmin":  claims["is_admin"],
		}
		log.Printf("전달되는 사용자 데이터: %+v", userData)

		c.Set("user", userData)
		c.Next()
	}
}
