package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTerm(t *testing.T) {
	ctx := context.Background()
	container, err := SetupMysqlContainerAndSetupDB(t, &ctx)
	require.NoError(t, err, "Failed to setup tester.")
	defer container.Terminate(ctx)
	defer DB.Close()

	// まずユーザーを作成しておく（外部キー制約対策）
	userId := TermUserId("01HGDJ5GZRJ2J5VEXR8HT8V9WF")
	_, err = DB.Exec(`INSERT INTO users (id, name, email, password, created_at, updated_at)
					  VALUES (?, ?, ?, ?, ?, ?)`,
		userId, "TestUser", "testuser@example.com", "hashedpassword", time.Now(), time.Now())
	require.NoError(t, err, "Failed to create test user.")

	testCases := []struct {
		name        string
		userId      TermUserId
		termName    TermName
		description TermDescription
		categories  []string
		wantErr     bool
	}{
		{
			name:        "Normal term creation",
			userId:      userId,
			termName:    "Test Term",
			description: "Description of term",
			categories:  []string{},
			wantErr:     false,
		},
		{
			name:        "Empty term name",
			userId:      userId,
			termName:    "",
			description: "Description",
			categories:  []string{},
			wantErr:     true,
		},
		{
			name:        "With categories",
			userId:      userId,
			termName:    "Term with categories",
			description: "Has categories",
			categories:  []string{"cat1", "cat2"},
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ダミーのカテゴリも作成しておく（必要であれば）
			for _, catId := range tc.categories {
				_, err := DB.Exec(`INSERT IGNORE INTO categories (id, user_id, name, created_at, updated_at)
								   VALUES (?, ?, ?, ?, ?)`,
					catId, tc.userId, "dummy", time.Now(), time.Now())
				require.NoError(t, err)
			}

			term, err := CreateTerm(tc.userId, tc.termName, tc.description, tc.categories)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but did not get one.")
				assert.Nil(t, term)
				return
			}

			assert.NoError(t, err, "Expected no error, but got one.")
			assert.NotNil(t, term, "Returned term is nil")

			// 検証: データベースに保存された内容を取得して検証
			var stored Term
			err = DB.QueryRow(`
				SELECT id, user_id, name, description, created_at, updated_at
				FROM terms WHERE id = ?`, term.ID).
				Scan(&stored.ID, &stored.FKUserId, &stored.Name, &stored.Description, &stored.CreatedAt, &stored.UpdatedAt)

			assert.NoError(t, err, "Failed to fetch created term.")
			assert.Equal(t, tc.userId, stored.FKUserId, "UserId mismatch")
			assert.Equal(t, tc.termName, stored.Name, "Name mismatch")
			assert.Equal(t, tc.description, stored.Description, "Description mismatch")

			// カテゴリ関連が存在する場合の検証
			if len(tc.categories) > 0 {
				rows, err := DB.Query(`SELECT category_id FROM term_category_relations WHERE term_id = ?`, term.ID)
				require.NoError(t, err)

				var count int
				for rows.Next() {
					var cid string
					err := rows.Scan(&cid)
					assert.NoError(t, err)
					assert.Contains(t, tc.categories, cid, "Category ID not expected")
					count++
				}
				assert.Equal(t, len(tc.categories), count, "Number of category relations mismatch")
			}
		})
	}
}

func TestUpdateTerm(t *testing.T) {
	ctx := context.Background()
	container, err := SetupMysqlContainerAndSetupDB(t, &ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer DB.Close()

	// 前提となるユーザーとカテゴリを作成
	userId := TermUserId("01HGDJ5GZRJ2J5VEXR8HT8V9WF")
	category1 := "cat1"
	category2 := "cat2"
	_, err = DB.Exec(`INSERT INTO users (id, name, email, password, created_at, updated_at)
					  VALUES (?, ?, ?, ?, ?, ?)`,
		userId, "TestUser", "testuser@example.com", "hashedpassword", time.Now(), time.Now())
	require.NoError(t, err)

	for _, cat := range []string{category1, category2} {
		_, err := DB.Exec(`INSERT INTO categories (id, user_id, name, created_at, updated_at)
						   VALUES (?, ?, ?, ?, ?)`,
			cat, userId, "Category "+cat, time.Now(), time.Now())
		require.NoError(t, err)
	}

	// 初期 Term 作成
	term, err := CreateTerm(userId, "Original Term", "Original Desc", []string{category1})
	require.NoError(t, err)

	// 更新処理
	term.Name = "Updated Term"
	term.Description = "Updated Description"
	term.UpdatedAt = time.Now()

	updated, err := term.Update([]string{category2})
	assert.NoError(t, err)
	assert.Equal(t, "Updated Term", string(updated.Name))
	assert.Equal(t, "Updated Description", string(updated.Description))

	// カテゴリ関連が正しく更新されたか
	rows, err := DB.Query(`SELECT category_id FROM term_category_relations WHERE term_id = ?`, term.ID)
	require.NoError(t, err)
	defer rows.Close()

	var cats []string
	for rows.Next() {
		var cid string
		err := rows.Scan(&cid)
		assert.NoError(t, err)
		cats = append(cats, cid)
	}
	assert.Equal(t, []string{category2}, cats)
}

func TestGetTermsByUserId(t *testing.T) {
	ctx := context.Background()
	container, err := SetupMysqlContainerAndSetupDB(t, &ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer DB.Close()

	userId := TermUserId("01HGDJ5GZRJ2J5VEXR8HT8V9WF")
	category := "searchcat"
	_, err = DB.Exec(`INSERT INTO users (id, name, email, password, created_at, updated_at)
					  VALUES (?, ?, ?, ?, ?, ?)`,
		userId, "SearchUser", "search@example.com", "password", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = DB.Exec(`INSERT INTO categories (id, user_id, name, created_at, updated_at)
					  VALUES (?, ?, ?, ?, ?)`,
		category, userId, "Search Category", time.Now(), time.Now())
	require.NoError(t, err)

	// 複数のTermを作成
	_, _ = CreateTerm(userId, "Alpha Term", "Desc A", []string{category})
	_, _ = CreateTerm(userId, "Beta Term", "Desc B", []string{})
	_, _ = CreateTerm(userId, "Gamma Term", "Desc C", []string{category})

	tests := []struct {
		name     string
		query    *string
		category *string
		sort     *string
		expected int
	}{
		{
			name:     "All terms",
			expected: 3,
		},
		{
			name:     "Filter by name",
			query:    ptrStr("Beta"),
			expected: 1,
		},
		{
			name:     "Filter by category",
			category: &category,
			expected: 2,
		},
		{
			name:     "Sort by name desc",
			sort:     ptrStr("term_desc"),
			expected: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			terms, err := GetTermsByUserId(userId, tc.query, tc.category, tc.sort, nil)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, len(terms), "Unexpected number of terms returned")

			if tc.sort != nil && *tc.sort == "term_desc" {
				assert.True(t, string(terms[0].Name) > string(terms[1].Name), "Terms are not sorted descending")
			}
		})
	}
}

func ptrStr(s string) *string {
	return &s
}
