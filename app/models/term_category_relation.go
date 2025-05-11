package models

import (
	"log/slog"

	"github.com/takuchi17/term-keeper/app/models/queries"
)

func GetCategoryIdsByTermId(db SQLExecutor, termId TermId) ([]CategoryId, error) {
	rows, err := db.Query(queries.GetCategoryIdsByTermId, termId)
	if err != nil {
		slog.Error("Failed to get category ids by term id", "err", err)
		return nil, err
	}
	defer rows.Close()

	var categoryIds []CategoryId
	for rows.Next() {
		var categoryId CategoryId
		if err := rows.Scan(&categoryId); err != nil {
			slog.Error("Failed to scan category id", "err", err)
			return nil, err
		}
		categoryIds = append(categoryIds, categoryId)
	}

	return categoryIds, nil
}

func LinkTermWithCategories(db SQLExecutor, termId TermId, categoryIds []CategoryId) error {
	for _, categoryId := range categoryIds {
		_, err := db.Exec(queries.CreateTermCategoryRelation, termId, categoryId)
		if err != nil {
			slog.Error("Failed to create term-category relation", "err", err)
			return err
		}
	}
	return nil
}

func CreateTermCategoryRelation(db SQLExecutor, termId TermId, categoryId CategoryId) error {
	_, err := db.Exec(queries.CreateTermCategoryRelation, termId, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTermCategoryRelations(db SQLExecutor, termId TermId) error {
	_, err := db.Exec(queries.DeleteTermCategoryRelations, termId)
	if err != nil {
		slog.Error("Failed to delete term-category relations", "err", err)
		return err
	}
	return nil
}

func UpdateTermCategories(db SQLExecutor, termId TermId, categoryIds []CategoryId) error {
	if err := DeleteTermCategoryRelations(db, termId); err != nil {
		return err
	}

	return LinkTermWithCategories(db, termId, categoryIds)
}
