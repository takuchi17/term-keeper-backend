package queries

const UpdateTermCategoryRelation = `
UPDATE
	term_category_relations
SET
	fk_category_id = ?
WHERE
	fk_term_id = ? AND fk_category_id = ?
`

const CreateTermCategoryRelation = `
INSERT INTO term_category_relations
(
	fk_term_id, 
	fk_category_id
) 
VALUES
(
	?, 
	?
)
`

const DeleteTermCategoryRelations = `
DELETE 
FROM 
	term_category_relations 
WHERE
	fk_term_id = ?;
`

const GetCategoryIdsByTermId = `
SELECT 
	fk_category_id 
FROM 
	term_category_relations 
WHERE 
	fk_term_id = ?
`
