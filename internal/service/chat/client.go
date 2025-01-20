package chat

import (
	"DBP/internal/model"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn // 웹 소켓 연결
	userID   string
	username string
	send     chan *model.ChatMessage // 메시지 전송 채널
	room     *Room
}

func NewClient(conn *websocket.Conn, userID, username string, room *Room) *Client {
	return &Client{
		conn:     conn,
		userID:   userID,
		username: username,
		send:     make(chan *model.ChatMessage, 256),
		room:     room,
	}
}

// 실제 통신-클라이언트가 websocker conn으로 보내는 메시지를 처리함
// 메시지 수신 담당 (웹소켓을 통해 브라우저로부터 메시지 받음)
// 받은 메시지를 ChatMessage구조체로 변환
// 사용자 정보(ID,이름)을 메시지에 추가
// Room의 broadcastMessage를 호출하여 메시지 전파
func (c *Client) ReadPump() {
	defer func() {
		c.room.mutex.Lock()
		delete(c.room.clients, c)
		c.room.mutex.Unlock()
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var chatMsg model.ChatMessage
		if err := json.Unmarshal(message, &chatMsg); err != nil {
			log.Printf("error: %v", err)
			continue
		}

		chatMsg.UserID = c.userID
		chatMsg.Username = c.username
		chatMsg.BookID = c.room.bookID
		chatMsg.CreatedAt = time.Now()

		c.room.broadcast <- &chatMsg
	}
}

// send로 들어오는 []byte메시지를 처리
// 메시지 전송 담당
// send채널을 통해 메시지 받음
// 받은 메시지를 통해 웹소켓을 통해 브라우저로 전송
func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		}
	}
}
