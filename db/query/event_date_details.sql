-- name: CreateEventDateDetail :one
INSERT INTO event_date_details (
        event_date_time_id,
        start_time,
        end_time,
        time_zone
    )
VALUES (
        $1, $2, $3, $4
    )
RETURNING *;

-- name: GetEventDateDetail :one
SELECT *
FROM event_date_details
WHERE event_date_time_id = $1;

-- name: GetEventDateDetailByOption :one
SELECT *
FROM event_date_details e_d_d
    JOIN event_date_times e_d_i on e_d_i.id = e_d_d.event_date_time_id
    JOIN options_infos o_i on o_i.id = e_d_i.event_info_id
    JOIN users u on u.id = o_i.host_id
WHERE e_d_d.event_date_time_id = $1 AND u.id = $2 AND o_i.is_complete = $3 AND e_d_i.is_active = true; 

-- name: UpdateEventDateDetail :one
UPDATE event_date_details
SET 
    time_zone = COALESCE(sqlc.narg(time_zone), time_zone),
    start_time = COALESCE(sqlc.narg(start_time), start_time),
    end_time = COALESCE(sqlc.narg(end_time), end_time),
    updated_at = NOW()
WHERE event_date_time_id = sqlc.arg(event_date_time_id) 
RETURNING *;

-- name: RemoveEventDateDetail :exec
DELETE FROM event_date_details
WHERE event_date_time_id=$1;