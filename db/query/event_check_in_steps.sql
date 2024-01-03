-- name: CreateEventCheckInStep :one
INSERT INTO event_check_in_steps (
        event_date_time_id,
        photo,
        des
    )
VALUES (
        $1, $2, $3
    )
RETURNING *;

-- name: GetEventCheckInStep :one
SELECT des, photo
FROM event_check_in_steps
WHERE id = $1 AND event_date_time_id=$2;

-- name: ListEventCheckInStepOrdered :many
SELECT des, photo, id
FROM event_check_in_steps
WHERE event_date_time_id = $1
ORDER BY created_at;

-- name: ListEventCheckInStepPhotos :many
SELECT photo
FROM event_check_in_steps
WHERE event_date_time_id = $1
ORDER BY created_at;

-- name: UpdateEventCheckInStepDes :one
UPDATE event_check_in_steps
    SET des = $1, 
    updated_at = NOW()
WHERE id = $2 AND event_date_time_id = $3
RETURNING des, photo, id;

-- name: UpdateEventCheckInStepPhoto :one
UPDATE event_check_in_steps
    SET photo = $1, 
    updated_at = NOW()
WHERE id = $2 AND event_date_time_id = $3
RETURNING des, photo, id;

-- name: RemoveEventCheckInStep :exec
DELETE FROM event_check_in_steps 
WHERE event_date_time_id = $1 AND id = $2;

-- name: RemoveAllEventCheckInStep :exec
DELETE FROM event_check_in_steps 
WHERE event_date_time_id = $1;