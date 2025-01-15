package repository

import (
	"DBP/internal/model"

	"log"
)

// 태그 조회
func (m *CommonRepo) GetTagList() []model.TagList {
	rows, err := m.mariaDB.Connection.Query("SELECT TAG_ID, TAG_NAME, DEL_YN FROM TAGS WHERE DEL_YN = FALSE")
	if err != nil {
		log.Println(err)
		return []model.TagList{}
	}
	all := []model.TagList{}
	var c1 int
	var c2 string
	var c3 bool
	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3)
		temp := model.TagList{TagID: c1, TagName: c2, DelYN: c3}
		all = append(all, temp)
	}
	log.Println("태그:", all)
	return all
}

// 태그 조회에 따른 책표지 조회
func (m *CommonRepo) GetTagWithBookList(tagID int) []model.TagWithBook {
	rows, err := m.mariaDB.Connection.Query(`
		SELECT TAG_ID, TAG_NAME, COVER_URL, BOOK_TITLE, BOOK_ID 
		FROM TAGS AS TA 
		INNER JOIN BOOK_INFOS AS BI 
		ON FIND_IN_SET(TA.TAG_ID, REPLACE(BI.TAGS, '&', ',')) > 0 
		WHERE TA.TAG_ID = ?`, tagID)
	if err != nil {
		log.Println(err)
		return []model.TagWithBook{}
	}
	all := []model.TagWithBook{}
	var c1, c5 int
	var c2, c3, c4 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5)
		if err != nil {
			log.Println("태그/책 데이터 스캔 중 오류 발생:", err)
			continue
		}
		temp := model.TagWithBook{TagID: c1, TagName: c2, CoverURL: c3, BookTitle: c4, BookID: c5}
		all = append(all, temp)
	}
	log.Println("태그관련책표지:", all)
	return all
}

// 개발자 리스트 조회
func (m *CommonRepo) GetRecommendDevList() []model.DevInfo {
	rows, err := m.mariaDB.Connection.Query(`
	SELECT ID, DEV_NAME, DEV_DETAIL_NAME, DEV_HISTORY, VIEW_YN, DEL_YN
FROM dbp.dev_infos
WHERE DEL_YN = FALSE
	`)
	if err != nil {
		log.Println(err)
		return []model.DevInfo{}
	}
	all := []model.DevInfo{}
	var c1, c2, c3, c4 string
	var c5, c6 bool
	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6)
		if err != nil {
			log.Println("개발자 데이터 스캔 중 오류 발생:", err)
			continue
		}
		temp := model.DevInfo{DevID: c1, DevName: c2, DevDetailName: c3, DevHistory: c4, ViewYN: c5, DelYN: c6}
		all = append(all, temp)
	}
	log.Println("개발자 조회:", all)
	return all
}

// 개발자 상세 정보 조회
func (m *CommonRepo) GetDevInfo(devID string) []model.DevDetail {
	row := m.mariaDB.Connection.QueryRow(`
	SELECT 
		ID, 
		DEV_NAME, 
		DEV_DETAIL_NAME, 
		DEV_HISTORY, 
		VIEW_YN, 
		DEL_YN
	FROM dbp.dev_infos
	WHERE ID = ? AND DEL_YN = FALSE`, devID)

	all := []model.DevDetail{}
	var c1, c2, c3, c4 string
	var c5, c6 bool

	err := row.Scan(&c1, &c2, &c3, &c4, &c5, &c6)
	if err != nil {
		log.Println("개발자 정보 조회 중 오류:", err)
		return []model.DevDetail{}
	}

	temp := model.DevDetail{
		DevID:         c1,
		DevName:       c2,
		DevDetailName: c3,
		DevHistory:    c4,
		ViewYN:        c5,
		DelYN:         c6,
	}
	all = append(all, temp)
	log.Println("개발자 상세정보:", all)
	return all
}

// 개발자가 추천한 책 목록 조회
func (m *CommonRepo) GetDevRecommendedBooks(devID string) []model.DevBookLikeReason {
	rows, err := m.mariaDB.Connection.Query(`
	SELECT 
		A.ID,
		A.DEV_RECOMMEND_REASON,
		B.BOOK_ID,
		B.BOOK_TITLE,
		B.COVER_URL
	FROM dbp.dev_recommends A
	INNER JOIN dbp.BOOK_INFOS B ON A.BOOK_ID = B.BOOK_ID 
	WHERE A.DEL_YN = 0 AND A.DEV_ID = ?`, devID)

	if err != nil {
		log.Println(err)
		return []model.DevBookLikeReason{}
	}

	var books []model.DevBookLikeReason
	for rows.Next() {
		var book model.DevBookLikeReason
		err = rows.Scan(&book.DevReasonID, &book.DevReason,
			&book.BookID, &book.BookTitle, &book.CoverURL)
		if err != nil {
			log.Println("책 정보 스캔 중 오류:", err)
			continue
		}
		books = append(books, book)
	}
	return books
}

// 단계별 책 정보(표지) 가져옴
func (m *CommonRepo) GetStepWithBook(step string) []model.BookInfo {
	rows, err := m.mariaDB.Connection.Query(`
SELECT BOOK_ID, BOOK_TITLE, COVER_URL, GRADE, DEL_YN
FROM dbp.book_infos
WHERE GRADE = ?
	`, step)
	if err != nil {
		log.Println(err)
		return []model.BookInfo{}
	}
	all := []model.BookInfo{}
	var bookID int
	var bookTitle, coverURL, grade string
	var delYN bool

	for rows.Next() {
		err = rows.Scan(&bookID, &bookTitle, &coverURL, &grade, &delYN)
		if err != nil {
			log.Println("태그/책 데이터 스캔 중 오류 발생:", err)
			continue
		}
		temp := model.BookInfo{
			BookID:    bookID,
			BookTitle: bookTitle,
			CoverURL:  coverURL,
			Grade:     grade,
			DelYN:     delYN,
		}
		all = append(all, temp)
	}
	log.Println("태그관련책표지:", all)
	return all
}
