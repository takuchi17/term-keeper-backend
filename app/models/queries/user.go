package queries

const RegisterUser = `
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
  NOW(), 
  NOW()
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
