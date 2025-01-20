package model

// 추천관련

// 태그 관련
type TagList struct {
	TagID   int    `json:"TAG_ID"`
	TagName string `json:"TAG_NAME"`
	DelYN   bool   `json:"DEL_YN"`
}

// 태그-책 조회
type TagWithBook struct {
	TagID     int    `json:"TAG_ID"`
	TagName   string `json:"TAG_NAME"`
	CoverURL  string `json:"COVER_URL"`
	BookID    int    `json:"BOOK_ID"`
	BookTitle string `json:"BOOK_TITLE"`
}

// 책정보
type BookInfo struct {
	BookID         int    `json:"BOOK_ID"`
	BookTitle      string `json:"BOOK_TITLE"`
	Publisher      string `json:"PUBLISHER"`
	Author         string `json:"AUTHOR"`
	PublishDate    string `json:"PUBLISH_DATE"`
	Price          int    `json:"Price"`
	ContentsList   string `json:"CONTENTS_LIST"`
	CoverURL       string `json:"COVER_URL"`
	Summary        string `json:"SUMMARY"`
	Description    string `json:"DESCRIPTION"`
	Recommendation string `json:"RECOMMENDATION"`
	Tags           string `json:"TAGS"`
	TagNames       string `json:"TAG_NAMES"`
	Grade          string `json:"GRADE"`
	ISBN           string `json:"ISBN"`
	DelYN          bool   `json:"DEL_YN"`
	IsbnAdd        string `json:"ISBN_ADD"`
	Page           int    `json:"PAGE"`
	EbookYN        string `json:"EBOOK_YN"`
}

// 책을 추천한 개발자정보
type DevInfo struct {
	DevID            string `json:"ID"`
	DevName          string `json:"DEV_NAME"`
	DevDetailName    string `json:"DEV_DETAIL_NAME"`
	DevHistory       string `json:"DEV_HISTORY"`
	ProfileImagePath string `json:"PROFILE_IMAGE_PATH"`
	ViewYN           bool   `json:"VIEW_YN"`
	DelYN            bool   `json:"DEL_YN"`
}

// 개발자 상세정보
type DevDetail struct {
	DevID            string `json:"ID"`
	DevName          string `json:"DEV_NAME"`
	DevDetailName    string `json:"DEV_DETAIL_NAME"`
	DevHistory       string `json:"DEV_HISTORY"`
	ProfileImagePath string `json:"PROFILE_IMAGE_PATH"`
	ViewYN           bool   `json:"VIEW_YN"`
	DelYN            bool   `json:"DEL_YN"`
}

// 개발자추천책 정보
type DevBookLikeReason struct {
	DevReasonID      int    `json:"ID"`
	DevReason        string `json:"DEV_RECOMMEND_REASON"`
	DelYN            bool   `json:"DEL_YN"`
	BookID           int    `json:"BOOK_ID"`
	ProfileImagePath string `json:"PROFILE_IMAGE_PATH"`
	DevID            string `json:"DEV_ID"`
	BookTitle        string `json:"BOOK_TITLE"`
	CoverURL         string `json:"COVER_URL"`
}
