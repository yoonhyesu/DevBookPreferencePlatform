package model

// 메인화면 노출 개발자 정보
type HomeDevBook struct {
	DevID            string `json:"ID"`
	DevName          string `json:"DEV_NAME"`
	DevHistory       string `json:"DEV_HISTORY"`
	DevDetailName    string `json:"DEV_DETAIL_NAME"`
	ProfileImagePath string `json:"PROFIILE_IMAGE_PATH"`
	Books            []Book
}

// 메인화면 노출 책 정보
type Book struct {
	BookID    int    `json:"BOOK_ID"`
	BookTitle string `json:"BOOK_TITLE"`
	CoverURL  string `json:"COVER_URL"`
}
