-- name: CreateAccount :exec
INSERT INTO accounts (
    user_id,
    currency
    )
VALUES (
    $1, $2
    );

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE user_id = $1 AND currency = $2
LIMIT 1;

-- name: ListAccount :many
SELECT *
FROM accounts
WHERE user_id = $1;

---- name: UpdateAccount :one
--UPDATE accounts
--SET balance = $2
--WHERE id = $1
--RETURNING *;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: AddAccountBalanceTopUp :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE user_id = sqlc.arg(user_id) AND currency = sqlc.arg(currency)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;