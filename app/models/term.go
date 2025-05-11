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
	TermCategoryId  CategoryId
)

type Term struct {
	ID          TermId
	FKUserId    TermUserId
	Name        TermName
	Description TermDescription
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TermAndCategory struct {
	Term       *Term
	Categories []*Category
}

func CreateTerm(db SQLExecutor, userId TermUserId, name TermName, description TermDescription, categoryIds []CategoryId) (*Term, error) {
	// cheack required fields
	if name == "" {
		return nil, errors.New("termname is required")
	}
	// generate ulid for termId
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	termId := TermId(ulid.MustNew(ulid.Timestamp(t), entropy).String())

	_, err := db.Exec(queries.CreateTerm, termId, userId, name, description, t, t)
	if err != nil {
		slog.Error("Failed to create a term", "err", err)
		return nil, err
	}

	// カテゴリIDを term_category_relations テーブルに挿入
	err = LinkTermWithCategories(db, termId, categoryIds)
	if err != nil {
		slog.Error("Failed to link term with categories", "err", err)
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

func GetTermsWithCategoriesByUserId(db SQLExecutor, userId TermUserId, query *string, category *string, sort *string, checked *bool) ([]*TermAndCategory, error) {
	terms, err := GetTermsByUserId(db, userId, query, category, sort, checked)
	if err != nil {
		return nil, err
	}

	var result []*TermAndCategory
	for _, term := range terms {
		// 中間テーブルからcategoryIdを取得
		categoryIds, err := GetCategoryIdsByTermId(db, term.ID)
		if err != nil {
			return nil, err
		}

		// カテゴリ情報を取得
		categories, err := GetCategoriesByIds(db, categoryIds)
		if err != nil {
			return nil, err
		}

		result = append(result, &TermAndCategory{
			Term:       term,
			Categories: categories,
		})
	}

	return result, nil
}

func GetTermsByUserId(db SQLExecutor, userId TermUserId, query *string, category *string, sort *string, checked *bool) ([]*Term, error) {
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
		sb.WriteString(queries.GetTermsFilterByChecked)
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
			sb.WriteString(queries.GetTermsSortByCreatedDesc)
		case "term_asc":
			sb.WriteString(queries.GetTermsSortByNameAsc)
		case "term_desc":
			sb.WriteString(queries.GetTermsSortByNameDesc)
		default:
		}
	}

	rows, err := db.Query(sb.String(), args...)
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

func (t *Term) Update(db SQLExecutor, categoryIds []CategoryId) (*Term, error) {
	if t.Name == "" {
		return nil, errors.New("termname is required")
	}

	_, err := db.Exec(queries.UpdateTerm, t.Name, t.Description, t.UpdatedAt, t.ID)
	if err != nil {
		slog.Error("Failed to update term", "err", err)
		return nil, err
	}

	err = UpdateTermCategories(db, t.ID, categoryIds)
	if err != nil {
		slog.Error("Failed to update term categpory relations", "err", err)
		return nil, err
	}

	return t, nil
}

func (t *Term) Delete(db SQLExecutor) error {
	_, err := db.Exec(queries.DeleteTerm, t.ID)
	if err != nil {
		slog.Error("Failed to delete term", "err", err)
		return err
	}

	_, err = db.Exec(queries.DeleteTermCategoryRelations, t.ID)
	if err != nil {
		slog.Error("Failed to delete term-category relations", "err", err)
		return err
	}

	return nil
}
