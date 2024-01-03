-- name: CreateChargeEventReference :one
INSERT INTO charge_event_references (
    user_id,
    option_user_id,
    total_fee,
    service_fee,
    total_absorb_fee,
    currency,
    date_booked,
    reference,
    payment_reference,
    is_complete
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;




-- name: UpdateChargeEventReferenceComplete :one
UPDATE charge_event_references
SET 
    is_complete = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: UpdateChargeEventReferenceCompleteByReference :one
UPDATE charge_event_references
SET 
    is_complete = $1,
    updated_at = NOW()
WHERE reference = $2
RETURNING *;
