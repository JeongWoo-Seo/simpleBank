package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/JeongWoo-Seo/simpleBank/util"
	_ "github.com/lib/pq"
)

var testDB *sql.DB
var testStore Store

// go test 실행 시 제일 먼저 실행됨
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config")
	}

	// DB 연결
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("❌ cannot connect to db: %v", err)
	}

	// 전역 Queries, Store 초기화
	testStore = NewStore(testDB)

	// 모든 테스트 실행
	code := m.Run()

	// 종료 시 DB 닫기
	_ = testDB.Close()

	os.Exit(code)
}
