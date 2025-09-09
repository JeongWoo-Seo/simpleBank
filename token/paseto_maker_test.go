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
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	// 정상 payload 생성
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// 잘못된 secretKey로 PasetoMaker 생성
	wrongMaker, err := NewPasetoMaker(util.RandomString(32)) // 공격자가 쓴 키
	require.NoError(t, err)

	// 공격자가 잘못된 key로 토큰 생성
	token, err := wrongMaker.CreateToken(payload.Username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// 서버가 가진 올바른 secretKey로 PasetoMaker 생성
	correctMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	// 서버에서 검증 시도 → 실패해야 함
	payload, err = correctMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
