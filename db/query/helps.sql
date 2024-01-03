-- name: CreateHelp :exec
INSERT INTO helps (
    email,
    subject,
    sub_subject,
    detail
) VALUES ($1, $2, $3, $4);

-- name: GetHelp :many
SELECT * 
FROM helps
WHERE email = $1;

-- name: RemoveHelp :exec
DELETE 
FROM helps
WHERE id = $1;