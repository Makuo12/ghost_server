-- name: CreateEventDateLocation :one
INSERT INTO event_date_locations (
        event_date_time_id,
        street,
        city,
        state,
        country,
        postcode,
        geolocation
    )
VALUES (
        $1, $2, $3, $4, $5, $6, $7
    )
RETURNING *;

-- name: UpdateEventDateLocation :one
UPDATE event_date_locations
SET 
    street = COALESCE(sqlc.narg(street), street),
    city = COALESCE(sqlc.narg(city), city),
    state = COALESCE(sqlc.narg(state), state),
    country = COALESCE(sqlc.narg(country), country),
    postcode = COALESCE(sqlc.narg(postcode), postcode),
    geolocation = COALESCE(sqlc.narg(geolocation), geolocation),
    updated_at = NOW()
WHERE event_date_time_id = sqlc.arg(event_date_time_id) 
RETURNING *;

-- name: UpdateEventDateLocationTwo :one
UPDATE event_date_locations
SET 
    street = $1,
    city = $2,
    state = $3,
    country = $4,
    postcode = $5,
    geolocation = $6,
    updated_at = NOW()
WHERE event_date_time_id = $7
RETURNING *;



-- name: GetEventDateLocation :one
SELECT *
FROM event_date_locations
WHERE event_date_time_id = $1;

-- name: RemoveEventDateLocation :exec
DELETE FROM event_date_locations
WHERE event_date_time_id = $1;

---- name: RemoveAllEventDateLocation :exec
--DELETE FROM event_date_locations
--WHERE event_date_time_id = $1;