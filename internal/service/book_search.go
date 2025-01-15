package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// 도서 검색 서비스

type BookSearchService struct {
	certKey string
	baseURL string
}

// 외부 API 응답 구조체
type BookSearchResponse struct {
	Result     string    `json:"RESULT"`
	ErrMsg     string    `json:"ERR_MESSAGE,omitempty"`
	PageNo     string    `json:"PAGE_NO"`
	PageSize   string    `json:"PAGE_SIZE"`
	TotalCount string    `json:"TOTAL_COUNT"`
	Documents  []BookDoc `json:"docs"`
}

// 도서 정보 구조체
type BookDoc struct {
	Title          string `json:"TITLE"`             // 도서명
	Author         string `json:"AUTHOR"`            // 저자
	Publisher      string `json:"PUBLISHER"`         // 출판사
	PublishPreDate string `json:"PUBLISH_PREDATE"`   // 출판예정일
	EaISBN         string `json:"EA_ISBN"`           // ISBN
	EaAddCode      string `json:"EA_ADD_CODE"`       // 부가기호
	PrePrice       string `json:"PRE_PRICE"`         // 예정가격
	Page           string `json:"PAGE"`              // 페이지 수
	TitleURL       string `json:"TITLE_URL"`         // 표지 이미지 URL
	Description    string `json:"BOOK_INTRODUCTION"` // 도서 소개
	Ebookyn        string `json:"EBOOK_YN"`          // 전자책 여부
	Summary        string `json:"BOOK_SUMMARY"`      // 도서 요약
	Cnt            string `json:"BOOK_TB_CNT"`       // 목차
}

func NewBookSearchService() *BookSearchService {
	return &BookSearchService{
		certKey: os.Getenv("CERT_KEY"),
		baseURL: os.Getenv("BASE_URL"),
	}
}

// 도서 검색
func (s *BookSearchService) SearchBooks(title string) (*BookSearchResponse, error) {
	// 검색어 전처리
	searchTitle := strings.TrimSpace(title)                // 앞뒤 공백 제거
	searchTitle = strings.ToLower(searchTitle)             // 소문자로 변환
	searchTitle = strings.ReplaceAll(searchTitle, " ", "") // 모든 공백 제거

	// URL 인코딩
	encodedTitle := url.QueryEscape(searchTitle)

	// API URL 생성
	apiURL := fmt.Sprintf("%s?cert_key=%s&result_style=json&page_no=1&page_size=50&title=%s",
		s.baseURL,
		s.certKey,
		encodedTitle)
	log.Printf("검색 API URL: %s", apiURL)

	// HTTP GET 요청
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 응답 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %v", err)
	}

	// JSON 파싱
	var result BookSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON 파싱 실패: %v", err)
	}

	// 검색 결과 후처리 (대소문자 무시하고 부분 일치 검색)
	var filteredDocs []BookDoc
	searchTitle = strings.ToLower(title) // 원본 검색어를 소문자로
	for _, doc := range result.Documents {
		bookTitle := strings.ToLower(doc.Title)
		// 공백을 제거한 버전과 원본 모두로 검색
		bookTitleNoSpace := strings.ReplaceAll(bookTitle, " ", "")
		searchTitleNoSpace := strings.ReplaceAll(searchTitle, " ", "")

		if strings.Contains(bookTitle, searchTitle) ||
			strings.Contains(bookTitleNoSpace, searchTitleNoSpace) {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	result.Documents = filteredDocs

	log.Printf("검색 결과: %d건", len(result.Documents))
	return &result, nil
}
