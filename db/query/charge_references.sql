-- name: CreateChargeReference :one
INSERT INTO charge_references (
    user_id,
    reason,
    charge,
    currency,
    is_complete,
    reference
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: UpdateChargeReferenceComplete :one
UPDATE charge_references 
SET
    is_complete = $1,
    updated_at = NOW()
WHERE user_id = $2 AND reference = $3
RETURNING *;
