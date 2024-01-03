-- name: CreateEventInfoCategory :one
INSERT INTO events_infos_category (
    option_id,
    event_type,
    highlight,
    event_sub_type,
    des,
    name
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetEventInfoCategory :one
SELECT * 
FROM events_infos_category
WHERE option_id = $1;

-- name: UpdateEventInfoCategory :one
UPDATE events_infos_category
SET 
    event_type = $1,
    event_sub_type = $2,
    highlight = $3,
    des = $4,
    name = $5,
    updated_at = NOW()
WHERE option_id = $6
RETURNING *;

-- name: RemoveEventInfoCategory :exec
DELETE FROM events_infos_category 
WHERE option_id = $1;