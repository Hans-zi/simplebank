package token

import (
	"testing"
	"time"

	"github.com/Hans-zi/simple_bank/util"
	"github.com/aead/chacha20poly1305"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	symmetricKey := util.RandomString(chacha20poly1305.KeySize)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username, duration := util.RandomOwner(), time.Minute
	issueAt := time.Now()
	expiresAt := time.Now().Add(duration)
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(secretKeySize))
	require.NoError(t, err)
	username, duration := util.RandomOwner(), -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
}
