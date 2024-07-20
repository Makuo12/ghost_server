-- name: CreateChargeReference :one
INSERT INTO charge_references (
    user_id,
    reference,
    object_reference,
    has_object_reference,
    main_object_type,
    payment_medium,
    payment_channel,
    reason,
    charge,
    currency,
    payment_reference,
    is_complete
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetChargeReference :one
SELECT * 
FROM charge_references
WHERE user_id = $1 AND reference = $2;


-- name: UpdateChargeReferenceComplete :one
UPDATE charge_references 
SET
    is_complete = $1,
    updated_at = NOW()
WHERE user_id = $2 AND reference = $3
RETURNING *;
