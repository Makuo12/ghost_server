-- name: CreateOptionTripLength :one
INSERT INTO option_trip_lengths (
    option_id
) VALUES (
    $1
) RETURNING option_id;


-- name: UpdateOptionTripLength :one
UPDATE option_trip_lengths 
SET
    min_stay_day = $1,
    max_stay_night = $2,
    manual_approve_request_pass_max = $3,
    allow_reservation_request = $4,
    updated_at = NOW()
WHERE option_id = $5
RETURNING min_stay_day, max_stay_night, manual_approve_request_pass_max, allow_reservation_request;


-- name: GetOptionTripLength :one
SELECT min_stay_day, max_stay_night, manual_approve_request_pass_max, allow_reservation_request
FROM option_trip_lengths
WHERE option_id = $1;


-- name: RemoveOptionTripLength :exec
DELETE FROM option_trip_lengths
WHERE option_id = $1;