-- name: CreateOptionPrice :one
INSERT INTO options_prices (
    option_id,
    price,
    weekend_price
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetOptionPrice :one
SELECT * 
FROM options_prices
WHERE option_id = $1;

-- name: UpdateOptionPrice :one
UPDATE options_prices 
SET price = COALESCE(sqlc.narg(price), price),
    weekend_price = COALESCE(sqlc.narg(weekend_price), weekend_price),
    updated_at = NOW()
WHERE option_id = sqlc.arg(option_id)
RETURNING *;

-- name: RemoveOptionPrice :exec
DELETE FROM options_prices 
WHERE option_id = $1;