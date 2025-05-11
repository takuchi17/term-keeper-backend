package models

import (
		"github.com/takuchi17/term-keeper/app/models/queries"
)

type TermCategoryRelation struct {
	FKTermId     TermId
	FKCategoryId CategoryId
}

// termId に紐づく categoryId をすべて取得
func GetCategoryIdsByTermId(termId TermId) ([]string, error) {
	rows, err := DB.Query(queries.GetCategoryIdsByTermId, termId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryIds []string
	for rows.Next() {
		var categoryId string
		if err := rows.Scan(&categoryId); err != nil {
			return nil, err
		}
		categoryIds = append(categoryIds, categoryId)
	}

	return categoryIds, nil
}

// 関連を作成
func CreateTermCategoryRelation(termId TermId, categoryId CategoryId) error {
	_, err := DB.Exec(queries.CreateTermCategoryRelation, termId, categoryId)
	if err != nil {
		return err
	}
	return nil
}

// 関連を削除
func DeleteTermCategoryRelations(termId TermId) error {
	_, err := DB.Exec(queries.DeleteTermCategoryRelations, termId)
	if err != nil {
		return err
	}
	return nil
}
