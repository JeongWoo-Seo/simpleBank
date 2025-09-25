package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/JeongWoo-Seo/simpleBank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

// go test 실행 시 제일 먼저 실행됨
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config")
	}

	// DB 연결
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("❌ cannot connect to db: %v", err)
	}

	// 전역 Queries, Store 초기화
	testStore = NewStore(connPool)

	// 모든 테스트 실행
	code := m.Run()

	os.Exit(code)
}
