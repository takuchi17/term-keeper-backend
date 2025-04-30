package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()
	container, err := SetupMysqlContainerAndSetupDB(t, &ctx)

	require.NoError(t, err, "Failed to setup tester.")
	defer container.Terminate(ctx)
	defer DB.Close()

	testCases := []struct {
		name     string
		username UserName
		email    Email
		password Password
		wantErr  bool
	}{
		{
			name:     "Normal user registering",
			username: "example1",
			email:    "example1@gmail.com",
			password: "password",
			wantErr:  false,
		},
		{
			name:     "Empty user name",
			email:    "example2@gmail.com",
			password: "password",
			wantErr:  true,
		},
		{
			name:     "Empty email",
			username: "example3",
			password: "password",
			wantErr:  true,
		},
		{
			name:     "Empty password",
			username: "example4",
			email:    "example4@gmail.com",
			wantErr:  true,
		},
		{
			name:     "Dupulicate email",
			username: "example5",
			email:    "example1@gmail.com",
			password: "password",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := RegisterUser(tc.username, tc.email, tc.password)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")

			var (
				id           UserId
				name         UserName
				email        Email
				hashedPwd    Password
				createdAtStr string
				updatedAtStr string
			)

			DB.QueryRow(`
    					SELECT id, name, email, password, created_at, updated_at
    					FROM users WHERE email = ?
						`, tc.email).
				Scan(&id, &name, &email, &hashedPwd, &createdAtStr, &updatedAtStr)

			assert.NoError(t, err, "Failed to get just created user.")

			assert.NotEmpty(t, id, "The user id is empty.")
			assert.Equal(t, tc.username, name, "The user name doesn't match the before.")

			assert.Equal(t, tc.email, email, "The email doesn't match the before.")

			err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(tc.password))
			assert.NoError(t, err, "Password hash does not match original.")

			createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
			assert.NoError(t, err, "Failed to parse string created_at ( time.Time ).")

			updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
			assert.NoError(t, err, "Failed to parse string updated_at ( time.Time ).")

			assert.Equal(t, createdAt.Format("2006-01-02 15:04:05"), updatedAt.Format("2006-01-02 15:04:05"), "The created_at and updated_at timestamps are different.")
		})
	}
}

func TestIsDuplicateEmail(t *testing.T) {

	ctx := context.Background()
	container, err := SetupMysqlContainerAndSetupDB(t, &ctx)

	require.NoError(t, err, "Failed to setup tester.")
	defer container.Terminate(ctx)
	defer DB.Close()

	testCases := []struct {
		name     string
		username UserName
		email    Email
		password Password
		wantDup  bool
	}{
		{
			name:     "Duplicate email",
			username: "duplicate",
			email:    "yamada@example.com",
			password: "password",
			wantDup:  true,
		},
		{
			name:     "Not duplicate email",
			username: "notduplicate",
			email:    "new@example.com",
			password: "password",
			wantDup:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isDuplicate, err := IsDuplicateEmail(tc.email)

			assert.NoError(t, err, "Expected no error, but an error occurred.")

			assert.Equal(t, tc.wantDup, isDuplicate, "Unexpected duplication check result.")
		})
	}
}
