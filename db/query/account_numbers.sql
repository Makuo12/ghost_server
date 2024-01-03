-- name: CreateAccountNumber :one
INSERT INTO account_numbers (
    user_id,
    account_number,
    account_name,
    bank_name,
    bank_code,
    country,
    currency,
    recipient_code,
    type,
    bank_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetAccountNumberAny :one
SELECT account_number, recipient_code
FROM account_numbers
WHERE user_id = $1;

-- name: GetDefaultAccountNumber :one
SELECT account_number, recipient_code
FROM account_numbers
WHERE user_id = $1 AND id = $2;

-- name: ListAccountNumber :many
SELECT account_number, id, bank_name, currency, account_name 
FROM account_numbers
WHERE user_id = $1;

-- name: RemoveAccountNumber :exec
DELETE  FROM account_numbers WHERE user_id = $1 AND id = $2;
