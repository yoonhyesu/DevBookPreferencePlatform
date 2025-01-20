package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Domain string
}

var (
	config *Config
	once   sync.Once
)

// 싱글톤 패턴으로 Config 인스턴스 관리
func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			Domain: getOrDefault("DOMAIN", "localhost"),
		}
	})
	return config
}

// 환경변수 값을 가져오거나 기본값 반환
func getOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadEnv() error {
	// 현재 실행 파일의 경로를 기준으로 상위 디렉토리로 이동
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("현재 디렉토리 확인 실패: %v", err)
		return err
	}

	// cmd 폴더에서 실행되는 경우 상위 디렉토리로 이동
	if filepath.Base(currentDir) == "cmd" {
		currentDir = filepath.Dir(currentDir)
	}

	envPath := filepath.Join(currentDir, ".env")
	log.Printf("환경변수 파일 경로: %s", envPath) // 디버깅을 위한 로그 추가

	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("환경변수 로드 실패 (%s): %v", envPath, err)
		return err
	}

	log.Println("환경변수 로드 성공")
	return nil
}

// DSN을 반환하는 함수
func GetDSN() string {
	LoadEnv()
	username := os.Getenv("MARIADB_USERNAME")
	password := os.Getenv("MARIADB_PASSWORD")
	host := os.Getenv("MARIADB_HOSTNAME")
	port := os.Getenv("MARIADB_PORT")
	dbname := os.Getenv("MARIADB_DBNAME")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", username, password, host, port, dbname)
}
