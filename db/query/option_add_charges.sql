-- name: CreateOptionAddCharge :one
INSERT INTO option_add_charges (
    option_id,
    type,
    main_fee,
    extra_fee,
    num_of_guest
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, type, main_fee, extra_fee, num_of_guest;


-- name: UpdateOptionAddCharge :one
UPDATE option_add_charges 
SET
    main_fee = COALESCE(sqlc.narg(main_fee), main_fee),
    extra_fee = COALESCE(sqlc.narg(extra_fee), extra_fee),
    num_of_guest = COALESCE(sqlc.narg(num_of_guest), num_of_guest),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND option_id = sqlc.arg(option_id)
RETURNING id, type, main_fee, extra_fee, num_of_guest;


-- name: UpdateOptionAddChargeByType :one
UPDATE option_add_charges 
SET
    main_fee = $1,
    extra_fee = $2,
    num_of_guest = $3,
    updated_at = NOW()
WHERE type = $4 AND option_id = $5
RETURNING id, type, main_fee, extra_fee, num_of_guest;

-- name: ListOptionAddCharge :many
SELECT id, type, main_fee, extra_fee, num_of_guest
FROM option_add_charges
WHERE option_id = $1;

-- name: GetOptionAddCharge :one
SELECT id, type, main_fee, extra_fee, num_of_guest
FROM option_add_charges
WHERE option_id = $1 AND type=$2;

-- name: RemoveOptionAddChargeByType :exec
DELETE FROM option_add_charges
WHERE option_id = $1 AND type = $2;

-- name: RemoveOptionAddCharge :exec
DELETE FROM option_add_charges
WHERE id = $1;

-- name: RemoveOptionRemoveChargeByOptionID :exec
DELETE FROM option_add_charges
WHERE option_id = $1;
