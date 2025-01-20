package service

import (
	pkg "DBP/pkg/database"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"gopkg.in/yaml.v2"
)

// 사용할 llm 모델 명
const generativeModelName = "gemini-1.5-flash"
const embeddingModelName = "text-embedding-004"

type RagServer struct {
	ctx          context.Context
	geminiClient *googleai.GoogleAI
	key          *RagBookKey
	redisClient  *redis.Client
}

func NewRagServer() *RagServer {
	return &RagServer{}
}

// 환경변수
type RagBookKey struct {
	certKey string
	baseURL string
}

// 개별 도서 데이터
type RagBook struct {
	Title          string `json:"TITLE"`           // 도서명
	Author         string `json:"AUTHOR"`          // 저자
	Publisher      string `json:"PUBLISHER"`       // 출판사
	PublishPreDate string `json:"PUBLISH_PREDATE"` // 출판예정일
	EaISBN         string `json:"EA_ISBN"`         // ISBN
	EaAddCode      string `json:"EA_ADD_CODE"`     // 부가기호
	TitleURL       string `json:"TITLE_URL"`       // 표지 이미지 URL
	Cnt            string `json:"BOOK_TB_CNT_URL"` // 목차
	RecommendUser  string `json:"RECOMMEND_USER"`  // 추천 대상
	Keywords       string `json:"KEYWORDS"`        // 키워드
}

// 응답구조
type BookSearchResponse struct {
	Result     string    `json:"RESULT"`
	ErrMsg     string    `json:"ERR_MESSAGE,omitempty"`
	PageNo     string    `json:"PAGE_NO"`
	PageSize   string    `json:"PAGE_SIZE"`
	TotalCount string    `json:"TOTAL_COUNT"`
	Documents  []RagBook `json:"docs"`
}

// 레디스 키
const REDIS_KEY = "MAIN_NEW_2BOOKS"

// 제목 키워드
type Keywords struct {
	ITKeywords []string `yaml:"it_keywords"`
}

// 키워드 파일 경로 찾기
var configPath string

// 출판사 목록 상수 추가
var PUBLISHERS = []string{
	"씨엔아이파트너스",
	"한빛미디어",
	"위키북스",
	"이지스퍼블리싱",
	"길벗",
	"에이콘출판사",
	"제이펍",
	"프리렉",
}

func init() {
	// cmd 디렉토리 기준으로 상대 경로 설정
	defaultPath := filepath.Join("..", "config")
	flag.StringVar(&configPath, "config", defaultPath, "설정 파일 경로")
	flag.Parse()
}

// IT 키워드 경로 탐색 및 로드
func loadKeywords() string {
	filePath := filepath.Join(configPath, "keyword_it.yaml")

	log.Printf("키워드 파일 경로: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("키워드 파일 읽기 실패: %v", err)
		return ""
	}
	var keywords Keywords
	if err := yaml.Unmarshal(data, &keywords); err != nil {
		log.Printf("YAML 파싱 실패: %v", err)
		return ""
	}
	// 키워드들을 || 로 연결
	return strings.Join(keywords.ITKeywords, "||")
}

func GetRagBookService() *RagServer {
	ctx := context.Background()
	geminiKey := os.Getenv("GEMINI_API_KEY")
	geminiClient, err := googleai.New(ctx,
		googleai.WithAPIKey(geminiKey),
		googleai.WithDefaultEmbeddingModel(embeddingModelName))
	if err != nil {
		log.Fatal(err)
	}
	BookApiKey := &RagBookKey{
		certKey: os.Getenv("CERT_KEY"),
		baseURL: os.Getenv("BASE_URL"),
	}

	// Redis 클라이언트 가져오기
	redis := pkg.GetRedisInstance()

	return &RagServer{
		ctx:          ctx,
		key:          BookApiKey,
		geminiClient: geminiClient,
		redisClient:  redis.Client,
	}
}

// Gemini 도서 검증 플러그인
func (rs *RagServer) GeminiforBook(book RagBook) (bool, string, string) {
	// 키워드 체크 제거 (fetchLatestBooks에서 이미 체크함)
	prompt := fmt.Sprintf(`당신은 IT 도서 큐레이터입니다. 다음 도서가 컴퓨터/IT 관련 도서인지, 또는 개발자에게 도움될 책인지 판단해주세요.
	아래 형식으로만 답변해주세요:
	1. "예" 또는 "아니오"로 먼저 답변
	2. 만약 "예"라면 다음 줄에 "추천대상: [입문개발자/초급개발자/중급개발자/고급개발자]" 중 하나만 작성
	3. 만약 "예'라면 다음 줄에 "관련 키워드: [키워드1/키워드2]" 2개를 ,로 묶어서 작성
	예시, 추천대상 : 입문개발자 / 관련 키워드 : 프로그래밍, 파이썬
	
	제목:%s
	추천대상: %s
	키워드 : %s
	저자:%s 
	출판사:%s 
	출판예정일:%s 
	표지URL:%s`,
		book.Title, book.RecommendUser, book.Keywords, book.Author, book.Publisher, book.PublishPreDate, book.TitleURL)

	response, err := llms.GenerateFromSinglePrompt(rs.ctx, rs.geminiClient, prompt, llms.WithModel(generativeModelName))
	if err != nil {
		log.Printf("Gemini 호출 실패: %v", err)
		return false, "", ""
	}
	log.Printf("Gemini 응답: %s", response)

	// "예"가 포함되어 있는지 확인
	if !strings.Contains(strings.ToLower(response), "예") {
		return false, "", ""
	}

	// 키워드 추출
	keywordIndex := strings.Index(response, "관련 키워드:")
	var keywords string
	if keywordIndex != -1 {
		keywordPart := response[keywordIndex:]
		endIndex := strings.Index(keywordPart, "\n")
		if endIndex != -1 {
			keywords = strings.TrimSpace(keywordPart[len("관련 키워드:"):endIndex])
		} else {
			keywords = strings.TrimSpace(keywordPart[len("관련 키워드:"):])
		}
	}

	// 추천 대상 찾기
	recommendTargets := []string{
		"입문 개발자", "입문개발자",
		"초급 개발자", "초급개발자",
		"중급 개발자", "중급개발자",
		"고급 개발자", "고급개발자",
	}

	for _, target := range recommendTargets {
		if strings.Contains(response, target) {
			return true, strings.ReplaceAll(target, " ", ""), keywords
		}
	}
	return false, "", ""
}

// 캐시된 데이터 조회
func (rs *RagServer) GetCachedBooks() (*BookSearchResponse, error) {
	val, err := rs.redisClient.Get(rs.ctx, REDIS_KEY).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var response BookSearchResponse
	if err := json.Unmarshal([]byte(val), &response); err != nil {
		return nil, err
	}
	log.Print(response.Documents, response.Result)
	return &response, nil
}

// 데이터 캐시 레디스에 저장
func (rs *RagServer) casheTodayBooks(response *BookSearchResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	// 24T 유효시간 설정
	return rs.redisClient.Set(rs.ctx, REDIS_KEY, data, 24*time.Hour).Err()
}

// 키워드 포함 여부 확인하는 함수 추가
func containsITKeyword(title string, keywords []string) bool {
	title = strings.ToLower(title)
	for _, keyword := range keywords {
		if strings.Contains(title, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// 캐시 책정보 관리
func (rs *RagServer) GetCasheBook() *BookSearchResponse {
	// 캐시된 데이터 확인
	if cashedBooks, err := rs.GetCachedBooks(); err == nil && cashedBooks != nil {
		// 캐시는 있지만 도서가 없는 경우도 다시 fetch 하도록 수정!!!
		if len(cashedBooks.Documents) == 0 {
			log.Printf("캐시에 도서 정보가 없어 새로운 도서 정보를 가져옵니다")
			return rs.fetchAndCacheBooks()
		}
		return cashedBooks
	}

	return rs.fetchAndCacheBooks()
}

// 새로운 함수 추가: fetch와 cache를 함께 처리
func (rs *RagServer) fetchAndCacheBooks() *BookSearchResponse {
	response := rs.fetchLatestBooks()

	// 도서가 있는 경우에만 캐시 저장
	if len(response.Documents) > 0 {
		if err := rs.casheTodayBooks(response); err != nil {
			log.Printf("캐시 저장 실패: %v", err)
		}
	}

	return response
}

// 실제 API 호출
func (rs *RagServer) fetchLatestBooks() *BookSearchResponse {
	startDate := time.Now().AddDate(0, -1, 0).Format("20060102")
	endDate := time.Now().Format("20060102")

	var allBooks []RagBook

	// 각 출판사별로 API 호출
	for _, publisher := range PUBLISHERS {
		encodedPublisher := url.QueryEscape(publisher)
		apiURL := fmt.Sprintf("%s?cert_key=%s&result_style=json&page_no=1&page_size=1000&start_publish_date=%s&end_publish_date=%s&sort=INPUT_DATE&order_by=DESC&publisher=%s",
			rs.key.baseURL,
			rs.key.certKey,
			startDate,
			endDate,
			encodedPublisher,
		)

		resp, err := http.Get(apiURL)
		if err != nil {
			log.Printf("출판사 %s API 호출 실패: %v", publisher, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("출판사 %s 응답 읽기 실패: %v", publisher, err)
			continue
		}

		var searchResp BookSearchResponse
		if err := json.Unmarshal(body, &searchResp); err != nil {
			log.Printf("출판사 %s JSON 파싱 실패: %v", publisher, err)
			continue
		}

		// 검색된 도서들을 allBooks에 추가
		allBooks = append(allBooks, searchResp.Documents...)
	}

	// 키워드 로드는 한 번만 수행
	keywordsStr := loadKeywords()
	keywords := strings.Split(keywordsStr, "||")

	ResultBooks := make([]RagBook, 0)
	publisherCount := make(map[string]int) // 출판사별 선택된 책 수를 추적!!!

	for _, book := range allBooks {
		if book.TitleURL == "" || book.PublishPreDate > time.Now().Format("20060102") {
			continue
		}

		// 이미 해당 출판사에서 책을 선택했다면 스킵!!!
		if publisherCount[book.Publisher] > 0 {
			continue
		}

		hasKeyword := containsITKeyword(book.Title, keywords) || containsITKeyword(book.Cnt, keywords)
		if !hasKeyword {
			continue
		}

		if isValid, recommendUser, keywords := rs.GeminiforBook(book); isValid {
			book.RecommendUser = recommendUser
			book.Keywords = keywords
			ResultBooks = append(ResultBooks, book)
			publisherCount[book.Publisher]++ // 출판사 카운트 증가!!!

			if len(ResultBooks) >= 2 {
				break
			}
		}
	}

	// 최종 응답 생성
	return &BookSearchResponse{
		Documents:  ResultBooks,
		TotalCount: fmt.Sprintf("%d", len(ResultBooks)),
		Result:     "SUCCESS",
	}
}

// 24시간 유효기간 종료마다 캐시 자동 갱신
func (rs *RagServer) StartAutoRefresh() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			duration := next.Sub(now)

			time.Sleep(duration)
			// 캐시갱신
			response := rs.fetchLatestBooks()
			if err := rs.casheTodayBooks(response); err != nil {
				log.Printf("자동 캐시 갱신 실패: %v", err)
			}
		}
	}()
}
