package models

import (
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/takuchi17/term-keeper/app/models/queries"
	"golang.org/x/crypto/bcrypt"
)

type UserId string
type UserName string
type Email string
type Password string
type HashedPassword string

type User struct {
	ID        UserId
	Name      UserName
	Email     Email
	Password  HashedPassword
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateUser(
	name UserName,
	email Email,
	password Password,
) error {
	// cheack required fields
	if name == "" {
		return errors.New("username is required")
	}
	if email == "" {
		return errors.New("umail is required")
	}
	if password == "" {
		return errors.New("password is required")
	}
	// generate ulid for userId
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	userId := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	// generate hash from plain password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to generate hash", "err", err)
		return err
	}

	_, err = DB.Exec(queries.CreateUser, userId, name, email, hashedPassword, t, t)
	if err != nil {
		slog.Error("Failed to register user", "err", err)
		return err
	}
	return nil
}

func IsDuplicateEmail(email Email) (bool, error) {
	var count int
	err := DB.QueryRow(queries.IsDupulicateEmail, email).Scan(&count)
	if err != nil {
		slog.Error("Failed to check duplicate of email", "err", err)
		return false, err
	}
	return count > 0, nil
}

func IsSamePassword(hashedPassword HashedPassword, plainPassword Password) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func GetUserById(id UserId) (*User, error) {
	var (
		name      UserName
		email     Email
		createdAt time.Time
		updatedAt time.Time
	)
	err := DB.QueryRow(queries.GetUserById, id).Scan(&name, &email, &createdAt, &updatedAt)

	if err != nil {
		slog.Error("Failed to get user by id", "err", err)
		return nil, err
	}

	return &User{ID: id, Name: name, Email: email, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}

func GetUserByEmail(email Email) (*User, error) {
	var (
		name      UserName
		id        UserId
		password  HashedPassword
		createdAt time.Time
		updatedAt time.Time
	)
	err := DB.QueryRow(queries.GetUserByEmail, email).Scan(&name, &id, &password, &createdAt, &updatedAt)

	if err != nil {
		slog.Error("Failed to get user by email", "err", err)
		return nil, err
	}

	return &User{ID: id, Name: name, Email: email, Password: password, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}
