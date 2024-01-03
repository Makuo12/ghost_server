-- name: CreateEmContact :one
INSERT INTO em_contacts (
    user_id,
    name,
    relationship,
    email,
    dial_code,
    dial_country,
    phone_number,
    language
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING name, id;


-- name: ListEmContact :many
SELECT id, name
FROM em_contacts
WHERE user_id = $1;

-- name: GetEmContactByPhone :one
SELECT id, name
FROM em_contacts
WHERE phone_number = $1;


-- name: RemoveEmContact :exec
DELETE 
FROM em_contacts
WHERE id = $1;
