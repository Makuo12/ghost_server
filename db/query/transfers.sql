-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount,
    from_account_id_int,
    to_account_id_int
    )
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
    )
RETURNING *;

-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE id = $1
LIMIT 1;