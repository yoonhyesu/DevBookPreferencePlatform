package database

import (
	"DBP/config"
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type MariaDBHandler struct {
	Connection *sql.DB
}

var (
	instance *MariaDBHandler
	once     sync.Once
)

// 데이터베이스 연결 풀의 인스턴스를 싱글톤패턴으로 유지 -sync.Once사용
func GetDatabaseInstance() *MariaDBHandler {
	once.Do(func() {
		config.LoadEnv()
		username := os.Getenv("MARIADB_USERNAME")
		password := os.Getenv("MARIADB_PASSWORD")
		host := os.Getenv("MARIADB_HOSTNAME")
		port := os.Getenv("MARIADB_PORT")
		dbname := os.Getenv("MARIADB_DBNAME")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", username, password, host, port, dbname)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		instance = &MariaDBHandler{Connection: db}
	})
	return instance
}
