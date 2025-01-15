package repository

import (
	"DBP/internal/model"
	"DBP/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// 회원가입 시 아이디 중복체크 -> 사용자에게 보여줄 용도
func (m *CommonRepo) CheckDuplicateID(req model.RegisterRequest) error {
	// 디버깅을 위한 로그 추가
	log.Printf("검사할 ID: %s", req.UserID)

	var count int
	err := m.mariaDB.Connection.QueryRow("SELECT COUNT(*) FROM users WHERE user_id = ?", req.UserID).Scan(&count)
	if err != nil {
		log.Printf("데이터베이스 조회 에러: %v", err)
		return fmt.Errorf("데이터베이스 오류: %v", err)
	}

	if count > 0 {
		return errors.New("이미 존재하는 아이디 입니다")
	}

	return nil
}

// 회원가입
func (m *CommonRepo) Register(req model.RegisterRequest) error {
	// 트랜잭션 처리
	tx, err := m.mariaDB.Connection.Begin()
	if err != nil {
		return fmt.Errorf("트랜잭션 시작 실패: %v", err)
	}
	defer tx.Commit() // 함수 종료 전에 롤백 설정!(커밋되면 롤백은 무시됨)

	// 아이디 중복 체크
	checkUserIDQuery := `
	SELECT COUNT(*) FROM dbp.users WHERE LOWER(USER_ID) = LOWER(?)
	`
	var userIdCount int
	err = tx.QueryRow(checkUserIDQuery, req.UserID).Scan(&userIdCount) // 쿼리에 요소를 넣고 결과를 스캔해 userIdCount에 반환함.
	if err != nil {
		return fmt.Errorf("아이디 중복 확인 중 오류발생: %v", err)
	}
	if userIdCount > 0 {
		return fmt.Errorf("이미 존재하는 아이디 입니다")
	}

	// 이메일 중복 체크
	checkUserEmailQuery := `
	SELECT COUNT(*) FROM dbp.users WHERE LOWER(EMAIL ) = LOWER(?)
	`
	var userEmailCount int
	err = tx.QueryRow(checkUserEmailQuery, req.Email).Scan(&userEmailCount)
	if err != nil {
		return fmt.Errorf("이메일 중복 확인 중 오류발생: %v", err)
	}
	if userEmailCount > 0 {
		return fmt.Errorf("이미 존재하는 이메일 입니다")
	}

	// 전화번호 중복 체크
	checkPhoneQuery := `
	SELECT COUNT(*) FROM dbp.users WHERE LOWER(PHONE_NUMBER) = LOWER(?)
	`
	var PhoneCount int
	err = tx.QueryRow(checkPhoneQuery, req.PhoneNumber).Scan(&PhoneCount)
	if err != nil {
		return fmt.Errorf("전화번호 중복 확인 중 오류발생: %v", err)
	}
	if PhoneCount > 0 {
		return fmt.Errorf("이미 존재하는 전화번호 입니다")
	}

	// 비밀번호 해싱
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("해싱 실패")
		return err
	}

	// 회원가입 쿼리 실행
	query := `
	INSERT INTO dbp.users
(id, USER_ID, PASSWORD, USER_NAME, USER_STATUS, EMAIL, PHONE_NUMBER, CREATE_DATE, UPDATE_DATE, IS_ADMIN)
VALUES(null, ?, ?, ?, 0, ?, ?, ?, ?, 0)
	`
	_, err = tx.Exec(query,
		req.UserID,
		hashedPassword,
		req.UserName,
		req.Email,
		req.PhoneNumber,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Println("회원가입 처리 결과에 오류가 있습니다")
	}
	// 모든 작업 성공 시 커밋
	return tx.Commit()
}

// 로그인
func (m *CommonRepo) SignIn(req model.LoginRequest) (*utils.TokenDetails, error) {
	// 아이디를 기준으로 비밀번호가 맞는지 확인
	IdCheckQuery := `
	SELECT USER_ID, PASSWORD, EMAIL
	FROM dbp.users
	WHERE USER_STATUS = 0
	AND USER_ID = ?
	`
	var hashedPassword string
	err := m.mariaDB.Connection.QueryRow(IdCheckQuery, req.UserId).Scan(
		&req.UserId,
		&hashedPassword, // DB의 해시된 비밀번호
		&req.Email,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("존재하지 않는 아이디입니다")
	}

	if err != nil {
		return nil, fmt.Errorf("로그인 처리 중 오류 발생:%v", err)
	}

	// 입력받은 비밀번호와 해시된 비밀번호 비교
	if err := utils.ComparePassword(req.Password, hashedPassword); err != nil {
		return nil, fmt.Errorf("비밀번호가 일치하지 않습니다")
	}

	// 사용자 전체 정보 조회
	query := `
	SELECT USER_ID, PASSWORD, EMAIL, IS_ADMIN, USER_NAME
	FROM dbp.users
	WHERE USER_STATUS = 0
	AND USER_ID = ?
	`
	err = m.mariaDB.Connection.QueryRow(query, req.UserId).Scan(
		&req.UserId,
		&hashedPassword,
		&req.Email,
		&req.IsAdmin,
		&req.UserName,
	)
	if err != nil {
		return nil, fmt.Errorf("사용자 정보 조회 실패: %v", err)
	}

	// 토큰 생성 시 isAdmin 정보도 전달
	tokens, err := utils.GenerateTokens(req.UserId, req.Email, req.IsAdmin, req.UserName)
	if err != nil {
		return nil, fmt.Errorf("토큰 생성 실패: %v", err)
	}

	// 마지막 로그인 시간 업데이트
	updateQuery := `
	UPDATE dbp.users 
	SET LAST_LOGIN_DATE = NOW() 
	WHERE USER_ID = ?
	`
	_, err = m.mariaDB.Connection.Exec(updateQuery, req.UserId)
	if err != nil {
		log.Printf("마지막 로그인 시간 업데이트 실패: %v", err)
		// 업데이트 실패해도 로그인은 허용
	}

	// Refresh Token만 Redis에 저장
	err = m.redis.Client.Set(
		context.Background(),
		tokens.RefreshUuid,
		req.UserId,
		time.Until(time.Unix(tokens.RtExpires, 0)),
	).Err()

	if err != nil {
		return nil, fmt.Errorf("토큰 저장 실패: %v", err)
	}
	return tokens, nil
}

// 로그아웃
func (m *CommonRepo) SignOut(accessToken string) error {
	// access token을 블랙리스트에 추가
	err := m.redis.Client.Set(
		context.Background(),
		"blacklist:"+accessToken,
		"revoked",
		30*time.Minute, // access token의 만료 시간과 동일하게 설정함  why?
	).Err()

	if err != nil {
		return fmt.Errorf("로그아웃 처리 중 오류 발생: %v", err)
	}
	return nil
	// 여기서는 템플릿데이터 어캐하지/
}

// 사용자 정보 조회
func (m *CommonRepo) GetUserInfo(userId string) (*model.User, error) {
	query := `
	SELECT USER_ID, EMAIL, USER_NAME, IS_ADMIN 
FROM dbp.users 
WHERE USER_ID = ?
AND USER_STATUS = 0
	`
	var user model.User
	err := m.mariaDB.Connection.QueryRow(query, userId).Scan(
		&user.UserId,
		&user.Email,
		&user.UserName,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("사용자 정보 조회 실패; %v", err)
	}
	return &user, nil
}

// 사용자 정보 조회
func (m *CommonRepo) GetUserProfile(userID string) (*model.User, error) {
	var user model.User
	// NULL이 가능한 필드를 위한 변수들
	var userText, company, githubLink, blogLink sql.NullString
	var phoneNumber sql.NullString

	err := m.mariaDB.Connection.QueryRow(`
	SELECT PHONE_NUMBER, USER_TEXT, COMPANY, GITHUB_LINK, BLOG_LINK
	FROM dbp.users
	WHERE USER_ID = ? AND USER_STATUS = 0
	`, userID).Scan(
		&phoneNumber,
		&userText,
		&company,
		&githubLink,
		&blogLink,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("사용자를 찾을 수 없습니다")
		}
		return nil, fmt.Errorf("프로필 조회 중 오류 발생: %v", err)
	}

	// NULL 체크 후 값 할당
	user.PhoneNumber = phoneNumber.String
	user.UserText = userText.String
	user.Company = company.String
	user.GithubLink = githubLink.String
	user.BlogLink = blogLink.String

	return &user, nil
}

// 프로필 수정
func (m *CommonRepo) UpdateProfile(req model.ProfileEditRequest) error {
	_, err := m.mariaDB.Connection.Exec(`
	UPDATE dbp.users 
	SET 
		USER_NAME = ?,
		PHONE_NUMBER = ?,
		USER_TEXT = ?,
		COMPANY = ?,
		GITHUB_LINK = ?,
		BLOG_LINK = ?,
		UPDATE_DATE = NOW()
	WHERE USER_ID = ? AND USER_STATUS = 0
	`, req.UserName, req.PhoneNumber, req.UserText, req.Company,
		req.GithubLink, req.BlogLink, req.UserID)

	if err != nil {
		return fmt.Errorf("프로필 수정 실패: %v", err)
	}
	return nil
}

// 비밀번호 변경
func (m *CommonRepo) ChangePassword(req model.PasswordChangeRequest) error {
	// 1. 현재 비밀번호 확인
	var storedHash string
	err := m.mariaDB.Connection.QueryRow(`
	SELECT PASSWORD FROM dbp.users 
	WHERE USER_ID = ? AND USER_STATUS = 0`,
		req.UserID).Scan(&storedHash)
	if err != nil {
		return fmt.Errorf("사용자 조회 실패: %v", err)
	}

	// 2. 현재 비밀번호 검증
	if err := utils.ComparePassword(req.OldPassword, storedHash); err != nil {
		return fmt.Errorf("현재 비밀번호가 일치하지 않습니다")
	}

	// 3. 새 비밀번호 해시화
	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("비밀번호 해시화 실패: %v", err)
	}

	// 4. 새 비밀번호로 업데이트
	_, err = m.mariaDB.Connection.Exec(`
	UPDATE dbp.users 
	SET PASSWORD = ?, UPDATE_DATE = NOW()
	WHERE USER_ID = ? AND USER_STATUS = 0`,
		newHash, req.UserID)
	if err != nil {
		return fmt.Errorf("비밀번호 변경 실패: %v", err)
	}

	return nil
}

// 회원 탈퇴
func (m *CommonRepo) LeaveAccount(req model.LeaveRequest) error {

	_, err := m.mariaDB.Connection.Exec(`
	UPDATE dbp.users 
	SET USER_STATUS = 1, UPDATE_DATE = NOW()
	WHERE USER_ID = ?`,
		req.UserID)
	if err != nil {
		return fmt.Errorf("회원 탈퇴 처리 실패: %v", err)
	}

	return nil
}

// 현재 비밀번호 확인
func (m *CommonRepo) CheckCurrentPassword(userID string, currentPassword string) error {
	var storedHash string
	err := m.mariaDB.Connection.QueryRow(`
	SELECT PASSWORD 
	FROM dbp.users 
	WHERE 
	USER_ID = ?
	AND USER_STATUS  = 0`,
		userID).Scan(&storedHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("사용자를 찾을 수 없습니다")
		}
		return fmt.Errorf("비밀번호 확인 중 오류 발생: %v", err)
	}

	// 비밀번호 검증
	if err := utils.ComparePassword(currentPassword, storedHash); err != nil {
		return fmt.Errorf("비밀번호가 일치하지 않습니다")
	}

	return nil
}
