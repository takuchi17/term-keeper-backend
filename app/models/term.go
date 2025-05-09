package models

import (
	"errors"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/takuchi17/term-keeper/app/models/queries"
)

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

func CreateTerm(userId TermUserId, name TermName, description TermDescription) (*Term, error) {
	// cheack required fields
	if name == "" {
		return nil, errors.New("termname is required")
	}
	// generate ulid for termId
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	termId := TermId(ulid.MustNew(ulid.Timestamp(t), entropy).String())

	_, err := DB.Exec(queries.CreateTerm, termId, userId, name, description, t, t)
	if err != nil {
		slog.Error("Failed to create a term", "err", err)
		return nil, err
	}

	return &Term{
			ID:          termId,
			FKUserId:    userId,
			Name:        name,
			Description: description,
			CreatedAt:   t,
			UpdatedAt:   t,
		},
		nil
}

func GetTermsByUserId(userId TermUserId, query *string, category *string, sort *string, checked *bool) ([]*Term, error) {
	var sb strings.Builder
	var args []interface{}

	sb.WriteString(queries.GetTermsByUserIdBase)

	if category != nil {
		sb.WriteString(queries.GetTermsJoinWithCategory)
	}

	sb.WriteString(queries.GetTermsByUserIdWhere)
	args = append(args, userId)

	if query != nil {
		sb.WriteString(queries.GetTermsFillterByName)
		args = append(args, "%"+*query+"%")
	}
	if category != nil {
		sb.WriteString(queries.GetTermsFillterByCategory)
		args = append(args, *category)
	}
	if checked != nil {
		sb.WriteString(queries.GetTermsFillterByChecked)
		args = append(args, *checked)
	}

	// ソート
	if sort != nil {
		switch *sort {
		case "created_at_asc":
			sb.WriteString(queries.GetTermsSortByCreatedAsc)
		case "created_at_desc":
			sb.WriteString(queries.GetTermsSortByCreatedDesc)
		case "updated_at_asc":
			sb.WriteString(queries.GetTermsSortByUpdatedAsc)
		case "updated_at_desc":
			sb.WriteString(queries.GetTermsSortByCreatedAsc)
		case "term_asc":
			sb.WriteString(queries.GetTermsSortByNameAsc)
		case "term_desc":
			sb.WriteString(queries.GetTermsSortByNameDesc)
		default:
		}
	}

	rows, err := DB.Query(sb.String(), args...)
	if err != nil {
		slog.Error("Failed to get terms", "err", err)
		return nil, err
	}
	defer rows.Close()

	var terms []*Term
	for rows.Next() {
		var term Term
		if err := rows.Scan(&term.ID, &term.FKUserId, &term.Name, &term.Description, &term.CreatedAt, &term.UpdatedAt); err != nil {
			slog.Error("Failed to scan term", "err", err)
			return nil, err
		}
		terms = append(terms, &term)
	}

	return terms, nil
}

func (t *Term) Update() (*Term, error) {
	if t.Name == "" {
		return nil, errors.New("termname is required")
	}

	_, err := DB.Exec(queries.UpdateTerm, t.Name, t.Description, t.UpdatedAt, t.ID)
	if err != nil {
		slog.Error("Failed to update term", "err", err)
		return nil, err
	}

	return t, nil
}
