package router

import (
	"net/http"
	"os"
	"path/filepath"

	"DBP/internal/handler"
	"DBP/internal/middleware"
	"DBP/internal/repository"
	"DBP/internal/service/chat"
	service "DBP/internal/service/rag_book_service"
	"DBP/pkg/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine, db *database.MariaDBHandler, redis *database.RedisHandler) {
	g.Use(func(c *gin.Context) {
		frontendURL := os.Getenv("FRONTEND_URL")
		c.Writer.Header().Set("Access-Control-Allow-Origin", frontendURL) // 프론트엔드 도메인
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Cookie") // 추가

		// preflight request 처리
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.Next()
	})

	// 요청 본문 크기 제한을 10MB로 설정
	//g.MaxMultipartMemory = 10 << 20 // 10 MB

	// gin의 기본 요청 크기 제한 설정
	//g.Use(func(c *gin.Context) {
	// c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10MB
	// c.Next()
	//})
	// 정적 파일 제공 설정 수정
	g.Static("/node_modules", "../view/node_modules")
	g.Static("/assets", "../view/assets")
	// 도커 볼륨에 마운트된 스토리지 경로 설정
	profilePath := filepath.Join("/app", "storage", "image", "profile")
	devPath := filepath.Join("/app", "storage", "image", "dev")
	// log.Printf("업로드 경로: %s", uploadPath)
	// g.Static("/uploads", uploadPath)
	g.Static("/storage/image/profile", profilePath)
	g.Static("/storage/image/dev", devPath)
	g.StaticFile("/favicon.ico", "../favicon.ico")

	// 템플릿 경로 설정
	g.LoadHTMLGlob("../view/**/*.html")
	// 디버깅용 로그 미들웨어
	g.Use(gin.Logger())

	// 모든 요청에 대해 인증 미들웨어 적용
	g.Use(middleware.AuthMiddleware())

	// 공통 repository 생성
	repo := repository.NewRepository(db, redis)

	// 웹 핸들러 정의
	// 홈화면(공지사항 조회,IT 신간 조회)
	ragServer := service.GetRagBookService()
	ragServer.StartAutoRefresh()
	homeHandler := handler.NewHomeHandler(repo, ragServer)
	// 회원관련 메뉴 (로그인, 회원가입)
	userHandler := handler.NewUserHandler(repo)
	// 추천관련 메뉴(분야,개발자,단계)
	recommend := handler.NewRecommendHandler(repo)
	// 책 상세정보 관련 메뉴
	bookdetail := handler.NewBookHandler(repo)
	// 관리자 메뉴 (추천책등록, 개발자등록, 공지사항등록)
	adminHandler := handler.NewAdminHandler(repo)

	// 채팅 서버 및 핸들러 초기화
	chatServer := chat.NewChatServer(db, redis)
	chatHandler := handler.NewChatHandler(chatServer)

	// 인증이 필요없는 public 라우트
	public := g.Group("/")
	{
		// 홈 화면
		public.GET("/", homeHandler.GetHomeData)
		public.GET("/about", func(c *gin.Context) {
			c.HTML(http.StatusOK, "about.html", gin.H{})
		})

		// 회원가입/로그인 관련
		public.GET("/join", func(c *gin.Context) {
			c.HTML(http.StatusOK, "join.html", gin.H{})
		})
		public.POST("/join/check-duplicate", userHandler.CheckDuplicateID)
		public.POST("/join/submit", userHandler.Register)
		public.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{})
		})
		public.POST("/auth/signin", userHandler.SignIn)

		// 책 조회 관련 public API
		public.GET("/category/tag", recommend.GetTagHandler)
		public.GET("/category/tag/:TagID", recommend.GetTagWithBookHandler)
		public.GET("/category/dev", recommend.GetRecommendDevListHandler)
		public.GET("/category/dev/:devID", recommend.GetDevWithBooksHandler)
		public.GET("/category/step/:stepName", recommend.GetStepWithBookHandler)
		public.GET("/category/book/detail/:BookID", bookdetail.GetBookDetail)
	}

	// 인증이 필요한 private 라우트
	private := g.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		// 로그아웃
		private.POST("/auth/signout", userHandler.SignOut)
		private.POST("/logout", userHandler.SignOut)

		// 책챗리뷰채팅 관련 라우트
		private.GET("/ws/chat", chatHandler.HandleWebSocket)
		private.GET("/api/chat/messages/:bookId", chatHandler.GetChatMessages)

		// 계정 관리
		account := private.Group("/account")
		{
			account.GET("/profile", func(c *gin.Context) {
				user := c.MustGet("user")
				c.HTML(http.StatusOK, "profile.html", gin.H{"user": user})
			})
			account.GET("/profile-edit", userHandler.GetUserProfile)
			account.POST("/profile-edit/manage", userHandler.ProfileEdit)
			account.GET("/profile-setting", func(c *gin.Context) {
				user := c.MustGet("user")
				c.HTML(http.StatusOK, "account.html", gin.H{"user": user})
			})
			account.POST("/profile-setting/leave", userHandler.LeavePlattform)
			account.POST("/profile-setting/changpw", userHandler.ChangePW)
			account.POST("/check_pw", userHandler.CheckCurrentPassword)
		}

		// 관리자 전용 라우트
		admin := private.Group("/admin")
		admin.Use(middleware.AdminMiddleware()) // 관리자 권한 확인 미들웨어 추가 필요
		{
			dev := admin.Group("/dev")
			{
				dev.GET("", func(c *gin.Context) {
					user := c.MustGet("user")
					c.HTML(http.StatusOK, "devs_manager.html", gin.H{
						"user": user,
					})
				})
				dev.GET("/manage", adminHandler.GetDevList)
				dev.POST("/add", adminHandler.AddDevs)
				dev.POST("/update", adminHandler.UpdateDevs)
				dev.POST("/delete", adminHandler.DeleteDevs)
			}
			book := admin.Group("/book")
			{
				book.GET("", func(c *gin.Context) {
					user := c.MustGet("user")
					c.HTML(http.StatusOK, "book_manager.html", gin.H{
						"user": user,
					})
				})
				book.GET("/manage", adminHandler.GetBookList)
				book.POST("/add", adminHandler.AddBook)
				book.POST("/update", adminHandler.UpdateBook)
				book.POST("/delete", adminHandler.DeleteBook)
				book.GET("/tags", adminHandler.GetTagList)
				book.GET("/search", adminHandler.SearchBooks)
			}
			notice := admin.Group("/notice")
			{
				notice.GET("", func(c *gin.Context) {
					user := c.MustGet("user")
					c.HTML(http.StatusOK, "notice_manager.html", gin.H{
						"user": user,
					})
				})
				notice.GET("/manage", adminHandler.GetNoticeList)
				notice.POST("/add", adminHandler.AddNotice)
				notice.POST("/update", adminHandler.UpdateNotice)
				notice.POST("/delete", adminHandler.DeleteNotice)
			}
		}
	}
}
