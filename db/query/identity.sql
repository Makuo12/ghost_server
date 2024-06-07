-- name: CreateIdentity :exec
INSERT INTO identity (
    user_id,
    id_photo_list,
    id_back_photo_list,
    facial_photo_list
) VALUES ($1, $2, $3, $4);

-- name: GetIdentityStatus :one
SELECT status, is_verified
FROM identity
WHERE user_id = $1;

-- name: GetIdentity :one
SELECT *
FROM identity
WHERE user_id = $1;

-- name ListIdentityByAdmin :many
SELECT *
FROM identity;

-- name: UpdateIdentity :one
UPDATE identity 
SET 
    country = $1,
    type = $2,
    id_photo = $3,
    facial_photo = $4,
    status = $5, 
    is_verified = $6,
    id_photo_list = $7,
    facial_photo_list = $8,
    id_back_photo_list = $9,
    updated_at = NOW()
WHERE user_id = $10
RETURNING status, is_verified;

-- name: RemoveIdentity :exec
DELETE FROM identity
WHERE user_id = $1;