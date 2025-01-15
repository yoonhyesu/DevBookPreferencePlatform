package repository

import (
	"DBP/internal/model"
	"log"
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
