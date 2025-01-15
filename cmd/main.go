package main

import (
	"DBP/config"
	"DBP/internal/router"
	"DBP/pkg/database"
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

type (
	Database struct {
		Connection *sql.DB
	}
)

func main() {
	r := gin.Default()

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
