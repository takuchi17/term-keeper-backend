package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
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
			err := CreateUser(tc.username, tc.email, tc.password)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")

			var (
				id        UserId
				name      UserName
				email     Email
				hashedPwd Password
				createdAt time.Time
				updatedAt time.Time
			)

			DB.QueryRow(`
    					SELECT id, name, email, password, created_at, updated_at
    					FROM users WHERE email = ?
						`, tc.email).
				Scan(&id, &name, &email, &hashedPwd, &createdAt, &updatedAt)

			assert.NoError(t, err, "Failed to get just created user.")

			assert.NotEmpty(t, id, "The user id is empty.")
			assert.Equal(t, tc.username, name, "The user name doesn't match the before.")

			assert.Equal(t, tc.email, email, "The email doesn't match the before.")

			err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(tc.password))
			assert.NoError(t, err, "Password hash does not match original.")

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

func TestGetUserById(t *testing.T) {
	ctx := context.Background()
	container, err := SetupMysqlContainerAndSetupDB(t, &ctx)

	require.NoError(t, err, "Failed to setup tester.")
	defer container.Terminate(ctx)
	defer DB.Close()

	rows, err := DB.QueryContext(ctx, "SELECT id, name, email FROM users")
	require.NoError(t, err, "Failed to query users table.")
	defer rows.Close()

	t.Log("--- users テーブルの内容 ---")
	for rows.Next() {
		var id, name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			t.Logf("Error scanning user row: %v", err)
			continue
		}
		t.Logf("ID: %s, Name: %s, Email: %s", id, name, email)
	}
	t.Log("-------------------------")

	testCases := []struct {
		name      string
		id        UserId
		wantName  UserName
		wantEmail Email
		wantErr   bool
	}{
		{
			name:      "Found user",
			id:        "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			wantName:  "山田太郎",
			wantEmail: "yamada@example.com",
			wantErr:   false,
		},
		{
			name:    "Not found user",
			id:      "notnotnotnotnotnotnotnotno",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := GetUserById(tc.id)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")
			assert.NotNil(t, user, "User should not be nil")

			if user != nil {
				assert.Equal(t, tc.wantName, user.Name, "Username mismatch")
				assert.Equal(t, tc.wantEmail, user.Email, "Email mismatch")

				assert.False(t, user.Created_at.IsZero(), "Created_at should not be zero time")
				assert.False(t, user.Updated_at.IsZero(), "Updated_at should not be zero time")
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	container, err := SetupMysqlContainerAndSetupDB(t, &ctx)

	require.NoError(t, err, "Failed to setup tester.")
	defer container.Terminate(ctx)
	defer DB.Close()

	rows, err := DB.QueryContext(ctx, "SELECT id, name, email FROM users")
	require.NoError(t, err, "Failed to query users table.")
	defer rows.Close()

	t.Log("--- users テーブルの内容 ---")
	for rows.Next() {
		var id, name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			t.Logf("Error scanning user row: %v", err)
			continue
		}
		t.Logf("ID: %s, Name: %s, Email: %s", id, name, email)
	}
	t.Log("-------------------------")

	testCases := []struct {
		name     string
		email    Email
		wantId   UserId
		wantName UserName
		wantErr  bool
	}{
		{
			name:     "Found user",
			email:    "yamada@example.com",
			wantName: "山田太郎",
			wantId:   "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			wantErr:  false,
		},
		{
			name:    "Not found user",
			email:   "no@email.jp",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := GetUserByEmail(tc.email)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")
			assert.NotNil(t, user, "User should not be nil")

			if user != nil {
				assert.Equal(t, tc.wantName, user.Name, "Username mismatch")
				assert.Equal(t, tc.wantId, user.ID, "ID mismatch")

				assert.False(t, user.Created_at.IsZero(), "Created_at should not be zero time")
				assert.False(t, user.Updated_at.IsZero(), "Updated_at should not be zero time")
			}
		})
	}
}
