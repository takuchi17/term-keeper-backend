package models

import (
	"errors"
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

type User struct {
	ID         UserId
	Name       UserName
	Email      Email
	Password   Password
	Created_at time.Time
	Updated_at time.Time
}

func RegisterUser(
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
		return err
	}

	_, err = DB.Exec(queries.RegisterUser, userId, name, email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func IsDuplicateEmail(email Email) (bool, error) {
	var count int
	err := DB.QueryRow(queries.IsDupulicateEmail, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
