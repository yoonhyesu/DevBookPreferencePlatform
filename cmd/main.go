package main

import (
	"DBP/config"
	"DBP/internal/router"
	"DBP/pkg/database"
	"context"
	"database/sql"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type (
	Database struct {
		Connection *sql.DB
	}
)

func main() {
	// 반환된 *gin.Engine 구조체는 REST API 요청을 핸들링하는 method 등을 구현함
	r := gin.Default()

	gin.DisableConsoleColor()

	// 파일에 로그를 작성합니다.
	f, _ := os.Create("DBP.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	// 환경변수 로드
	config.LoadEnv()
	// db 연결 생성
	db := database.GetDatabaseInstance()
	defer db.Connection.Close()

	// Redis 연결
	redis := database.GetRedisInstance()
	defer redis.Client.Close()

	// Redis 연결 확인
	ctx := context.Background()
	if err := redis.Client.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis 연결 실패:", err)
	}
	log.Println("Redis 연결 성공!!")

	//라우터 설정 시 db 인스턴스 전달
	router.SetupRouter(r, db, redis)
	r.Run(":7777")
}
