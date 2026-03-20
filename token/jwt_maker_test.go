package token

import (
	"testing"
	"time"

	"github.com/Hans-zi/simple_bank/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	secretKey := util.RandomString(secretKeySize)
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)
	username, duration := util.RandomOwner(), time.Minute
	issueAt := time.Now()
	expireAt := time.Now().Add(duration)

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMaker.VerifyJWTToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expireAt, payload.ExpiresAt.Time, time.Second)

}

func TestExpiredToken(t *testing.T) {
	secretKey := util.RandomString(secretKeySize)
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)
	username, duration := util.RandomOwner(), -time.Minute

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMaker.VerifyJWTToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(secretKeySize))

	payload, err := NewJTWPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = maker.VerifyJWTToken(token)
	require.Nil(t, payload)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
}
