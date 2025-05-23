package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
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
			email:    "yamada@example.com",
			password: "password",
			wantErr:  true,
		},
		{
			name:     "Duplicate username allowed",
			username: "example1",
			email:    "example_duplicate_username@gmail.com",
			password: "password",
			wantErr:  false,
		},
		{
			name:     "Invalid email format",
			username: "invalidemail",
			email:    "invalid-email",
			password: "password",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			err = CreateUser(tx, tc.username, tc.email, tc.password)

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

			tx.QueryRow(`
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
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			isDuplicate, err := IsDuplicateEmail(tx, tc.email)

			assert.NoError(t, err, "Expected no error, but an error occurred.")

			assert.Equal(t, tc.wantDup, isDuplicate, "Unexpected duplication check result.")
		})
	}
}

func TestGetUserById(t *testing.T) {
	rows, err := DB.Query("SELECT id, name, email FROM users")
	require.NoError(t, err, "Failed to query users table.")
	defer rows.Close()

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
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			user, err := GetUserById(tx, tc.id)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")
			assert.NotNil(t, user, "User should not be nil")

			if user != nil {
				assert.Equal(t, tc.wantName, user.Name, "Username mismatch")
				assert.Equal(t, tc.wantEmail, user.Email, "Email mismatch")

				assert.False(t, user.CreatedAt.IsZero(), "Created_at should not be zero time")
				assert.False(t, user.UpdatedAt.IsZero(), "Updated_at should not be zero time")
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
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
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			user, err := GetUserByEmail(tx, tc.email)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")
			assert.NotNil(t, user, "User should not be nil")

			if user != nil {
				assert.Equal(t, tc.wantName, user.Name, "Username mismatch")
				assert.Equal(t, tc.wantId, user.ID, "ID mismatch")

				assert.False(t, user.CreatedAt.IsZero(), "Created_at should not be zero time")
				assert.False(t, user.UpdatedAt.IsZero(), "Updated_at should not be zero time")
			}
		})
	}
}

func TestIsSamePassword(t *testing.T) {
	raw := "mypassword"
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	require.NoError(t, err)

	t.Run("Correct password", func(t *testing.T) {
		err := IsSamePassword(DB, HashedPassword(hashed), Password(raw))
		assert.NoError(t, err)
	})

	t.Run("Wrong password", func(t *testing.T) {
		err := IsSamePassword(DB, HashedPassword(hashed), Password("wrongpass"))
		assert.Error(t, err)
	})
}
