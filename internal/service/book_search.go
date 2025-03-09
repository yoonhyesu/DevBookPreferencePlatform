package service // 도서 검색 서비스

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
	// 환경 변수 확인 및 로깅
	certKey := os.Getenv("CERT_KEY")
	baseURL := os.Getenv("BASE_URL")

	// 환경 변수 확인 로그
	log.Printf("BookSearchService 초기화: CERT_KEY=%s, BASE_URL=%s",
		maskSecret(certKey), baseURL)

	// 환경 변수가 비어있는지 확인
	if certKey == "" {
		log.Println("경고: CERT_KEY 환경 변수가 설정되지 않았습니다")
	}

	if baseURL == "" {
		log.Println("경고: BASE_URL 환경 변수가 설정되지 않았습니다")
		// 기본값 설정 (선택사항)
		baseURL = "https://api.default-book-service.com"
	}

	return &BookSearchService{
		certKey: certKey,
		baseURL: baseURL,
	}
}

// 비밀 값 마스킹 헬퍼 함수
func maskSecret(secret string) string {
	if len(secret) <= 4 {
		return "****"
	}
	return secret[:2] + "****" + secret[len(secret)-2:]
}

// 도서 검색
func (s *BookSearchService) SearchBooks(title string) (*BookSearchResponse, error) {
	// 서비스 상태 확인
	if s.baseURL == "" {
		return nil, fmt.Errorf("BASE_URL이 설정되지 않았습니다")
	}

	if s.certKey == "" {
		return nil, fmt.Errorf("CERT_KEY가 설정되지 않았습니다")
	}

	// 검색어 전처리
	searchTitle := strings.TrimSpace(title)    // 앞뒤 공백 제거
	searchTitle = strings.ToLower(searchTitle) // 소문자로 변환

	// URL 인코딩
	encodedTitle := url.QueryEscape(searchTitle)

	// API URL 생성
	apiURL := fmt.Sprintf("%s?cert_key=%s&result_style=json&page_no=1&page_size=1000&title=%s&sort=PUBLISH_PREDATE&order_by=DESC",
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

	// 검색어에서 연속된 공백을 하나로 정규화
	searchTitle = strings.Join(strings.Fields(searchTitle), " ")

	for _, doc := range result.Documents {
		bookTitle := strings.ToLower(doc.Title)
		// 도서 제목에서도 연속된 공백을 하나로 정규화
		bookTitle = strings.Join(strings.Fields(bookTitle), " ")

		// 정규화된 문자열로 검색
		if strings.Contains(bookTitle, searchTitle) {
			filteredDocs = append(filteredDocs, doc)
			continue
		}
		// 공백을 완전히 제거하고 검색
		bookTitleNoSpace := strings.ReplaceAll(bookTitle, " ", "")
		searchTitleNoSpace := strings.ReplaceAll(searchTitle, " ", "")
		if strings.Contains(bookTitleNoSpace, searchTitleNoSpace) {
			filteredDocs = append(filteredDocs, doc)
		}
	}

	result.Documents = filteredDocs

	log.Printf("검색 결과: %d건", len(result.Documents))
	return &result, nil
}
