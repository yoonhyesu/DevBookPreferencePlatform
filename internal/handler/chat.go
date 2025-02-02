package handler

import (
	"DBP/internal/service/chat"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	chatServer *chat.ChatServer
}

// 웹 소켓 연결 설정
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			os.Getenv("ALLOWED_ORIGIN1"),
			os.Getenv("ALLOWED_ORIGIN2"),
			os.Getenv("ALLOWED_ORIGIN3"),
		}
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}
		log.Printf("거부된 Origin: %s", origin)
		return false
	},
}

func NewChatHandler(chatServer *chat.ChatServer) *ChatHandler {
	return &ChatHandler{
		chatServer: chatServer,
	}
}

// 웹 소켓 연결 처리
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	log.Printf("웹소켓 연결 시도 - Headers: %v", c.Request.Header)

	// 사용자 인증 확인
	user, exists := c.Get("user")
	if !exists {
		log.Printf("인증되지 않은 사용자의 웹소켓 연결 시도")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
		return
	}
	userData := user.(gin.H)

	// BookID 파라미터 확인
	bookID, err := strconv.Atoi(c.Query("BOOK_ID"))
	if err != nil {
		log.Printf("잘못된 BookID: %v", c.Query("BOOK_ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 책 id입니다"})
		return
	}

	// 웹 소켓 연결 수립
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("웹소켓 연결 실패 - Error: %v, Headers: %v", err, c.Request.Header)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "웹소켓연결에 실패했습니다"})
		return
	}

	log.Printf("웹소켓 연결 성공 - User: %v, BookID: %d", user, bookID)
	// 채팅방 가져오기 또는 생성
	room := h.chatServer.GetOrCreateRoom(bookID)

	// 클라이언트 생성
	client := chat.NewClient(
		conn,
		userData["UserId"].(string),
		userData["Username"].(string),
		room,
	)

	// 클라이언트를 채팅방에 추가
	room.AddClient(client)

	// 고루틴으로 메시지 송수신 처리 시작
	go client.ReadPump()
	go client.WritePump()
}

// GetChatMessages 채팅 메시지 조회
func (h *ChatHandler) GetChatMessages(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 책 ID입니다"})
		return
	}

	messages, err := h.chatServer.GetMessages(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "메시지 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
