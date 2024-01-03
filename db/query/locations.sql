-- name: CreateLocation :one
INSERT INTO locations (
   option_id,
   street,
   city,
   state,
   country,
   postcode,
   geolocation,
   show_specific_location
) VALUES (
   $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetLocation :one
SELECT * FROM locations
WHERE option_id = $1
LIMIT 1;


-- name: UpdateLocation :one
UPDATE locations
SET 
   street = COALESCE(sqlc.narg(street), street),
   city = COALESCE(sqlc.narg(city), city),
   state = COALESCE(sqlc.narg(state), state),
   country = COALESCE(sqlc.narg(country), country),
   postcode = COALESCE(sqlc.narg(postcode), postcode),
   geolocation = COALESCE(sqlc.narg(geolocation), geolocation),
   show_specific_location = COALESCE(sqlc.narg(show_specific_location), show_specific_location),
   updated_at = NOW()
WHERE option_id = sqlc.arg(option_id) 
RETURNING *;


-- name: UpdateLocationTwo :one
UPDATE locations
SET 
   street = $1,
   city = $2,
   state = $3,
   country = $4,
   postcode = $5,
   geolocation = $6,
   updated_at = NOW()
WHERE option_id = $7
RETURNING *;

-- name: UpdateSpecificLocation :one
UPDATE locations
SET 
   show_specific_location = $1,
   updated_at = NOW()
WHERE option_id = $2 
RETURNING *;

-- name: DeleteLocation :exec
DELETE FROM locations 
WHERE option_id = $1;

-- name: RemoveLocation :exec
DELETE FROM locations 
WHERE option_id = $1;

--WHERE (point(35.697933, 139.707318) <@> point(longitude, latitude)) < 3