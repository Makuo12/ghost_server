-- name: CreateCard :one
INSERT INTO cards (
    user_id,
    email,
    authorization_code,
    card_type,
    last4,
    exp_month,
    exp_year,
    bank,
    currency,
    country_code,
    reusable,
    channel,
    card_signature,
    account_name,
    bin
    )
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15
    )
RETURNING *;

-- name: GetCard :one
SELECT *
FROM cards
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: GetCardAny :one
SELECT *
FROM cards
WHERE user_id = $1
LIMIT 1;

-- name: GetCardByLast4 :one
SELECT id
FROM cards
WHERE last4 = $1 AND currency = $2 AND user_id = $3
LIMIT 1;

-- name: GetCardCount :one
SELECT COUNT(*)
FROM cards
WHERE user_id = $1
LIMIT 1;

-- name: ListCard :many
SELECT * FROM cards
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: RemoveCard :exec
DELETE FROM cards
WHERE user_id = $1 AND id = $2;

