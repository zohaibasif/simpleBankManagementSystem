package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zohaibAsif/simple_bank_management_system/util"
)

func createRandomUser(t *testing.T) User {
	args := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: "secured",
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}

	user, err := queries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {

	createRandomUser(t)

}

func TestGetUser(t *testing.T) {

	user := createRandomUser(t)

	response, err := queries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)

	require.NotEmpty(t, response)

	require.Equal(t, user.Username, response.Username)

	require.WithinDuration(t, user.PasswordChangedAt, response.PasswordChangedAt, time.Second)

	require.WithinDuration(t, user.CreatedAt, response.CreatedAt, time.Second)

}
