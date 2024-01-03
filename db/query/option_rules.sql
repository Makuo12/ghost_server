-- name: CreateOptionRule :one
INSERT INTO option_rules (
        option_id,
        tag,
        type,
        checked,
        start_time,
        end_time
    )
VALUES (
        $1, $2, $3, $4, $5, $6
    )
RETURNING *;

-- name: GetOptionRule :one
SELECT *
FROM option_rules
WHERE id = $1 AND option_id = $2;

-- name: UpdateOptionRule :one
UPDATE option_rules
    SET checked = $1, 
    updated_at = NOW()
WHERE id = $2
RETURNING tag, type, checked, id, des;

-- name: GetOptionRuleDetail :one
SELECT id, tag, type, checked, des
FROM option_rules
WHERE option_id = $1 AND type = $2 AND tag = $3;


-- name: UpdateOptionRuleDetail :one
UPDATE option_rules
SET
    des = COALESCE(sqlc.narg(des), des),
    start_time = COALESCE(sqlc.narg(start_time), start_time),
    end_time = COALESCE(sqlc.narg(end_time), end_time),
    updated_at = NOW()
WHERE id = sqlc.arg(id)  AND option_id = sqlc.arg(option_id) 
RETURNING id, tag, type, checked, des, start_time, end_time;


-- name: GetOptionRuleByType :one
SELECT *
FROM option_rules
WHERE option_id = $1 AND type = $2 AND tag = $3;

-- name: ListOptionRule :many
SELECT *
FROM option_rules
WHERE option_id = $1 AND checked = $2;

-- name: ListAllOptionRule :many
SELECT *
FROM option_rules
WHERE option_id = $1;

-- name: ListOptionRuleTag :many
SELECT tag
FROM option_rules
WHERE option_id = $1 AND checked = $2;

-- name: ListOptionRuleOne :many
SELECT tag, type, checked, id, des
FROM option_rules
WHERE option_id = $1;


-- name: RemoveAllOptionRule :exec
DELETE FROM option_rules
WHERE option_id = $1;

-- name: RemoveOptionRule :exec
DELETE FROM option_rules
WHERE id = $1;