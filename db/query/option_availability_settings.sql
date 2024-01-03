-- name: CreateOptionAvailabilitySetting :one
INSERT INTO option_availability_settings (
    option_id
) VALUES (
    $1
) RETURNING option_id;


-- name: UpdateOptionAvailabilitySetting :one
UPDATE option_availability_settings 
SET
    advance_notice = $1,
    advance_notice_condition = $2,
    preparation_time = $3,
    availability_window = $4,
    updated_at = NOW()
WHERE option_id = $5
RETURNING advance_notice, advance_notice_condition, preparation_time, availability_window, auto_block_dates;


-- name: GetOptionAvailabilitySetting :one
SELECT advance_notice, advance_notice_condition, preparation_time, availability_window, auto_block_dates
FROM option_availability_settings
WHERE option_id = $1;


-- name: RemoveOptionAvailabilitySetting :exec
DELETE FROM option_availability_settings
WHERE option_id = $1;