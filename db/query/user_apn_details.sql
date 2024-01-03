-- name: CreateUserAPNDetail :one
INSERT INTO user_apn_details (
    user_id,
    device_name,
    model,
    identifier_for_vendor,
    token
) VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateUserAPNDetailToken :exec
UPDATE user_apn_details
SET 
    token = $1,
    updated_at = NOW()
WHERE id = $2;


-- name: ListUidAPNDetail :many
SELECT *
FROM user_apn_details
WHERE user_id = $1;

-- name: ListUserIdAPNDetail :many
SELECT *
FROM user_apn_details ua
JOIN users u on u.id = ua.user_id
WHERE u.user_id = $1;



-- name: RemoveUserAPNDetail :exec
DELETE FROM user_apn_details
WHERE id = $1;

-- name: RemoveAllUserAPNDetail :exec
DELETE FROM user_apn_details
WHERE user_id = $1;

-- name: RemoveAllUserAPNDetailButOne :exec
DELETE FROM user_apn_details
WHERE user_id = $1 AND id != $2;