-- name: CreateOptionMessage :one
INSERT INTO option_messages (
    option_id,
    user_id,
    message,
    type
) VALUES (
    $1, $2, $3, $4
) RETURNING id, option_id, seen, message, user_id, type;


-- name: UpdateOptionMessage :one
UPDATE option_messages 
SET
    seen = COALESCE(sqlc.narg(seen), seen),
    message = COALESCE(sqlc.narg(message), message),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND option_id = sqlc.arg(id)
RETURNING id, option_id, seen, message, user_id, type;


-- name: UpdateOptionMessageByType :one
UPDATE option_messages 
SET
    seen = $1,
    updated_at = NOW()
WHERE option_id
RETURNING id, option_id, seen, message, user_id, type;

-- name: GetOptionMessage :one
SELECT id, option_id, seen, message, user_id, type
FROM option_messages
WHERE option_id = $1 AND id = $2;

-- name: ListOptionMessage :many
SELECT id, option_id, seen, message, user_id, type
FROM option_messages
WHERE option_id = $1;

-- name: ListOptionMessageBySeen :many
SELECT id, option_id, seen, message, user_id, type
FROM option_messages
WHERE option_id = $1 AND seen=$2;

-- name: RemoveAllOptionMessage :exec
DELETE FROM option_messages
WHERE option_id = $1;

-- name: RemoveOptionMessage :exec
DELETE FROM option_messages
WHERE option_id = $1 AND id = $2;

