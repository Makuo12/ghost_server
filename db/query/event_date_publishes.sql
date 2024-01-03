-- name: CreateEventDatePublish :one
INSERT INTO event_date_publishes (
        event_date_time_id,
        event_going_public_date,
        event_going_public_time
    )
VALUES (
        $1, $2, $3
    )
RETURNING *;

-- name: GetEventDatePublish :one
SELECT event_public, event_going_public, event_going_public_time, event_going_public_date, event_date_time_id
FROM event_date_publishes
WHERE event_date_time_id = $1;

-- name: UpdateEventDatePublish :one
UPDATE event_date_publishes
SET 
   event_public = $1,
   event_going_public = $2,
   event_going_public_date = $3,
   event_going_public_time = $4,
   updated_at = NOW()
WHERE event_date_time_id = $5
RETURNING *;

-- name: UpdateEventDatePublishTwo :one
UPDATE event_date_publishes
SET 
   event_public = COALESCE(sqlc.narg(event_public), event_public),
   event_going_public = COALESCE(sqlc.narg(event_going_public), event_going_public),
   event_going_public_date = COALESCE(sqlc.narg(event_going_public_date), event_going_public_date),
   event_going_public_time = COALESCE(sqlc.narg(event_going_public_time), event_going_public_time),
   updated_at = NOW()
WHERE event_date_time_id = sqlc.arg(event_date_time_id)
RETURNING event_public, event_going_public, event_going_public_time, event_going_public_date, event_date_time_id;

-- name: GetEventDatePublishByOption :one
SELECT *
FROM event_date_publishes e_d_p
    JOIN event_date_times e_d_i on e_d_i.id = e_d_p.event_date_time_id
    JOIN options_infos o_i on o_i.id = e_d_p.event_date_time_id
    JOIN users u on u.id = o_i.host_id
WHERE e_d_p.event_date_time_id = $1 AND u.id = $2 AND o_i.is_complete = $3 AND e_d_i.is_active = true; 

-- name: RemoveEventDatePublish :exec
DELETE FROM event_date_publishes
WHERE event_date_time_id = $1;