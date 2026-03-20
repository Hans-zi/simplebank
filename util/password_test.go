package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := createPassword()
	hashed_password1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed_password1)

	err = CheckPassword(hashed_password1, password)
	require.NoError(t, err)

	wrong_password := createPassword()
	err = CheckPassword(hashed_password1, wrong_password)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashed_password2, err := HashPassword(hashed_password1)
	err = CheckPassword(hashed_password2, password)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}

func createPassword() string {
	return RandomString(10)
}
