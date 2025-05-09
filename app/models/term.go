package models

import "time"

type (
	TermId          string
	TermUserId      UserId
	TermName        string
	TermDescription string
)

type Term struct {
	ID          TermId
	FKUserId    TermUserId
	Name        TermName
	Description TermDescription
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
