package model

import "time"

type ChatMessage struct {
	ID        int64     `json:"ID"`
	BookID    int       `json:"BOOK_ID"`
	UserID    string    `json:"USER_ID"`
	Username  string    `json:"USER_NAME"`
	Message   string    `json:"MESSAGE"`
	CreatedAt time.Time `json:"CREATED_AT"`
}

type ChatRoom struct {
	BookID    int    `json:"BOOK_ID"`
	BookTitle string `json:"BOOK_TITLE"`
}
