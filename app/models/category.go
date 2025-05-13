package models

import (
	"strings"
	"time"
)

type (
	CategoryId           string
	CategoryUserId       UserId
	CategoryName         string
	CategoryHexColorCode string
)

type Category struct {
	ID           CategoryId
	Name         CategoryName
	FKUserId     CategoryUserId
	HexColorCode CategoryHexColorCode
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

// 単体取得
func GetCategoryById(db SQLExecutor, id CategoryId) (*Category, error) {
	// DBからカテゴリを取得
	row := db.QueryRow("SELECT id, name, fk_user_id, hex_color_code, created_at, updated_at FROM categories WHERE id = ?", id)
	var category Category
	err := row.Scan(&category.ID, &category.Name, &category.FKUserId, &category.HexColorCode, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func GetCategoriesByIds(db SQLExecutor, ids []CategoryId) ([]*Category, error) {
	if len(ids) == 0 {
		return []*Category{}, nil
	}

	// プレースホルダ作成 (?, ?, ...)
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = strings.TrimRight(placeholders, ",")

	query := "SELECT id, name, fk_user_id, hex_color_code, created_at, updated_at FROM categories WHERE id IN (" + placeholders + ")"

	// []string → []interface{}
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.FKUserId, &category.HexColorCode, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}
