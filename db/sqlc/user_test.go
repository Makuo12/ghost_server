package db

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/makuo12/ghost_server/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	var username string = uuid.New().String()
	var firebasePassword string = uuid.New().String()
	hashedPassword, err := utils.HashedPassword(utils.RandomString(6))
	require.NoError(t, err)
	firstName := utils.RandomString(8)
	lastName := utils.RandomString(8)
	email := utils.RandomEmail()
	arg := CreateUserParams{
		HashedPassword:   hashedPassword,
		Email:            strings.ToLower(email),
		Username:         username,
		FirebasePassword: firebasePassword,
		DateOfBirth:      time.Now(),
		FirstName:        strings.ToLower(strings.TrimSpace(firstName)),
		LastName:         strings.ToLower(strings.TrimSpace(lastName)),
		Currency:         utils.NGN,
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
