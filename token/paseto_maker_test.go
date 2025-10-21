package token

import (
	"testing"
	"time"

	"github.com/JeongWoo-Seo/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), util.DepositorRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	// 정상 payload 생성
	payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// 잘못된 secretKey로 PasetoMaker 생성
	wrongMaker, err := NewPasetoMaker(util.RandomString(32)) // 공격자가 쓴 키
	require.NoError(t, err)

	// 공격자가 잘못된 key로 토큰 생성
	token, _, err := wrongMaker.CreateToken(payload.Username, util.DepositorRole, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// 서버가 가진 올바른 secretKey로 PasetoMaker 생성
	correctMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = correctMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

var benchmarkPasetoMaker Maker
var benchmarkPasetoSecretKey string = util.RandomString(32)

func init() {
	var err error

	benchmarkPasetoMaker, err = NewPasetoMaker(benchmarkPasetoSecretKey)
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------------------
// 1. PASETO 토큰 생성 속도 측정
// ----------------------------------------------------

func BenchmarkPasetoCreateToken(b *testing.B) {
	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	for i := 0; i < b.N; i++ {
		_, _, err := benchmarkPasetoMaker.CreateToken(username, role, duration)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ----------------------------------------------------
// 2. PASETO 토큰 검증 속도 측정
// ----------------------------------------------------

func BenchmarkPasetoVerifyToken(b *testing.B) {
	b.StopTimer()
	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	token, _, err := benchmarkPasetoMaker.CreateToken(username, role, duration)
	require.NoError(b, err)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := benchmarkPasetoMaker.VerifyToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}
