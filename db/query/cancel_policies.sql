-- name: CreateCancelPolicy :one
INSERT INTO cancel_policies (
    option_id
) VALUES (
    $1
) RETURNING option_id;


-- name: UpdateCancelPolicy :one
UPDATE cancel_policies 
SET
    type_one = $1,
    type_two = $2,
    request_a_refund = $3,
    updated_at = NOW()
WHERE option_id = $4
RETURNING type_one, type_two, request_a_refund;


-- name: GetCancelPolicy :one
SELECT type_one, type_two, request_a_refund
FROM cancel_policies
WHERE option_id = $1;


-- name: RemoveCancelPolicy :exec
DELETE FROM cancel_policies
WHERE option_id = $1;