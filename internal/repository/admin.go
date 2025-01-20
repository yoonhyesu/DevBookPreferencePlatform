package repository

import (
	"DBP/internal/model"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

// 개발자 조회
func (m *CommonRepo) GetDevList() []model.AddDevs {
	rows, err := m.mariaDB.Connection.Query(`
	SELECT ID, DEV_NAME, DEV_DETAIL_NAME, DEV_HISTORY, VIEW_YN, PROFILE_IMAGE_PATH, DEL_YN
	FROM dbp.dev_infos
	WHERE DEL_YN = 0
	ORDER BY ID
	`)
	if err != nil {
		log.Println(err)
		return []model.AddDevs{}
	}
	defer rows.Close()

	var devs []model.AddDevs
	for rows.Next() {
		var dev model.AddDevs
		var profileImagePath sql.NullString // NULL을 허용하는 타입으로 변경!!!

		err := rows.Scan(
			&dev.DevID,
			&dev.DevName,
			&dev.DevDetailName,
			&dev.DevHistory,
			&dev.ViewYN,
			&profileImagePath, // NullString으로 스캔!!!
			&dev.DelYN,
		)
		if err != nil {
			log.Println("개발자 데이터 스캔 오류:", err)
			continue
		}

		// NULL 체크 후 값 할당!!!
		if profileImagePath.Valid {
			dev.ProfileImagePath = profileImagePath.String
		} else {
			dev.ProfileImagePath = ""
		}

		devs = append(devs, dev)
	}
	return devs
}

// 특정 개발자 조회
func (m *CommonRepo) GetDevID(devID string) (model.AddDevs, error) {
	var dev model.AddDevs
	var profileImagePath sql.NullString

	err := m.mariaDB.Connection.QueryRow(`
    SELECT ID, DEV_NAME, DEV_DETAIL_NAME, DEV_HISTORY, VIEW_YN, PROFILE_IMAGE_PATH, DEL_YN
    FROM dbp.dev_infos
    WHERE DEL_YN = 0 AND ID = ?
    `, devID).Scan(
		&dev.DevID,
		&dev.DevName,
		&dev.DevDetailName,
		&dev.DevHistory,
		&dev.ViewYN,
		&profileImagePath,
		&dev.DelYN,
	)

	if err != nil {
		return model.AddDevs{}, err
	}

	if profileImagePath.Valid {
		dev.ProfileImagePath = profileImagePath.String
	} else {
		dev.ProfileImagePath = ""
	}

	return dev, nil
}

// 마지막 개발자 ID 조회
func (m *CommonRepo) getLastDevID() (string, error) {
	var lastID string
	err := m.mariaDB.Connection.QueryRow(`
	SELECT ID 
	FROM dbp.dev_infos 
	ORDER BY ID DESC 
	LIMIT 1
	`).Scan(&lastID)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "DEV_000", nil // 데이터가 없는 경우 초기값 반환
		}
		return "", err
	}
	return lastID, nil
}

// 새로운 개발자 ID 생성
func generateNextDevID(lastID string) string {
	// DEV_001 형식에서 숫자부분만 추출
	numStr := lastID[4:] // "001"
	num, _ := strconv.Atoi(numStr)
	nextNum := num + 1
	// 새로운 ID 생성 (3자리 숫자 형식 유지)
	return fmt.Sprintf("DEV_%03d", nextNum)
}

// 개발자 등록
func (m *CommonRepo) AddDev(dev model.AddDevs) error {
	// 마지막 ID 조회
	lastID, err := m.getLastDevID()
	if err != nil {
		return err
	}

	// 새로운 ID 생성
	newID := generateNextDevID(lastID)

	query := `
	INSERT INTO dbp.dev_infos 
	(ID, DEV_NAME, DEV_DETAIL_NAME, DEV_HISTORY, PROFILE_IMAGE_PATH, VIEW_YN) 
	VALUES (?, ?, ?, ?, ?, ?)
	`
	// 개발자 등록
	_, err = m.mariaDB.Connection.Exec(query,
		newID, dev.DevName, dev.DevDetailName, dev.DevHistory, dev.ProfileImagePath, dev.ViewYN)
	return err
}

// 개발자 수정
func (m *CommonRepo) UpdateDev(dev model.AddDevs) error {
	_, err := m.mariaDB.Connection.Exec(`
        UPDATE dbp.dev_infos 
        SET DEV_NAME = ?, DEV_DETAIL_NAME = ?, DEV_HISTORY = ?, PROFILE_IMAGE_PATH = ?, VIEW_YN = ?
        WHERE ID = ?`,
		dev.DevName, dev.DevDetailName, dev.DevHistory, dev.ProfileImagePath, dev.ViewYN, dev.DevID)
	if err != nil {
		return fmt.Errorf("개발자 업데이트 실패: %v", err)
	}
	return nil
}

// 개발자 삭제 (soft delete)
func (m *CommonRepo) DeleteDev(devID string) error {
	// dev_infos 테이블 업데이트
	_, err1 := m.mariaDB.Connection.Exec(`
	UPDATE dbp.dev_infos 
	SET DEL_YN = 1 
	WHERE ID = ?`, devID)

	// dev_recommends 테이블에서 삭제
	_, err2 := m.mariaDB.Connection.Exec(`
	DELETE FROM dbp.dev_recommends
	WHERE DEV_ID = ?`, devID)

	// 두 에러 중 하나라도 있다면 해당 에러 반환
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}

// 공지 조회
func (m *CommonRepo) GetNoticeList() []model.Notice {
	rows, err := m.mariaDB.Connection.Query(`
SELECT ID, NOTICE_TITLE, NOTICE_CONTENT, TOP_YN, CREATED_AT, UPDATED_AT, DEL_YN
FROM dbp.notices
WHERE DEL_YN = FALSE
ORDER BY ID
	`)
	if err != nil {
		log.Printf("공지사항 조회 오류: %v", err)
		return []model.Notice{}
	}
	defer rows.Close()

	var notices []model.Notice
	for rows.Next() {
		var notice model.Notice
		err := rows.Scan(
			&notice.ID,
			&notice.Title,
			&notice.Content,
			&notice.TopYn,
			&notice.CreatedDate,
			&notice.UpdatedDate,
			&notice.DelYn,
		)
		if err != nil {
			log.Printf("공지사항 데이터 스캔 오류: %v", err)
			continue
		}
		notices = append(notices, notice)
		log.Printf("스캔된 공지사항: %+v", notice)
	}
	return notices
}

// 공지 등록
func (m *CommonRepo) AddNotice(notice model.Notice) error {
	_, err := m.mariaDB.Connection.Exec(`
	INSERT INTO dbp.notices 
	(NOTICE_TITLE, NOTICE_CONTENT, TOP_YN, CREATED_AT, UPDATED_AT) 
	VALUES (?, ?, ?, NOW(), NOW())`,
		notice.Title, notice.Content, notice.TopYn)
	return err
}

// 공지 수정
func (m *CommonRepo) UpdateNotice(notice model.Notice) error {
	_, err := m.mariaDB.Connection.Exec(`
	UPDATE dbp.notices 
	SET NOTICE_TITLE = ?, NOTICE_CONTENT = ?, TOP_YN = ?, UPDATED_AT = NOW()
	WHERE ID = ?`,
		notice.Title, notice.Content, notice.TopYn, notice.ID)
	return err
}

// 공지 삭제 (soft delete)
func (m *CommonRepo) DeleteNotice(noticeID int) error {
	_, err := m.mariaDB.Connection.Exec(`
	UPDATE dbp.notices 
	SET DEL_YN = 1, UPDATED_AT = NOW()
	WHERE ID = ?`, noticeID)
	return err
}

// 책 조회
func (m *CommonRepo) GetBookList() []model.AddBooks {
	// 책 정보 조회
	rows, err := m.mariaDB.Connection.Query(`
	SELECT 
		BOOK_ID, BOOK_TITLE, PUBLISHER, AUTHOR, PUBLISH_DATE, 
		PRICE, CONTENTS_LIST, COVER_URL, SUMMARY, DESCRIPTION, 
		RECOMMENDATION, TAGS, GRADE, ISBN, 
		DEL_YN, ISBN_ADD, PAGE, EBOOK_YN
	FROM dbp.book_infos
	WHERE DEL_YN = 0
	ORDER BY BOOK_ID
	`)
	if err != nil {
		log.Printf("책 조회 오류: %v", err)
		return []model.AddBooks{}
	}
	defer rows.Close()

	var books []model.AddBooks
	for rows.Next() {
		var book model.AddBooks
		err := rows.Scan(
			&book.BookID, &book.BookTitle, &book.Publisher, &book.Author,
			&book.PublishDate, &book.Price, &book.ContentsList, &book.CoverURL,
			&book.Summary, &book.Description, &book.Recommendation,
			&book.Tags, &book.Grade, &book.ISBN,
			&book.DelYN, &book.IsbnAdd, &book.Page, &book.EbookYN,
		)
		if err != nil {
			log.Printf("책 데이터 스캔 오류: %v", err)
			continue
		}

		// 해당 책의 추천 개발자 정보 조회
		devRows, err := m.mariaDB.Connection.Query(`
		SELECT d.ID, d.DEV_NAME, dr.DEV_RECOMMEND_REASON
		FROM dbp.dev_recommends dr
		JOIN dbp.dev_infos d ON dr.DEV_ID = d.ID
		WHERE dr.BOOK_ID = ? AND d.DEL_YN = 0
		`, book.BookID)
		if err != nil {
			log.Printf("개발자 추천 정보 조회 오류: %v", err)
			continue
		}
		defer devRows.Close()

		var devRecommends []model.DevRecommend
		for devRows.Next() {
			var dev model.DevRecommend
			err := devRows.Scan(&dev.DevID, &dev.DevName, &dev.RecommendReason)
			if err != nil {
				log.Printf("개발자 추천 데이터 스캔 오류: %v", err)
				continue
			}
			devRecommends = append(devRecommends, dev)
		}
		book.DevRecommends = devRecommends

		books = append(books, book)
	}
	return books
}

// 책 등록 (트랜잭션 처리)
func (m *CommonRepo) AddBook(book model.AddBooks, devRecommends []model.DevRecommend) error {
	// 입력 데이터 로깅
	log.Printf("받은 책 정보: %+v", book)
	log.Printf("받은 개발자 추천 정보: %+v", devRecommends)

	tx, err := m.mariaDB.Connection.Begin()
	if err != nil {
		log.Printf("트랜잭션 시작 오류: %v", err)
		return err
	}
	defer tx.Rollback()

	// 1. 책 정보 등록
	result, err := tx.Exec(`
	INSERT INTO dbp.book_infos (
		BOOK_TITLE, PUBLISHER, AUTHOR, PUBLISH_DATE, PRICE,
		CONTENTS_LIST, COVER_URL, SUMMARY, DESCRIPTION,
		RECOMMENDATION, TAGS, GRADE, ISBN,
		ISBN_ADD, PAGE, EBOOK_YN
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		book.BookTitle, book.Publisher, book.Author, book.PublishDate,
		book.Price, book.ContentsList, book.CoverURL, book.Summary,
		book.Description, book.Recommendation,
		book.Tags, book.Grade, book.ISBN, book.IsbnAdd, book.Page, book.EbookYN,
	)

	// SQL 쿼리와 바인딩된 값들 로깅
	log.Printf("SQL 바인딩 값들: %q, %q, %q, %q, %d, %q, %q, %q, %q, %q, %q, %q, %q, %q, %d, %q",
		book.BookTitle, book.Publisher, book.Author, book.PublishDate,
		book.Price, book.ContentsList, book.CoverURL, book.Summary,
		book.Description, book.Recommendation,
		book.Tags, book.Grade, book.ISBN, book.IsbnAdd, book.Page, book.EbookYN)

	if err != nil {
		log.Printf("책 정보 등록 오류: %v", err)
		return err
	}

	bookID, err := result.LastInsertId()
	if err != nil {
		log.Printf("LastInsertId 오류: %v", err)
		return err
	}

	// 2. 개발자 추천 정보 등록
	for _, rec := range devRecommends {
		if rec.DevID == "" {
			continue // 빈 DevID는 건너뛰기
		}
		_, err = tx.Exec(`
		INSERT INTO dbp.dev_recommends (
			DEV_ID, BOOK_ID, DEV_RECOMMEND_REASON
		) VALUES (?, ?, ?)`,
			rec.DevID, bookID, rec.RecommendReason,
		)
		if err != nil {
			log.Printf("개발자 추천 정보 등록 오류: %v", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("트랜잭션 커밋 오류: %v", err)
		return err
	}

	return nil
}

// 책 수정 (트랜잭션 처리)
func (m *CommonRepo) UpdateBook(book model.AddBooks, devRecommends []model.DevRecommend) error {
	// 입력 데이터 로깅
	log.Printf("수정할 책 정보: %+v", book)
	log.Printf("수정할 개발자 추천 정보: %+v", devRecommends)

	tx, err := m.mariaDB.Connection.Begin()
	if err != nil {
		log.Printf("트랜잭션 시작 오류: %v", err)
		return err
	}
	defer tx.Rollback()

	// 1. 책 정보 수정
	_, err = tx.Exec(`
	UPDATE dbp.book_infos SET
		BOOK_TITLE = ?, PUBLISHER = ?, AUTHOR = ?, PUBLISH_DATE = ?, PRICE = ?,
		CONTENTS_LIST = ?, COVER_URL = ?, SUMMARY = ?, DESCRIPTION = ?,
		RECOMMENDATION = ?, TAGS = ?, GRADE = ?, ISBN = ?,
		ISBN_ADD = ?, PAGE = ?, EBOOK_YN = ?
	WHERE BOOK_ID = ? AND DEL_YN = 0`,
		book.BookTitle, book.Publisher, book.Author, book.PublishDate,
		book.Price, book.ContentsList, book.CoverURL, book.Summary,
		book.Description, book.Recommendation,
		book.Tags, book.Grade, book.ISBN, book.IsbnAdd, book.Page, book.EbookYN,
		book.BookID,
	)

	// SQL 쿼리와 바인딩된 값들 로깅
	log.Printf("SQL 바인딩 값들: %q, %q, %q, %q, %d, %q, %q, %q, %q, %q, %q, %q, %q, %q, %d, %q, %d",
		book.BookTitle, book.Publisher, book.Author, book.PublishDate,
		book.Price, book.ContentsList, book.CoverURL, book.Summary,
		book.Description, book.Recommendation,
		book.Tags, book.Grade, book.ISBN, book.IsbnAdd, book.Page, book.EbookYN,
		book.BookID)

	if err != nil {
		log.Printf("책 정보 수정 오류: %v", err)
		return err
	}

	// 2. 기존 개발자 추천 정보 삭제
	_, err = tx.Exec(`DELETE FROM dbp.dev_recommends WHERE BOOK_ID = ?`, book.BookID)
	if err != nil {
		log.Printf("기존 추천 정보 삭제 오류: %v", err)
		return err
	}

	// 3. 새로운 개발자 추천 정보 등록
	for _, rec := range devRecommends {
		if rec.DevID == "" {
			continue // 빈 DevID는 건너뛰기
		}
		_, err = tx.Exec(`
		INSERT INTO dbp.dev_recommends (
			DEV_ID, BOOK_ID, DEV_RECOMMEND_REASON
		) VALUES (?, ?, ?)`,
			rec.DevID, book.BookID, rec.RecommendReason,
		)
		if err != nil {
			log.Printf("개발자 추천 정보 등록 오류: %v", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("트랜잭션 커밋 오류: %v", err)
		return err
	}

	return nil
}

// 책 삭제 (soft delete)
func (m *CommonRepo) DeleteBook(bookID int) error {
	_, err := m.mariaDB.Connection.Exec(`
	UPDATE dbp.book_infos 
	SET DEL_YN = 1 
	WHERE BOOK_ID = ?`, bookID)
	return err
}

// 태그 목록 조회
func (m *CommonRepo) GetAdminTagList() []model.TagList {
	rows, err := m.mariaDB.Connection.Query(`
	SELECT TAG_ID, TAG_NAME
	FROM dbp.tags
	WHERE DEL_YN = 0
	ORDER BY TAG_ID ASC
	`)
	if err != nil {
		log.Printf("태그 조회 오류: %v", err)
		return []model.TagList{}
	}
	defer rows.Close()

	var tags []model.TagList
	for rows.Next() {
		var tag model.TagList
		err := rows.Scan(
			&tag.TagID,
			&tag.TagName,
		)
		if err != nil {
			log.Printf("태그 데이터 스캔 오류: %v", err)
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}
