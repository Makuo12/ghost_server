-- name: CreateUserLocation :one
INSERT INTO users_locations (
    user_id,
    street,
    city,
    state,
    country,
    postcode,
    geolocation
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING state, country;


-- name: UpdateUserLocation :one
UPDATE users_locations 
SET
    street = $1,
    city = $2,
    state = $3,
    country = $4,
    postcode = $5,
    geolocation = $6,
    updated_at = NOW()
WHERE user_id = $7
RETURNING state, country;


-- name: GetUserLocationHalf :one
SELECT state, country
FROM users_locations
WHERE user_id = $1;

-- name: RemoveUserLocation :exec
DELETE  FROM users_locations
WHERE user_id = $1;