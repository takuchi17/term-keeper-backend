package queries

const UpdateTermCategoryRelation = `
UPDATE
	term_category_relations
SET
	fk_category_id = ?
WHERE
	fk_term_id = ? AND fk_category_id = ?
`
