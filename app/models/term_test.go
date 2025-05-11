package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTerm(t *testing.T) {
	testCases := []struct {
		name        string
		userId      TermUserId
		termName    TermName
		description TermDescription
		categoryIds []CategoryId
		wantErr     bool
	}{
		{
			name:        "Normal term creation",
			userId:      "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			termName:    "Go",
			description: "A programming language created at Google.",
			categoryIds: []CategoryId{"CATE001PROG000000000000001"},
			wantErr:     false,
		},
		{
			name:        "Term with multiple categories",
			userId:      "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			termName:    "PostgreSQL",
			description: "An open-source relational database.",
			categoryIds: []CategoryId{"CATE001PROG000000000000001", "CATE002DBS0000000000000001"},
			wantErr:     false,
		},
		{
			name:        "Term without categories",
			userId:      "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			termName:    "HTML",
			description: "HyperText Markup Language",
			categoryIds: []CategoryId{},
			wantErr:     false,
		},
		{
			name:        "Empty term name",
			userId:      "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			termName:    "",
			description: "This should fail",
			categoryIds: []CategoryId{"CATE001PROG000000000000001"},
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			term, err := CreateTerm(tx, tc.userId, tc.termName, tc.description, tc.categoryIds)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")
			assert.NotNil(t, term, "Term should not be nil")

			// Verify created term
			assert.Equal(t, tc.userId, term.FKUserId, "User ID mismatch")
			assert.Equal(t, tc.termName, term.Name, "Term name mismatch")
			assert.Equal(t, tc.description, term.Description, "Term description mismatch")
			assert.NotEmpty(t, term.ID, "Term ID should not be empty")
			assert.False(t, term.CreatedAt.IsZero(), "CreatedAt should not be zero time")
			assert.False(t, term.UpdatedAt.IsZero(), "UpdatedAt should not be zero time")

			// Verify term in database
			var (
				id          TermId
				userId      TermUserId
				name        TermName
				description TermDescription
				createdAt   time.Time
				updatedAt   time.Time
			)

			err = tx.QueryRow(`
				SELECT id, fk_user_id, name, description, created_at, updated_at
				FROM terms WHERE id = ?
			`, term.ID).Scan(&id, &userId, &name, &description, &createdAt, &updatedAt)

			assert.NoError(t, err, "Failed to get created term from database")
			assert.Equal(t, term.ID, id, "Term ID mismatch in database")
			assert.Equal(t, tc.userId, userId, "User ID mismatch in database")
			assert.Equal(t, tc.termName, name, "Term name mismatch in database")
			assert.Equal(t, tc.description, description, "Term description mismatch in database")

			// Verify category relations
			if len(tc.categoryIds) > 0 {
				rows, err := tx.Query(`
					SELECT fk_category_id FROM term_category_relations
					WHERE fk_term_id = ?
				`, term.ID)

				assert.NoError(t, err, "Failed to query term categories")
				defer rows.Close()

				var foundCategories []CategoryId
				for rows.Next() {
					var categoryId CategoryId
					err := rows.Scan(&categoryId)
					assert.NoError(t, err, "Failed to scan category ID")
					foundCategories = append(foundCategories, categoryId)
				}

				assert.Len(t, foundCategories, len(tc.categoryIds), "Number of categories mismatch")

				// Check if all expected categories are found
				for _, expectedCategoryId := range tc.categoryIds {
					found := false
					for _, foundCategoryId := range foundCategories {
						if expectedCategoryId == foundCategoryId {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected category ID not found: %s", expectedCategoryId)
				}
			}
		})
	}
}

func TestGetTermsByUserId(t *testing.T) {
	testCases := []struct {
		name          string
		userId        TermUserId
		query         *string
		category      *string
		sort          *string
		checked       *bool
		expectedCount int
		expectedTerms []string
		wantErr       bool
	}{
		{
			name:          "Get all terms for user",
			userId:        "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			query:         nil,
			category:      nil,
			sort:          nil,
			checked:       nil,
			expectedCount: 5,
			expectedTerms: []string{"SQL", "TCP/IP", "Docker", "AWS", "TLS"},
			wantErr:       false,
		},
		{
			name:          "Filter terms by name query",
			userId:        "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			query:         stringPtr("SQL"),
			category:      nil,
			sort:          nil,
			checked:       nil,
			expectedCount: 1,
			expectedTerms: []string{"SQL"},
			wantErr:       false,
		},
		{
			name:          "Filter terms by category",
			userId:        "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			query:         nil,
			category:      stringPtr("CATE001PROG000000000000001"),
			sort:          nil,
			checked:       nil,
			expectedCount: 1,
			expectedTerms: []string{"Docker"},
			wantErr:       false,
		},
		{
			name:          "Sort terms by name ascending",
			userId:        "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			query:         nil,
			category:      nil,
			sort:          stringPtr("term_asc"),
			checked:       nil,
			expectedCount: 5,
			expectedTerms: []string{"AWS", "Docker", "SQL", "TCP/IP", "TLS"},
			wantErr:       false,
		},
		{
			name:          "No terms for non-existent user",
			userId:        "NONEXISTENTUSERID000000000",
			query:         nil,
			category:      nil,
			sort:          nil,
			checked:       nil,
			expectedCount: 0,
			expectedTerms: []string{},
			wantErr:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			terms, err := GetTermsByUserId(tx, tc.userId, tc.query, tc.category, tc.sort, tc.checked)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")
			assert.Len(t, terms, tc.expectedCount, "Unexpected number of terms returned")

			if tc.expectedCount > 0 {
				// If sorting is specified, check that the order matches expected
				if tc.sort != nil {
					for i, expectedTerm := range tc.expectedTerms {
						assert.Equal(t, TermName(expectedTerm), terms[i].Name, "Term at position %d doesn't match expected", i)
					}
				} else {
					// Otherwise just check that all expected terms are present
					var termNames []string
					for _, term := range terms {
						termNames = append(termNames, string(term.Name))
					}

					for _, expectedTerm := range tc.expectedTerms {
						assert.Contains(t, termNames, expectedTerm, "Expected term not found: %s", expectedTerm)
					}
				}
			}
		})
	}
}

func TestGetTermsWithCategoriesByUserId(t *testing.T) {
	testCases := []struct {
		name          string
		userId        TermUserId
		query         *string
		category      *string
		sort          *string
		checked       *bool
		expectedCount int
		expectedTerm  string
		expectedCats  []string
		wantErr       bool
	}{
		{
			name:          "Get term with categories",
			userId:        "01HGDJ5GZRJ2J5VEXR8HT8V9WF",
			query:         stringPtr("SQL"),
			category:      nil,
			sort:          nil,
			checked:       nil,
			expectedCount: 1,
			expectedTerm:  "SQL",
			expectedCats:  []string{"データベース"},
			wantErr:       false,
		},
		{
			name:          "Term with multiple categories",
			userId:        "01HGDJ5HXZD3K6WFYS9JU0A1XG", // 佐藤花子's ID
			query:         stringPtr("Python"),
			category:      nil,
			sort:          nil,
			checked:       nil,
			expectedCount: 1,
			expectedTerm:  "Python",
			expectedCats:  []string{"プログラミング", "機械学習"},
			wantErr:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			termAndCategories, err := GetTermsWithCategoriesByUserId(tx, tc.userId, tc.query, tc.category, tc.sort, tc.checked)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")
			assert.Len(t, termAndCategories, tc.expectedCount, "Unexpected number of terms returned")

			if tc.expectedCount > 0 {
				assert.Equal(t, TermName(tc.expectedTerm), termAndCategories[0].Term.Name, "Term name mismatch")

				// Check categories
				var categoryNames []string
				for _, cat := range termAndCategories[0].Categories {
					categoryNames = append(categoryNames, string(cat.Name))
				}

				assert.Len(t, categoryNames, len(tc.expectedCats), "Category count mismatch")

				for _, expectedCat := range tc.expectedCats {
					assert.Contains(t, categoryNames, expectedCat, "Expected category not found: %s", expectedCat)
				}
			}
		})
	}
}

func TestTermUpdate(t *testing.T) {
	testCases := []struct {
		name           string
		termId         TermId
		newName        TermName
		newDescription TermDescription
		newCategories  []CategoryId
		wantErr        bool
	}{
		{
			name:           "Update term name and description",
			termId:         "TERM001SQL000000000000001",
			newName:        "SQL (Updated)",
			newDescription: "Updated description for SQL",
			newCategories:  []CategoryId{"CATE002DBS0000000000000001"},
			wantErr:        false,
		},
		{
			name:           "Update term categories",
			termId:         "TERM002TCP000000000000001",
			newName:        "TCP/IP",
			newDescription: "インターネット通信の基盤となるプロトコル群。",
			newCategories:  []CategoryId{"CATE003NET0000000000000001", "CATE006SEC0000000000000001"},
			wantErr:        false,
		},
		{
			name:           "Empty term name",
			termId:         "TERM003DOCK00000000000001",
			newName:        "",
			newDescription: "Updated description",
			newCategories:  []CategoryId{"CATE001PROG000000000000001"},
			wantErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			// First get the original term
			var term Term
			err = tx.QueryRow(`
				SELECT id, fk_user_id, name, description, created_at, updated_at
				FROM terms WHERE id = ?
			`, tc.termId).Scan(&term.ID, &term.FKUserId, &term.Name, &term.Description, &term.CreatedAt, &term.UpdatedAt)

			assert.NoError(t, err, "Failed to get original term")

			// Apply changes
			term.Name = tc.newName
			term.Description = tc.newDescription
			term.UpdatedAt = time.Now()

			// Update term
			updatedTerm, err := term.Update(tx, tc.newCategories)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")
			assert.NotNil(t, updatedTerm, "Updated term should not be nil")

			// Verify term in database
			var (
				id          TermId
				userId      TermUserId
				name        TermName
				description TermDescription
				createdAt   time.Time
				updatedAt   time.Time
			)

			err = tx.QueryRow(`
				SELECT id, fk_user_id, name, description, created_at, updated_at
				FROM terms WHERE id = ?
			`, tc.termId).Scan(&id, &userId, &name, &description, &createdAt, &updatedAt)

			assert.NoError(t, err, "Failed to get updated term from database")
			assert.Equal(t, tc.termId, id, "Term ID mismatch in database")
			assert.Equal(t, tc.newName, name, "Term name mismatch in database")
			assert.Equal(t, tc.newDescription, description, "Term description mismatch in database")

			// Verify category relations
			rows, err := tx.Query(`
				SELECT fk_category_id FROM term_category_relations
				WHERE fk_term_id = ?
			`, tc.termId)

			assert.NoError(t, err, "Failed to query term categories")
			defer rows.Close()

			var foundCategories []CategoryId
			for rows.Next() {
				var categoryId CategoryId
				err := rows.Scan(&categoryId)
				assert.NoError(t, err, "Failed to scan category ID")
				foundCategories = append(foundCategories, categoryId)
			}

			assert.Len(t, foundCategories, len(tc.newCategories), "Number of categories mismatch")

			// Check if all expected categories are found
			for _, expectedCategoryId := range tc.newCategories {
				found := false
				for _, foundCategoryId := range foundCategories {
					if expectedCategoryId == foundCategoryId {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected category ID not found: %s", expectedCategoryId)
			}
		})
	}
}

func TestTermDelete(t *testing.T) {
	testCases := []struct {
		name    string
		termId  TermId
		wantErr bool
	}{
		{
			name:    "Delete existing term",
			termId:  "TERM004AWS000000000000001",
			wantErr: false,
		},
		{
			name:    "Delete term with multiple categories",
			termId:  "TERM006PYTH00000000000001",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := DB.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			// First verify term exists
			var term Term
			err = tx.QueryRow(`
				SELECT id, fk_user_id, name, description, created_at, updated_at
				FROM terms WHERE id = ?
			`, tc.termId).Scan(&term.ID, &term.FKUserId, &term.Name, &term.Description, &term.CreatedAt, &term.UpdatedAt)

			assert.NoError(t, err, "Term should exist before deletion")

			// Delete term
			err = term.Delete(tx)

			if tc.wantErr {
				assert.Error(t, err, "Expected error, but an error did not occur.")
				return
			}

			assert.NoError(t, err, "Expected no error, but an error occurred.")

			// Verify term is deleted from database
			var count int
			err = tx.QueryRow(`
				SELECT COUNT(*) FROM terms WHERE id = ?
			`, tc.termId).Scan(&count)

			assert.NoError(t, err, "Error counting terms")
			assert.Equal(t, 0, count, "Term should be deleted")

			// Verify category relations are deleted
			err = tx.QueryRow(`
				SELECT COUNT(*) FROM term_category_relations WHERE fk_term_id = ?
			`, tc.termId).Scan(&count)

			assert.NoError(t, err, "Error counting term category relations")
			assert.Equal(t, 0, count, "Term category relations should be deleted")
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
