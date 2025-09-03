package db

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func dbSource() string {
	if s := os.Getenv("TEST_DB_SOURCE"); s != "" {
		return s
	}
	return "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
}

func testContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// DB에 연결 시도 — 실패하면 테스트 스킵
func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	src := dbSource()

	conn, err := sql.Open("postgres", src)
	if err != nil {
		t.Skipf("db 접속 실패 (sql.Open): %v — 테스트를 스킵합니다", err)
	}

	// 실제 연결 가능한지 체크
	ctx, cancel := testContext()
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		t.Skipf("db ping 실패: %v — 테스트를 스킵합니다", err)
	}

	return conn
}

func TestAccount_CRUD_List_Update_Delete(t *testing.T) {
	conn := setupDB(t)
	defer conn.Close()

	q := New(conn)

	// 1) Create
	ctx, cancel := testContext()
	defer cancel()

	owner := "testuser_" + randomString(6)
	createParams := CreateAccountParams{
		Owner:    owner,
		Balance:  10000,
		Currency: "USD",
	}

	created, err := q.CreateAccount(ctx, createParams)
	require.NoError(t, err, "CreateAccount 에러 발생")
	require.NotZero(t, created.ID, "생성된 ID가 0")
	require.Equal(t, owner, created.Owner)
	require.Equal(t, int64(10000), created.Balance)
	require.Equal(t, "USD", created.Currency)

	// 안전한 정리
	createdID := created.ID
	defer func() {
		_ = q.DeleteAccount(context.Background(), createdID)
	}()

	// 2) Get
	got, err := q.GetAccount(ctx, created.ID)
	require.NoError(t, err, "GetAccount 에러 발생")
	require.Equal(t, created.ID, got.ID)
	require.Equal(t, created.Owner, got.Owner)

	// 3) List
	listParams := ListAccountsParams{
		Limit:  10,
		Offset: 0,
	}
	accounts, err := q.ListAccounts(ctx, listParams)
	require.NoError(t, err, "ListAccounts 에러 발생")
	require.NotEmpty(t, accounts, "ListAccounts 결과가 비어 있음")

	found := false
	for _, a := range accounts {
		if a.ID == created.ID {
			found = true
			break
		}
	}
	require.True(t, found, "ListAccounts에 생성된 계정이 없음")

	// 5) UpdateAccountBalance
	setParams := UpdateAccountBalanceParams{
		ID:      created.ID,
		Balance: 12345,
	}
	updated2, err := q.UpdateAccountBalance(ctx, setParams)
	require.NoError(t, err, "UpdateAccountBalance 에러 발생")
	require.Equal(t, int64(12345), updated2.Balance)

	// 6) DeleteAccount
	err = q.DeleteAccount(ctx, created.ID)
	require.NoError(t, err, "DeleteAccount 에러 발생")

	// 7) 삭제 확인
	_, err = q.GetAccount(ctx, created.ID)
	require.Error(t, err, "삭제 후 GetAccount가 성공함 (원치않음)")
	require.Equal(t, sql.ErrNoRows, err, "삭제 후 예상 에러와 불일치")
}
