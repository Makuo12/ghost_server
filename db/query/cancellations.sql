-- name: CreateCancellation :one
INSERT INTO cancellations (
    charge_id,
    charge_type,
    type,
    cancel_user_id,
    reason_one,
    reason_two,
    main_option_type,
    message
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING charge_id;