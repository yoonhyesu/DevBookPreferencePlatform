package repository

import (
	"DBP/internal/model"
	"os"

	"log"
	"strings"
)

// 태그 조회 (가나다순 정렬)
func (m *CommonRepo) GetTagList() []model.TagList {
	rows, err := m.mariaDB.Connection.Query("SELECT TAG_ID, TAG_NAME, DEL_YN FROM TAGS WHERE DEL_YN = FALSE ORDER BY TAG_NAME ASC")
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

// 개발자 리스트 조회 - DevInfo 타입으로 변경
func (m *CommonRepo) GetRecommendDevList() []model.DevInfo {
	rows, err := m.mariaDB.Connection.Query(`
	SELECT 
		ID, 
		DEV_NAME, 
		DEV_DETAIL_NAME, 
		DEV_HISTORY,
		PROFILE_IMAGE_PATH,
		VIEW_YN, 
		DEL_YN
	FROM dbp.dev_infos
	WHERE VIEW_YN = TRUE AND DEL_YN = FALSE`)

	if err != nil {
		log.Println(err)
		return []model.DevInfo{}
	}

	var devList []model.DevInfo
	for rows.Next() {
		var dev model.DevInfo
		err = rows.Scan(
			&dev.DevID,
			&dev.DevName,
			&dev.DevDetailName,
			&dev.DevHistory,
			&dev.ProfileImagePath,
			&dev.ViewYN,
			&dev.DelYN,
		)
		if err != nil {
			log.Println("개발자 정보 스캔 중 오류:", err)
			continue
		}

		// 이미지 경로 처리
		if dev.ProfileImagePath != "" {
			dev.ProfileImagePath = strings.ReplaceAll(dev.ProfileImagePath, "\\", "/")
			log.Printf("개발자 %s의 이미지 경로 변환: %s", dev.DevName, dev.ProfileImagePath)
		}

		dev.ProfileImagePath = os.Getenv("SELECT_DEV_PATH") + "/" + dev.ProfileImagePath

		devList = append(devList, dev)
	}
	return devList
}

// 개발자 상세 정보 조회
func (m *CommonRepo) GetDevInfo(devID string) []model.DevDetail {
	row := m.mariaDB.Connection.QueryRow(`
	SELECT 
		ID, 
		DEV_NAME, 
		DEV_DETAIL_NAME, 
		DEV_HISTORY,
		PROFILE_IMAGE_PATH,
		VIEW_YN, 
		DEL_YN
	FROM dbp.dev_infos
	WHERE ID = ? AND DEL_YN = FALSE`, devID)

	all := []model.DevDetail{}
	var c1, c2, c3, c4, c5 string
	var c6, c7 bool

	err := row.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7)
	if err != nil {
		log.Println("개발자 정보 조회 중 오류:", err)
		return []model.DevDetail{}
	}

	// 이미지 경로 처리
	if c5 != "" {
		c5 = strings.ReplaceAll(c5, "\\", "/")
		log.Printf("개발자 %s의 이미지 경로 변환: %s", c2, c5)
	}
	c5 = os.Getenv("SELECT_DEV_PATH") + "/" + c5

	temp := model.DevDetail{
		DevID:            c1,
		DevName:          c2,
		DevDetailName:    c3,
		DevHistory:       c4,
		ProfileImagePath: c5,
		ViewYN:           c6,
		DelYN:            c7,
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
		B.COVER_URL,
		C.PROFILE_IMAGE_PATH
	FROM dbp.dev_recommends A
	INNER JOIN dbp.book_infos B ON A.BOOK_ID = B.BOOK_ID
	INNER JOIN dbp.dev_infos C ON A.DEV_ID = C.ID
	WHERE A.DEL_YN = 0 AND A.DEV_ID = ?`, devID)

	if err != nil {
		log.Println(err)
		return []model.DevBookLikeReason{}
	}
	defer rows.Close()

	var books []model.DevBookLikeReason
	for rows.Next() {
		var book model.DevBookLikeReason
		err = rows.Scan(
			&book.DevReasonID,
			&book.DevReason,
			&book.BookID,
			&book.BookTitle,
			&book.CoverURL,
			&book.ProfileImagePath,
		)
		if err != nil {
			log.Println("책 정보 스캔 중 오류:", err)
			continue
		}
		book.ProfileImagePath = os.Getenv("SELECT_PROFILE_PATH") + "/" + book.ProfileImagePath
		books = append(books, book)
	}
	return books
}

// 단계별 책 정보(표지) 가져옴
func (m *CommonRepo) GetStepWithBook(step string) []model.BookInfo {
	rows, err := m.mariaDB.Connection.Query(`
SELECT BOOK_ID, BOOK_TITLE, COVER_URL, GRADE, DEL_YN
FROM dbp.book_infos
WHERE GRADE = ? AND DEL_YN = FALSE
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
