package token

import (
	"testing"
	"time"

	"github.com/JeongWoo-Seo/simpleBank/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
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
	require.Equal(t, payload.Role, role)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
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

func TestInvalidJWTToken(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

var benchmarkMaker Maker
var benchmarkSecretKey string = util.RandomString(32)

func init() {
	var err error
	benchmarkMaker, err = NewJWTMaker(benchmarkSecretKey)
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------------------
// 1. 토큰 생성 속도 측정
// ----------------------------------------------------

func BenchmarkCreateToken(b *testing.B) {
	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	// b.N은 테스트 프레임워크가 결정하는 반복 횟수입니다.
	for i := 0; i < b.N; i++ {
		_, _, err := benchmarkMaker.CreateToken(username, role, duration)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ----------------------------------------------------
// 2. 토큰 검증 속도 측정
// ----------------------------------------------------

func BenchmarkVerifyToken(b *testing.B) {
	b.StopTimer()
	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	token, _, err := benchmarkMaker.CreateToken(username, role, duration)
	require.NoError(b, err)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := benchmarkMaker.VerifyToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}
