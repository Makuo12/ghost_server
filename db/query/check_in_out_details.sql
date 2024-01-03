-- name: CreateCheckInOutDetail :one
INSERT INTO check_in_out_details (
    option_id,
    restricted_check_in_days,
    restricted_check_out_days
) VALUES (
    $1, $2, $3
) RETURNING option_id;


-- name: UpdateCheckInOutDetail :one
UPDATE check_in_out_details 
SET
    arrive_after = $1,
    arrive_before = $2,
    leave_before = $3,
    restricted_check_in_days = $4,
    restricted_check_out_days = $5,
    updated_at = NOW()
WHERE option_id = $6
RETURNING arrive_after, arrive_before, leave_before, restricted_check_in_days, restricted_check_out_days;


-- name: GetCheckInOutDetail :one
SELECT arrive_after, arrive_before, leave_before, restricted_check_in_days, restricted_check_out_days
FROM check_in_out_details
WHERE option_id = $1;


-- name: RemoveCheckInOutDetail :exec
DELETE FROM check_in_out_details
WHERE option_id = $1;