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

	payload, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.RegisteredClaims.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expireAt, payload.RegisteredClaims.ExpiresAt.Time, time.Second)

}

func TestExpiredToken(t *testing.T) {
	secretKey := util.RandomString(secretKeySize)
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)
	username, duration := util.RandomOwner(), -time.Minute

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(secretKeySize))

	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Nil(t, payload)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
}
