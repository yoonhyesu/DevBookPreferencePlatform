package repository

import (
	"DBP/internal/model"
	"DBP/internal/utils"
	"DBP/pkg/database"
)

type Repository interface {
	// mariadb
	// 홈화면(개발자추천책, 공지사항)
	GetHomeDevBook() model.HomeDevBook
	GetNotice() []model.Notice

	// 회원가입
	Register(req model.RegisterRequest) error
	CheckDuplicateID(req model.RegisterRequest) error

	// 메뉴(분야별,개발자별,단계별추천)
	GetTagList() []model.TagList
	GetTagWithBookList(tagID int) []model.TagWithBook
	GetBookDetail(bookID int) []model.BookInfo
	///'......ㅇㄴ() []model.DevInfo
	GetDevInfo(devID string) []model.DevDetail
	GetDevRecommendedBooks(devID string) []model.DevBookLikeReason
	GetStepWithBook(step string) []model.BookInfo

	// 관리자(개발자관리)
	GetDevList() []model.AddDevs

	AddDev(dev model.AddDevs) error
	UpdateDev(dev model.AddDevs) error
	DeleteDev(devID string) error
	GetDevID(devID string) (model.AddDevs, error)
	// 관리자(공지사항관리)
	GetNoticeList() []model.Notice
	AddNotice(notice model.Notice) error
	UpdateNotice(notice model.Notice) error
	DeleteNotice(noticeID int) error

	// 관리자(책관리)
	GetBookList() []model.AddBooks
	AddBook(book model.AddBooks, devRecommends []model.DevRecommend) error
	UpdateBook(book model.AddBooks, devRecommends []model.DevRecommend) error
	DeleteBook(bookID int) error

	// 관리자(태그관리)
	GetAdminTagList() []model.TagList

	// redis
	// 로그인&로그아웃
	SignIn(req model.LoginRequest) (*utils.TokenDetails, error)
	SignOut(accessToken string) error

	// 사용자 정보 조회(세션유지용)
	GetUserInfo(userId string) (*model.User, error)

	// 사용자 관리
	GetUserProfile(userID string) (*model.User, error)
	UpdateProfile(req model.ProfileEditRequest) error
	ChangePassword(req model.PasswordChangeRequest) error
	LeaveAccount(req model.LeaveRequest) error
	CheckCurrentPassword(userID string, currentPassword string) error
}

// 공통 레포지토리 구조체
type CommonRepo struct {
	mariaDB *database.MariaDBHandler
	redis   *database.RedisHandler
}

// 싱글톤 패턴으로 repository 인스턴스 관리
func NewRepository(db *database.MariaDBHandler, redis *database.RedisHandler) Repository {
	return &CommonRepo{
		mariaDB: db,
		redis:   redis,
	}
}
