-- name: CreateEventDatePrivateAudience :one
INSERT INTO event_date_private_audiences (
        event_date_time_id,
        name,
        type,
        email,
        number
    )
VALUES (
        $1, $2, $3, $4, $5
    )
RETURNING *;

-- name: GetEventDatePrivateAudience :one
SELECT *
FROM event_date_private_audiences
WHERE event_date_time_id = $1 AND id=$2;

-- name: ListEventDatePrivateAudience :many
SELECT  name, type, email, number, id, event_date_time_id, sent
FROM event_date_private_audiences
WHERE event_date_time_id = $1;

-- name: UpdateEventDatePrivateAudience :one
UPDATE event_date_private_audiences
SET 
    name = $1,
    type = $2,
    email = $3,
    number = $4,
    updated_at = NOW()
WHERE event_date_time_id = $5 AND id = $6 AND sent = $7
RETURNING name, type, email, number, id, event_date_time_id, sent;


-- name: UpdateEventDatePrivateAudienceSent :one
UPDATE event_date_private_audiences
SET 
    sent = $1,
    updated_at = NOW()
WHERE event_date_time_id = $2 AND id = $3
RETURNING  name, type, email, number, id, event_date_time_id, sent;

-- name: UpdateEventDatePrivateAudienceTwo :one
UPDATE event_date_private_audiences
SET 
    name = COALESCE(sqlc.narg(name), name),
    type = COALESCE(sqlc.narg(type), type),
    email = COALESCE(sqlc.narg(email), email),
    number = COALESCE(sqlc.narg(number), number),
    updated_at = NOW()
WHERE event_date_time_id = sqlc.arg(event_date_time_id) AND id = sqlc.arg(id) AND sent = sqlc.arg(sent)
RETURNING *;

-- name: RemoveEventDatePrivateAudience :exec
DELETE FROM event_date_private_audiences
WHERE event_date_time_id=$1 AND id=$2;

-- name: RemoveEventDatePrivateAudienceBySent :exec
DELETE FROM event_date_private_audiences
WHERE event_date_time_id=$1 AND id=$2 AND sent=$3;

-- name: RemoveAllEventDatePrivateAudience :exec
DELETE FROM event_date_private_audiences
WHERE event_date_time_id=$1;