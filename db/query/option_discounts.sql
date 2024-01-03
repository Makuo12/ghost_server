-- name: CreateOptionDiscount :one
INSERT INTO option_discounts (
    option_id,
    main_type,
    type,
    percent,
    name,
    extra_type,
    des
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, main_type, type, percent, extra_type, name, des;


-- name: UpdateOptionDiscount :one
UPDATE option_discounts 
SET
    main_type = COALESCE(sqlc.narg(main_type), main_type),
    type = COALESCE(sqlc.narg(type), type),
    percent = COALESCE(sqlc.narg(percent), percent),
    extra_type = COALESCE(sqlc.narg(extra_type), extra_type),
    name = COALESCE(sqlc.narg(name), name),
    des = COALESCE(sqlc.narg(des), des),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND option_id = sqlc.arg(id)
RETURNING id, main_type, type, percent, extra_type, name, des;


-- name: UpdateOptionDiscountByType :one
UPDATE option_discounts 
SET
    main_type = $1,
    type = $2,
    percent = $3,
    updated_at = NOW()
WHERE type = $4 AND option_id = $5 AND main_type = $6
RETURNING id, main_type, type, percent, extra_type, name, des;

-- name: ListOptionDiscount :many
SELECT id, main_type, type, percent, extra_type, name, des
FROM option_discounts
WHERE option_id = $1;

-- name: ListOptionDiscountByMainType :many
SELECT id, main_type, type, percent, extra_type, name, des
FROM option_discounts
WHERE option_id = $1 AND main_type=$2;

-- name: GetOptionDiscount :one
SELECT id, main_type, type, percent, extra_type, name, des
FROM option_discounts
WHERE option_id = $1 AND type=$2;

-- name: RemoveOptionDiscountByMainType :exec
DELETE FROM option_discounts
WHERE option_id = $1 AND main_type = $2;

-- name: RemoveOptionDiscount :exec
DELETE FROM option_discounts
WHERE id = $1;

