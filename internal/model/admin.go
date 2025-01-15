package model

import "time"

// 공지사항
type Notice struct {
	ID          int       `json:"NOTICE_ID"`
	Title       string    `json:"TITLE"`
	Content     string    `json:"CONTENT"`
	TopYn       bool      `json:"TOP_YN"`
	CreatedDate time.Time `json:"CREATED_DATE"`
	UpdatedDate time.Time `json:"UPDATED_DATE"`
	DelYn       bool      `json:"DEL_YN"`
}

// 책 등록용
type AddBooks struct {
	BookID         int            `json:"BOOK_ID"`
	BookTitle      string         `json:"BOOK_TITLE"`
	Publisher      string         `json:"PUBLISHER"`
	Author         string         `json:"AUTHOR"`
	PublishDate    string         `json:"PUBLISH_DATE"`
	Price          int            `json:"PRICE"`
	ContentsList   string         `json:"CONTENTS_LIST"`
	CoverURL       string         `json:"COVER_URL"`
	Summary        string         `json:"SUMMARY"`
	Description    string         `json:"DESCRIPTION"`
	Recommendation string         `json:"RECOMMENDATION"`
	Tags           string         `json:"TAGS"`
	Grade          string         `json:"GRADE"`
	ISBN           string         `json:"ISBN"`
	DelYN          int            `json:"DEL_YN"`
	IsbnAdd        string         `json:"ISBN_ADD"`
	Page           int            `json:"PAGE"`
	EbookYN        string         `json:"EBOOK_YN"`
	DevRecommends  []DevRecommend `json:"DEV_RECOMMENDS"`
}

// 개발자 등록용
type AddDevs struct {
	DevID         string `json:"ID"`
	DevName       string `json:"DEV_NAME"`
	DevDetailName string `json:"DEV_DETAIL_NAME"`
	DevHistory    string `json:"DEV_HISTORY"`
	ViewYN        bool   `json:"VIEW_YN"`
	DelYN         bool   `json:"DEL_YN"`
}

type DevRecommend struct {
	DevID           string `json:"DEV_ID"`
	RecommendReason string `json:"DEV_RECOMMEND_REASON"`
	DevName         string `json:"DEV_NAME"`
}
