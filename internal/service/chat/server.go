package chat

import (
	"DBP/internal/model"
	"DBP/pkg/database"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type ChatServer struct {
	db          *database.MariaDBHandler
	redisClient *database.RedisHandler
	rooms       map[int]*Room // 채팅방 관리 :  o(1) 시간복잡도=>bookid로 채팅방 검색 빨리하려고 사용함
	mutex       sync.RWMutex  // 동시성 제어
}

type Room struct {
	bookID    int                     // 책 ID
	clients   map[*Client]bool        // 접속 중인 클라이언트
	broadcast chan *model.ChatMessage // 메시지 전달 채널
	mutex     sync.RWMutex            // 동시성 제어
}

func NewChatServer(db *database.MariaDBHandler, redisClient *database.RedisHandler) *ChatServer {
	return &ChatServer{
		db:          db,
		redisClient: redisClient,
		rooms:       make(map[int]*Room),
		mutex:       sync.RWMutex{},
	}
}

func (s *ChatServer) GetOrCreateRoom(bookID int) *Room {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if room, exists := s.rooms[bookID]; exists { // 채팅방 찾기
		return room
	}
	room := &Room{
		bookID:    bookID,
		clients:   make(map[*Client]bool),
		broadcast: make(chan *model.ChatMessage),
	}
	s.rooms[bookID] = room

	// 채팅방 메시지 처리 고루틴
	go room.run(s)
	return room
}

func (r *Room) run(s *ChatServer) {
	for msg := range r.broadcast {
		// db에 메시지 저장
		if err := s.saveMessage(msg); err != nil {
			log.Printf("메시지 저장 실패: %v", err)
			continue
		}

		// redis에 메시지 발행
		if err := s.publishMessage(msg); err != nil {
			log.Printf("메시지 발행 실패: %v", err)
		}
		// 모든 클라이언트에게 메시지 전송(같은채널)
		r.mutex.Lock()
		for client := range r.clients {
			select { // non-blocking 방식
			case client.send <- msg: // client.send : 각 클라이언트마다 갖고이는 메시지 전송용 채널
			default: // 전송 실패 시..
				close(client.send)
				delete(r.clients, client)
			}
		}
		r.mutex.Unlock() // 여러 고루틴이 동시에 클라이언트 맵을 수정하는것 방지
	}
}

// mariadb에 메시지저장
func (s *ChatServer) saveMessage(msg *model.ChatMessage) error {
	_, err := s.db.Connection.Exec(`
		INSERT INTO book_chat_messages
		(BOOK_ID, USER_ID, USER_NAME, MESSAGE, CREATED_AT)
		VALUES (?, ?, ?, ?, ?)`,
		msg.BookID, msg.UserID, msg.Username, msg.Message, msg.CreatedAt)
	return err
}

// 채팅메시지 redis pub/sub시스템을 통해 발행(publish)
func (s *ChatServer) publishMessage(msg *model.ChatMessage) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.redisClient.Client.Publish(context.Background(),
		fmt.Sprintf("chat:%d", msg.BookID), // 채널이름 chat:{bookID}형식으로 생성
		string(msgJSON)).Err()
}

// Room에 클라이언트 추가 메서드
func (r *Room) AddClient(client *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.clients[client] = true
}

// GetMessages 채팅 메시지 조회
func (s *ChatServer) GetMessages(bookID int) ([]model.ChatMessage, error) {
	rows, err := s.db.Connection.Query(`
		SELECT ID, BOOK_ID, USER_ID, USER_NAME, MESSAGE, CREATED_AT
		FROM book_chat_messages
		WHERE BOOK_ID = ?
		ORDER BY CREATED_AT DESC
		LIMIT 100
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.ChatMessage
	for rows.Next() {
		var msg model.ChatMessage
		err := rows.Scan(
			&msg.ID,
			&msg.BookID,
			&msg.UserID,
			&msg.Username,
			&msg.Message,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
