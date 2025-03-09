package repository

import (
	"DBP/internal/model"
	"os"
	"strings"
	"time"

	"log"
)

// 공지사항 조회
func (m *CommonRepo) GetNotice() []model.Notice {
	rows, err := m.mariaDB.Connection.Query(
		`SELECT ID, NOTICE_TITLE, NOTICE_CONTENT, TOP_YN, CREATED_AT, UPDATED_AT, DEL_YN
FROM dbp.notices
WHERE DEL_YN = FALSE
ORDER BY TOP_YN DESC, CREATED_AT DESC`,
	)
	if err != nil {
		log.Println(err)
		return []model.Notice{}
	}
	all := []model.Notice{}
	var c1 int
	var c2, c3 string
	var c4, c5 time.Time
	var c6, c7 bool

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7)
		temp := model.Notice{ID: c1, Title: c2, Content: c3, CreatedDate: c4, UpdatedDate: c5, DelYn: c6, TopYn: c7}
		all = append(all, temp)
	}
	log.Println("공지사항:", all)
	return all
}

// 메인 노출 개발자 조회
// 개발자가 추천한 책 목록 조회
func (m *CommonRepo) GetHomeDevBook() model.HomeDevBook {
	var devBook model.HomeDevBook

	// 먼저 표시할 개발자 한 명을 선택합니다 (랜덤 또는 특정 조건으로)
	var devID string
	err := m.mariaDB.Connection.QueryRow(`
		SELECT ID 
		FROM dbp.dev_infos 
		WHERE VIEW_YN = 1 AND DEL_YN = 0
		ORDER BY RAND() LIMIT 1`).Scan(&devID)

	if err != nil {
		log.Println("개발자 선택 오류:", err)
		return model.HomeDevBook{}
	}

	// 선택된 개발자 정보와 그 개발자가 추천한 책만 조회
	rows, err := m.mariaDB.Connection.Query(`
    SELECT 
        B.BOOK_ID,
        B.BOOK_TITLE,
        B.COVER_URL,
        C.ID,
        C.DEV_NAME,
        C.DEV_HISTORY,
        C.DEV_DETAIL_NAME,
		C.PROFILE_IMAGE_PATH
    FROM dbp.dev_recommends A
    INNER JOIN dbp.BOOK_INFOS B ON A.BOOK_ID = B.BOOK_ID 
    INNER JOIN dbp.dev_infos C ON A.DEV_ID = C.ID
    WHERE A.DEL_YN = 0 
        AND C.VIEW_YN = 1
        AND C.DEL_YN = 0
        AND C.ID = ?`, devID)

	if err != nil {
		log.Println(err)
		return model.HomeDevBook{}
	}
	defer rows.Close()

	first := true
	for rows.Next() {
		var bookID int
		var bookTitle, coverURL, devID, devName, devHistory, devDetailName, profileImagePath string
		err = rows.Scan(&bookID, &bookTitle, &coverURL, &devID, &devName, &devHistory, &devDetailName, &profileImagePath)
		if err != nil {
			log.Println("책 정보 스캔 중 오류:", err)
			continue
		}

		// 첫 번째 행에서만 개발자 정보 설정
		if first {
			devBook.DevID = devID
			devBook.DevName = devName
			devBook.DevHistory = devHistory
			devBook.DevDetailName = devDetailName
			devBook.ProfileImagePath = os.Getenv("SELECT_DEV_PATH") + "/" + profileImagePath
			devBook.Books = []model.Book{}
			first = false
		}

		// 책 정보 추가
		book := model.Book{
			BookID:    bookID,
			BookTitle: bookTitle,
			CoverURL:  coverURL,
		}
		devBook.Books = append(devBook.Books, book)
	}

	// 이미지 경로 처리
	if devBook.ProfileImagePath != "" {
		// 백슬래시를 슬래시로 변환
		devBook.ProfileImagePath = strings.ReplaceAll(devBook.ProfileImagePath, "\\", "/")
	}

	return devBook
}
