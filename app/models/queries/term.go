package queries

const CreateTerm = `
INSERT INTO terms
(
	id,
	fk_user_id,
	name,
	description,
	created_at,
	updated_at
)
VALUES
(
	?,
	?,
	?,
	?,
	?,
	?
)
`

const GetTermsByUserIdBase = `
SELECT 
	t.id, t.fk_user_id, t.name, t.description, t.created_at, t.updated_at
FROM 
	terms t
`

const GetTermsJoinWithCategory = `
INNER JOIN 
	term_category_relations r 
ON
	t.id = r.fk_term_id
`

const GetTermsByUserIdWhere = `
WHERE
	t.fk_user_id = ?
`

const GetTermsFillterByName = `
AND
	t.name LIKE ?
`

const GetTermsFillterByCategory = `
AND
	r.fk_category_id = ?
`

const GetTermsFilterByChecked = `
AND 
	(CASE WHEN ? = true THEN (t.description IS NULL OR t.description = '') 
	      ELSE (t.description IS NOT NULL AND t.description != '')
	END)
`

const GetTermsSortByCreatedAsc = `
ORDER BY 
	t.created_at ASC
`

const GetTermsSortByCreatedDesc = `
ORDER BY 
	t.created_at DESC
`

const GetTermsSortByUpdatedAsc = `
ORDER BY 
	t.updated_at ASC
`

const GetTermsSortByUpdatedDesc = `
ORDER BY 
	t.updated_at Desc
`

const GetTermsSortByNameAsc = `
ORDER BY 
	t.name ASC
`

const GetTermsSortByNameDesc = `
ORDER BY 
	t.name Desc
`

const UpdateTerm = `
UPDATE
	terms
SET
	name = ?, description = ?, category = ?, updated_at = ?
WHERE
	id = ?
`

const DeleteTerm = `
DELETE
FROM
	terms
WHERE
	id = ?
`
