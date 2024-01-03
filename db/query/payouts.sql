-- name: CreatePayout :one
INSERT INTO payouts (
    payout_ids,
    send_medium,
    user_id,
    amount,
    amount_payed,
    parent_type,
    account_number
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;


-- name: UpdatePayout :one
UPDATE payouts
SET 
    transfer_code = COALESCE(sqlc.narg(transfer_code), transfer_code),
    amount_payed = COALESCE(sqlc.narg(amount_payed), amount_payed),
    is_complete = COALESCE(sqlc.narg(is_complete), is_complete),
    time_paid = COALESCE(sqlc.narg(time_paid), time_paid),
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING payout_ids, user_id, parent_type, account_number, time_paid;