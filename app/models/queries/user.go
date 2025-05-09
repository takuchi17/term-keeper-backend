package queries

const CreateUser = `
INSERT INTO users 
(
  id, 
  name, 
  email, 
  password,
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

const IsDupulicateEmail = `
SELECT 
  COUNT(*)
FROM
  users
WHERE
  email = ?
`

const GetUserById = `
SELECT
  name, email, created_at, updated_at
FROM
  users
WHERE
  id = ?
`

const GetUserByEmail = `
SELECT
  name, id, created_at, updated_at
FROM
  users
WHERE
  email = ?
`
