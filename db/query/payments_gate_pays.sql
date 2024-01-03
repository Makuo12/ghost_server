-- name: CreatePaymentGatePay :exec
INSERT INTO payments_gate_pays (
    user_id,
    transaction_id,
    reference,
    requested_amount,
    amount,
    payment_gate_fee,
    currency,
    authorization_code,
    payment_gate_paid_at,
    channel,
    payment_gate_created_at
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
    $11
    );
-- name: GetPaymentGatePay :one
SELECT *
FROM payments_gate_pays
WHERE id = $1
LIMIT 1;

-- name: ListPaymentGatePayUser :many
SELECT *
FROM payments_gate_pays
WHERE user_id = $1
ORDER BY created_at
LIMIT $2 OFFSET $3;

