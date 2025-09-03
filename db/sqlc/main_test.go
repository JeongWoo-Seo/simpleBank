package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB
var testQueries *Queries

// DB 연결 소스 (환경변수 없으면 기본값 사용)
func testDBSource() string {
	source := os.Getenv("TEST_DB_SOURCE")
	if source == "" {
		source = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	}
	return source
}

// go test 실행 시 제일 먼저 실행됨
func TestMain(m *testing.M) {
	var err error

	// DB 연결
	testDB, err = sql.Open("postgres", testDBSource())
	if err != nil {
		log.Fatalf("❌ cannot connect to db: %v", err)
	}

	// 전역 Queries, Store 초기화
	testQueries = New(testDB)

	// 모든 테스트 실행
	code := m.Run()

	// 종료 시 DB 닫기
	_ = testDB.Close()

	os.Exit(code)
}
