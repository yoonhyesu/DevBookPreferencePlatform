package repository

import (
	"DBP/internal/model"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// 책 상세페이지
func (m *CommonRepo) GetBookDetail(bookID int) []model.BookInfo {
	row := m.mariaDB.Connection.QueryRow(`
		SELECT 
    bi.*,
    GROUP_CONCAT(t.tag_name SEPARATOR ', ') as tag_names
FROM 
    dbp.book_infos bi
    LEFT JOIN dbp.tags t ON t.tag_id IN (
        SELECT SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING(TAGS, 2, LENGTH(TAGS)-2), '&', n.n), '&', -1) as tag_id
        FROM 
            (SELECT 1 n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5) n
        WHERE SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING(bi.TAGS, 2, LENGTH(bi.TAGS)-2), '&', n.n), '&', -1) != ''
    )
WHERE bi.BOOK_ID = ?
GROUP BY 
    bi.BOOK_ID`, bookID)

	all := []model.BookInfo{}
	var c1, c6, c17 int
	var c2, c3, c4, c5, c7, c8, c9, c10, c11, c12, c13, c14, c16, c18, c19 string
	var c15 bool

	err := row.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8, &c9, &c10, &c11, &c12, &c13, &c14, &c15, &c16, &c17, &c18, &c19)
	if err != nil {
		log.Println(err)
		return []model.BookInfo{}
	}
	temp := model.BookInfo{
		BookID:         c1,
		BookTitle:      c2,
		Publisher:      c3,
		Author:         c4,
		PublishDate:    c5,
		Price:          c6,
		ContentsList:   c7,
		CoverURL:       c8,
		Summary:        c9,
		Description:    c10,
		Recommendation: c11,
		Tags:           c12,
		Grade:          c13,
		ISBN:           c14,
		DelYN:          c15,
		IsbnAdd:        c16,
		Page:           c17,
		EbookYN:        c18,
		TagNames:       c19,
	}
	all = append(all, temp)
	log.Println("책 상세정보:", all)
	return all
}

// 목차 가져오기
func (m *CommonRepo) GetBookContentsList(bookID int) string {
	row := m.mariaDB.Connection.QueryRow(`
	SELECT CONTENTS_LIST FROM dbp.book_infos WHERE BOOK_ID = ?`, bookID)
	var contentsList string
	err := row.Scan(&contentsList)
	if err != nil {
		log.Println(err)
		return ""
	}
	return contentsList
}

// 사용자가 책을 좋아요했는지 확인
func (m *CommonRepo) CheckBookLiked(userID, bookID string) bool {
	var likeList sql.NullString // NULL 값을 처리하기 위해 sql.NullString 사용
	err := m.mariaDB.Connection.QueryRow("SELECT LIKE_BOOK_LIST FROM users WHERE USER_ID = ?", userID).Scan(&likeList)
	if err != nil {
		log.Printf("좋아요 목록 조회 실패: %v", err)
		return false
	}

	// NULL 값이거나 빈 문자열이면 좋아요하지 않은 것
	if !likeList.Valid || likeList.String == "" {
		return false
	}

	// &23&7& 형식에서 &7& 부분이 있는지 확인
	bookIDWithDelimiters := "&" + bookID + "&"
	return strings.Contains(likeList.String, bookIDWithDelimiters)
}

// 책 좋아요 상태 업데이트
func (m *CommonRepo) UpdateBookLike(userID, bookID string, liked bool) error {
	log.Printf("UpdateBookLike 시작: userID=%s, bookID=%s, liked=%v", userID, bookID, liked)

	// 현재 좋아요 목록 가져오기
	var likeList sql.NullString
	err := m.mariaDB.Connection.QueryRow("SELECT LIKE_BOOK_LIST FROM users WHERE USER_ID = ?", userID).Scan(&likeList)
	if err != nil {
		log.Printf("좋아요 목록 조회 실패: %v", err)
		return fmt.Errorf("좋아요 목록 조회 실패: %w", err)
	}

	// 좋아요 목록 문자열 처리
	var likeListStr string
	if likeList.Valid {
		likeListStr = likeList.String
	} else {
		likeListStr = ""
	}

	log.Printf("현재 좋아요 목록: %s", likeListStr)

	// 좋아요 추가 또는 제거
	bookIDPattern := fmt.Sprintf("&%s&", bookID)

	if liked {
		// 이미 좋아요한 경우 중복 추가 방지
		if !strings.Contains(likeListStr, bookIDPattern) {
			if likeListStr == "" {
				// 첫 번째 좋아요인 경우
				likeListStr = "&" + bookID + "&"
			} else {
				// 기존 목록에 추가 (마지막 & 다음에 번호 추가)
				likeListStr = strings.TrimSuffix(likeListStr, "&") + "&" + bookID + "&"
			}
		}
	} else {
		// 좋아요 취소: 해당 ID만 제거
		likeListStr = strings.Replace(likeListStr, bookIDPattern, "&", 1)

		// 연속된 && 처리
		for strings.Contains(likeListStr, "&&") {
			likeListStr = strings.Replace(likeListStr, "&&", "&", -1)
		}

		// 빈 목록이 되면 &만 남는 경우 처리
		if likeListStr == "&" {
			likeListStr = ""
		}
	}

	log.Printf("업데이트할 좋아요 목록: %s", likeListStr)

	// 데이터베이스 업데이트
	_, err = m.mariaDB.Connection.Exec("UPDATE users SET LIKE_BOOK_LIST = ? WHERE USER_ID = ?", likeListStr, userID)
	if err != nil {
		log.Printf("좋아요 목록 업데이트 실패: %v", err)
		return fmt.Errorf("좋아요 목록 업데이트 실패: %w", err)
	}

	log.Printf("좋아요 목록 업데이트 성공")
	return nil
}

// 사용자가 좋아요한 책 목록 가져오기
func (m *CommonRepo) GetLikedBooks(userID string) []model.BookInfo {
	var likeList sql.NullString
	err := m.mariaDB.Connection.QueryRow("SELECT LIKE_BOOK_LIST FROM users WHERE USER_ID = ?", userID).Scan(&likeList)
	if err != nil {
		log.Println("좋아요 목록 조회 실패:", err)
		return []model.BookInfo{}
	}

	// 좋아요한 책이 없으면 빈 배열 반환
	if !likeList.Valid || likeList.String == "" {
		return []model.BookInfo{}
	}

	// &23&7& 형식에서 ID만 추출
	likeListStr := likeList.String
	// 첫 번째와 마지막 & 제거
	if len(likeListStr) > 0 && likeListStr[0] == '&' {
		likeListStr = likeListStr[1:]
	}
	if len(likeListStr) > 0 && likeListStr[len(likeListStr)-1] == '&' {
		likeListStr = likeListStr[:len(likeListStr)-1]
	}

	// 남은 &로 분리
	bookIDs := strings.Split(likeListStr, "&")

	// 빈 ID 제거
	var validBookIDs []string
	for _, id := range bookIDs {
		if id != "" {
			validBookIDs = append(validBookIDs, id)
		}
	}

	if len(validBookIDs) == 0 {
		return []model.BookInfo{}
	}

	// 책 정보 조회를 위한 IN 쿼리 생성
	placeholders := make([]string, len(validBookIDs))
	args := make([]interface{}, len(validBookIDs))
	for i, id := range validBookIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	// 테이블 이름을 book_infos로 수정
	query := fmt.Sprintf(`
		SELECT 
			bi.BOOK_ID, bi.BOOK_TITLE, bi.AUTHOR, bi.PUBLISHER, 
			bi.PUBLISH_DATE, bi.PRICE, bi.PAGE, bi.GRADE, 
			bi.EBOOK_YN, bi.COVER_URL, bi.ISBN, bi.DESCRIPTION, 
			bi.SUMMARY, bi.RECOMMENDATION, GROUP_CONCAT(t.TAG_NAME) as TAG_NAMES
		FROM 
			book_infos bi
		LEFT JOIN 
			tags t ON t.TAG_ID IN (
				SELECT SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING(bi.TAGS, 2, LENGTH(bi.TAGS)-2), '&', n.n), '&', -1) as tag_id
				FROM 
					(SELECT 1 n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5) n
				WHERE SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING(bi.TAGS, 2, LENGTH(bi.TAGS)-2), '&', n.n), '&', -1) != ''
			)
		WHERE 
			bi.BOOK_ID IN (%s)
		GROUP BY 
			bi.BOOK_ID
	`, strings.Join(placeholders, ","))

	rows, err := m.mariaDB.Connection.Query(query, args...)
	if err != nil {
		log.Println("좋아요한 책 목록 조회 실패:", err)
		return []model.BookInfo{}
	}
	defer rows.Close()

	var books []model.BookInfo
	for rows.Next() {
		var book model.BookInfo
		err := rows.Scan(
			&book.BookID, &book.BookTitle, &book.Author, &book.Publisher,
			&book.PublishDate, &book.Price, &book.Page, &book.Grade,
			&book.EbookYN, &book.CoverURL, &book.ISBN, &book.Description,
			&book.Summary, &book.Recommendation, &book.TagNames,
		)
		if err != nil {
			log.Println("책 정보 스캔 실패:", err)
			continue
		}
		books = append(books, book)
	}

	return books
}
