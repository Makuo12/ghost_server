-- name: CreateOptionExtraInfo :one
INSERT INTO options_extra_infos (
    option_id,
    type,
    info
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetOptionExtraInfo :one
SELECT * 
FROM options_extra_infos
WHERE option_id = $1 AND type = $2;

-- name: GetOptionExtraInfoByID :one
SELECT * 
FROM options_extra_infos
WHERE id = $1;

-- name: UpdateOptionExtraInfo :one
UPDATE options_extra_infos
SET 
    info = $1,
    updated_at = NOW()
WHERE option_id = $2 AND type = $3
RETURNING *;

-- name: RemoveOptionExtraInfo :exec
DELETE FROM options_extra_infos 
WHERE option_id = $1;