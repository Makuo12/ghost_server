-- name: CreateEventInfo :one
INSERT INTO event_infos (
        option_id,
        event_type
    )
VALUES (
        $1, $2
    )
RETURNING *;

-- name: GetEventInfo :one
SELECT *
FROM event_infos
WHERE option_id = $1;

-- name: GetEventInfoAnyPolicy :one
SELECT *
FROM event_infos e_i
   JOIN cancel_policies c_p on c_p.option_id = e_i.option_id
WHERE e_i.option_id = $1;

-- name: GetEventInfoByOption :one
SELECT *
FROM event_infos e_i
   JOIN options_infos o_i on o_i.id = e_i.option_id
   JOIN users u on u.id = o_i.host_id
WHERE e_i.option_id = $1 AND u.id = $2 AND o_i.is_complete = $3; 

-- name: UpdateEventInfo :one
UPDATE event_infos
SET 
    event_type = COALESCE(sqlc.narg(event_type), event_type),
    sub_category_type = COALESCE(sqlc.narg(sub_category_type), sub_category_type),
    updated_at = NOW()
WHERE option_id = sqlc.arg(option_id) 
RETURNING *;

-- name: RemoveEventInfo :exec
DELETE FROM event_infos
WHERE option_id=$1;